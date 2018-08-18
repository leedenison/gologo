package gologo

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"regexp"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/pkg/errors"
)

/////////////////////////////////////////////////////////////
// Templates
//

type Template struct {
	Name                string
	Primitive           Primitive
	Renderer            Renderer
	InitialisePrimitive bool
	CloneRenderer       bool
}

func CreateTemplateObject(templateType string, position mgl32.Vec3) (*Object, error) {
	template, ok := templates[templateType]
	if !ok {
		return nil, errors.Errorf("invalid object template: %v", templateType)
	}

	object := CreateObject(position)

	if template.Renderer != nil {
		object.SetRenderer(template.Renderer, template.CloneRenderer)
	}

	if template.Primitive != nil {
		object.SetPrimitive(template.Primitive, true)
		if template.InitialisePrimitive == true {
			err := object.InitialisePrimitive()
			if err != nil {
				return nil, errors.Wrap(err, "Failed to initialise primitive")
			}
		}
	}

	return object, nil
}

// LoadObjectTemplates : read the object config files and
// build the templates. If the supplied path is empty string
// then use the default path
func LoadObjectTemplates(path string) {
	var err error

	if path == "" {
		path = GetResourcePath()
	}

	if err = loadConfigs(path); err != nil {
		Error.Fatalln("Failed to load resources:", err)
	}

	// Set up shaders for defined object types
	if err = configureTemplates(); err != nil {
		Error.Fatalln("Failed to load resources:", err)
	}
}

func configureTemplates() error {
	for _, config := range configs {
		templates[config.Name] = &Template{
			Name: config.Name,
		}

		if config.RendererConfig != nil {
			renderer, err := config.RendererConfig.Create()
			if err != nil {
				return err
			}
			templates[config.Name].Renderer = renderer
			templates[config.Name].CloneRenderer = config.CloneRenderer
		}

		if config.PhysicsPrimitiveConfig != nil {
			primitive, err := config.PhysicsPrimitiveConfig.Create()
			if err != nil {
				return err
			}
			templates[config.Name].Primitive = primitive
			templates[config.Name].InitialisePrimitive = config.InitialisePrimitive
		}
	}

	return nil
}

/////////////////////////////////////////////////////////////
// Template config
//

type TemplateConfig struct {
	Name                   string
	RendererType           string
	Renderer               json.RawMessage
	RendererConfig         RendererConfig
	CloneRenderer          bool
	PhysicsPrimitiveType   string
	PhysicsPrimitive       json.RawMessage
	PhysicsPrimitiveConfig PhysicsPrimitiveConfig
	InitialisePrimitive    bool
}

type RendererConfig interface {
	Create() (Renderer, error)
}

type PhysicsPrimitiveConfig interface {
	Create() (Primitive, error)
}

func RegisterRendererConfig(name string, rendererType reflect.Type) {
	rendererTypes[name] = rendererType
}

func RegisterPhysicsConfig(name string, physicsType reflect.Type) {
	physicsTypes[name] = physicsType
}

func loadConfigs(resourceDir string) error {
	files, err := ioutil.ReadDir(resourceDir)

	if err != nil {
		return errors.Wrap(err, "Failed to load resources.")
	}

	for _, file := range files {
		Trace.Printf("Config file: %v\n", file.Name())
		matched, _ := regexp.MatchString(".*\\.json$", file.Name())
		if file.IsDir() || !matched {
			continue
		}

		filePath := resourceDir + "/" + file.Name()
		config, err := loadConfig(filePath)
		if err != nil {
			Warning.Println("Skipping resource:", err)
			continue
		}

		if config.Name == "" {
			return errors.New("Template is missing required field: 'Name'")
		}

		configs[config.Name] = config
	}

	return nil
}

func loadConfig(resourcePath string) (*TemplateConfig, error) {
	parseResult := TemplateConfig{}

	resourceJSON, err := ioutil.ReadFile(resourcePath)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to read resource file: %s", resourcePath)
	}

	err = json.Unmarshal(resourceJSON, &parseResult)
	if err != nil {
		return nil, errors.Wrapf(err, "Failed to parse resource: %s", resourcePath)
	}

	if parseResult.RendererType != "" && parseResult.RendererType != "NONE" {
		if rendererType, exists := rendererTypes[parseResult.RendererType]; exists {
			untypedConfig := reflect.New(rendererType).Elem().Addr().Interface()
			rendererConfig := untypedConfig.(RendererConfig)
			err = json.Unmarshal(parseResult.Renderer, rendererConfig)
			if err != nil {
				return nil, errors.Wrapf(err, "Failed to parse renderer config: %s", resourcePath)
			}
			parseResult.RendererConfig = rendererConfig
		} else {
			return nil, errors.Errorf("Unknown RenderType: %v\n", parseResult.RendererType)
		}
	}

	if parseResult.PhysicsPrimitiveType != "" && parseResult.PhysicsPrimitiveType != "NONE" {
		if physicsType, exists := physicsTypes[parseResult.PhysicsPrimitiveType]; exists {
			untypedConfig := reflect.New(physicsType).Elem().Addr().Interface()
			physicsConfig := untypedConfig.(PhysicsPrimitiveConfig)
			if parseResult.PhysicsPrimitive != nil {
				err = json.Unmarshal(parseResult.PhysicsPrimitive, &physicsConfig)
				if err != nil {
					return nil, errors.Wrapf(err, "Failed to parse primitive config: %s", resourcePath)
				}
			} else {
				// Catch that we haven't initialised the primitive with values
				parseResult.InitialisePrimitive = true
			}
			parseResult.PhysicsPrimitiveConfig = physicsConfig
		} else {
			return nil, errors.Errorf("Unknown PhysicsPrimitiveType: %v\n", parseResult.PhysicsPrimitiveType)
		}
	}

	return &parseResult, nil
}

/////////////////////////////////////////////////////////////
// Physics primitive config
//

type CircleConfig struct {
	Radius float32
}

func (config *CircleConfig) Create() (Primitive, error) {
	return &Circle{
		Radius: config.Radius,
	}, nil
}

/////////////////////////////////////////////////////////////
// Render config
//

type MeshRendererConfig struct {
	VertexShader   string
	FragmentShader string
	Color          mgl32.Vec4
	Texture        string
	MeshVertices   []float32
}

func (config *MeshRendererConfig) Create() (Renderer, error) {
	var err error
	uniformValues := map[int]interface{}{}

	var uniforms []int

	if config.Texture != "" {
		uniforms = append(uniforms, uniformTexture)
		uniformValues[uniformTexture], err = CreateTexture(config.Texture)
		if err != nil {
			return nil, err
		}
	}

	// TODO: Should this be elsif - which order with texture?
	// TODO: deal with the fact that Vec4 always has a 4 length value and 0,0,0,0 is valid
	if len(config.Color) > 0 {
		uniforms = append(uniforms, uniformColor)
		uniformValues[uniformColor] = config.Color
	}

	err = validateMeshRenderConfig(
		config.VertexShader, config.FragmentShader, config.MeshVertices)
	if err != nil {
		return nil, err
	}

	return CreateMeshRenderer(
		config.VertexShader,
		config.FragmentShader,
		uniforms,
		uniformValues,
		config.MeshVertices)
}

func validateMeshRenderConfig(
	vertexShader string,
	fragmentShader string,
	meshVertices []float32) error {
	if vertexShader == "" {
		return errors.New("Missing required field: 'VertexShader'")
	} else if _, ok := shaders[vertexShader]; !ok {
		return errors.Errorf("Unknown 'VertexShader': %v", vertexShader)
	}

	if fragmentShader == "" {
		return errors.New("Missing required field: 'FragmentShader'")
	} else if _, ok := shaders[fragmentShader]; !ok {
		return errors.Errorf("Unknown 'FragmentShader': %v", fragmentShader)
	}

	if len(meshVertices) == 0 {
		return errors.New("Missing required field: 'MeshVertices'")
	}

	return nil
}

/////////////////////////////////////////////////////////////
// Sprite renderer config
//

type SpriteRendererConfig struct {
	VertexShader   string
	FragmentShader string
	Texture        string
	TextureOrigin  []int32
	MeshScaling    float32
}

func (config *SpriteRendererConfig) Create() (Renderer, error) {
	uniformValues := map[int]interface{}{}

	if config.Texture == "" {
		return nil, errors.New("Missing required field 'Texture'")
	}

	uniform := []int{uniformTexture}
	texture, err := CreateTexture(config.Texture)
	if err != nil {
		return nil, err
	}
	uniformValues[uniformTexture] = texture

	meshVertices := CalcMeshFromSprite(
		float32(config.TextureOrigin[0]),
		float32(config.TextureOrigin[1]),
		float32(texture.Size[0]),
		float32(texture.Size[1]),
		config.MeshScaling)

	return CreateMeshRenderer(
		config.VertexShader,
		config.FragmentShader,
		uniform,
		uniformValues,
		meshVertices)
}

func CalcMeshFromSprite(originX, originY, spriteSizeX, spriteSizeY, scaleFactor float32) []float32 {
	return []float32{
		// Bottom left
		(-spriteSizeX/2 - originX) * scaleFactor,
		(-spriteSizeY/2 + originY) * scaleFactor,
		0.0,
		0.0,
		1.0,
		// Top right
		(spriteSizeX/2 - originX) * scaleFactor,
		(spriteSizeY/2 + originY) * scaleFactor,
		0.0,
		1.0,
		0.0,
		// Top left
		(-spriteSizeX/2 - originX) * scaleFactor,
		(spriteSizeY/2 + originY) * scaleFactor,
		0.0,
		0.0,
		0.0,
		// Bottom left
		(-spriteSizeX/2 - originX) * scaleFactor,
		(-spriteSizeY/2 + originY) * scaleFactor,
		0.0,
		0.0,
		1.0,
		// Bottom right
		(spriteSizeX/2 - originX) * scaleFactor,
		(-spriteSizeY/2 + originY) * scaleFactor,
		0.0,
		1.0,
		1.0,
		// Top right
		(spriteSizeX/2 - originX) * scaleFactor,
		(spriteSizeY/2 + originY) * scaleFactor,
		0.0,
		1.0,
		0.0,
	}
}

/////////////////////////////////////////////////////////////
// Text renderer config
//

type TextRendererConfig struct {
	MeshRenderers map[string]CharRendererConfig
	CharSpacer    float32
}

type CharRendererConfig struct {
	VertexShader   string
	FragmentShader string
	Texture        string
	TextureSize    [2]float32
	TextureRect    [][2]float32
	CharRect       [][2]float32
}

func (config *TextRendererConfig) Create() (Renderer, error) {
	uniformValues := map[int]interface{}{}

	uniform := []int{uniformTexture}

	result := TextRenderer{
		CharWidths:    map[byte]float32{},
		MeshRenderers: map[byte]*MeshRenderer{},
		CharSpacer:    config.CharSpacer,
	}

	for char, charConfig := range config.MeshRenderers {
		bytes := []byte(char)

		if len(bytes) != 1 {
			Trace.Printf("Ignoring multibyte character '%v'\n", char)
			continue
		}

		texture, err := CreateTexture(charConfig.Texture)
		if err != nil {
			return nil, err
		}
		uniformValues[uniformTexture] = texture

		charWidth := charConfig.CharRect[1][0] - charConfig.CharRect[0][0]
		meshVertices := CalcMeshFromChar(
			texture.Size[0],
			texture.Size[1],
			charConfig.TextureRect,
			charConfig.CharRect)

		meshRenderer, err := CreateMeshRenderer(
			charConfig.VertexShader,
			charConfig.FragmentShader,
			uniform,
			uniformValues,
			meshVertices)
		if err != nil {
			return nil, err
		}

		result.CharWidths[bytes[0]] = charWidth
		result.MeshRenderers[bytes[0]] = meshRenderer
	}

	return &result, nil
}

/////////////////////////////////////////////////////////////
// Explosion renderer config
//

type ExplosionRendererConfig struct {
	ParticleCount int
	MaxAge        float32
	MeshRenderers []MeshRendererConfig
}

func (config *ExplosionRendererConfig) Create() (Renderer, error) {
	var err error
	meshRenderers := []*MeshRenderer{}
	uniformValues := map[int]interface{}{}

	uniform := []int{uniformAlpha, uniformTexture}

	for _, meshRendererConfig := range config.MeshRenderers {
		uniformValues[uniformTexture], err = CreateTexture(meshRendererConfig.Texture)
		if err != nil {
			return nil, err
		}

		meshRenderer, err := CreateMeshRenderer(
			meshRendererConfig.VertexShader,
			meshRendererConfig.FragmentShader,
			uniform,
			uniformValues,
			meshRendererConfig.MeshVertices)
		if err != nil {
			return nil, err
		}

		meshRenderers = append(meshRenderers, meshRenderer)
	}

	return &ExplosionRenderer{
		ParticleCount: config.ParticleCount,
		MaxAge:        config.MaxAge,
		MeshRenderers: meshRenderers,
	}, nil
}
