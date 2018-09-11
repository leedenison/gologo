package gologo

import (
	"encoding/json"
	"io/ioutil"
	"reflect"
	"regexp"

	"github.com/go-gl/mathgl/mgl32"
	"github.com/leedenison/gologo/log"
	"github.com/leedenison/gologo/opengl"
	"github.com/leedenison/gologo/render"
	"github.com/pkg/errors"
)

/////////////////////////////////////////////////////////////
// Templates
//

type Template struct {
	Name                string
	Primitive           Primitive
	Renderer            render.Renderer
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
		log.Error.Fatalln("Failed to load resources:", err)
	}

	// Set up shaders for defined object types
	if err = configureTemplates(); err != nil {
		log.Error.Fatalln("Failed to load resources:", err)
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
	Create() (render.Renderer, error)
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
		log.Trace.Printf("Config file: %v\n", file.Name())
		matched, _ := regexp.MatchString(".*\\.json$", file.Name())
		if file.IsDir() || !matched {
			continue
		}

		filePath := resourceDir + "/" + file.Name()
		config, err := loadConfig(filePath)
		if err != nil {
			log.Warning.Println("Skipping resource:", err)
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

func (config *MeshRendererConfig) Create() (render.Renderer, error) {
	var err error
	uniformValues := map[int]interface{}{}

	var uniforms []int

	if config.Texture != "" {
		uniforms = append(uniforms, opengl.UniformTexture)
		uniformValues[opengl.UniformTexture], err = opengl.CreateTexture(config.Texture)
		if err != nil {
			return nil, err
		}
	}

	// TODO: Should this be elsif - which order with texture?
	// TODO: deal with the fact that Vec4 always has a 4 length value and 0,0,0,0 is valid
	if len(config.Color) > 0 {
		uniforms = append(uniforms, opengl.UniformColor)
		uniformValues[opengl.UniformColor] = config.Color
	}

	return opengl.CreateMeshRenderer(
		config.VertexShader,
		config.FragmentShader,
		uniforms,
		uniformValues,
		config.MeshVertices)
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

func (config *SpriteRendererConfig) Create() (render.Renderer, error) {
	uniformValues := map[int]interface{}{}

	if config.Texture == "" {
		return nil, errors.New("Missing required field 'Texture'")
	}

	uniform := []int{opengl.UniformTexture}
	texture, err := opengl.CreateTexture(config.Texture)
	if err != nil {
		return nil, err
	}
	uniformValues[opengl.UniformTexture] = texture

	meshVertices := CalcMeshFromSprite(
		float32(config.TextureOrigin[0]),
		float32(config.TextureOrigin[1]),
		float32(texture.Size[0]),
		float32(texture.Size[1]),
		config.MeshScaling)

	return opengl.CreateMeshRenderer(
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

func (config *TextRendererConfig) Create() (render.Renderer, error) {
	uniformValues := map[int]interface{}{}

	uniform := []int{opengl.UniformTexture}

	result := TextRenderer{
		CharWidths:    map[byte]float32{},
		MeshRenderers: map[byte]render.Renderer{},
		CharSpacer:    config.CharSpacer,
	}

	for char, charConfig := range config.MeshRenderers {
		bytes := []byte(char)

		if len(bytes) != 1 {
			log.Trace.Printf("Ignoring multibyte character '%v'\n", char)
			continue
		}

		texture, err := opengl.CreateTexture(charConfig.Texture)
		if err != nil {
			return nil, err
		}
		uniformValues[opengl.UniformTexture] = texture

		charWidth := charConfig.CharRect[1][0] - charConfig.CharRect[0][0]
		meshVertices := CalcMeshFromChar(
			texture.Size[0],
			texture.Size[1],
			charConfig.TextureRect,
			charConfig.CharRect)

		meshRenderer, err := opengl.CreateMeshRenderer(
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

func CalcMeshFromChar(
	textureSizeX uint32,
	textureSizeY uint32,
	textureRect [][2]float32,
	charRect [][2]float32) []float32 {

	textureWidth := textureRect[1][0] - textureRect[0][0]
	textureHeight := textureRect[1][1] - textureRect[0][1]
	deltaRight := textureRect[1][0] - charRect[1][0]
	deltaLeft := charRect[0][0] - textureRect[0][0]
	widthDelta := deltaRight - deltaLeft
	deltaTop := charRect[0][1] - textureRect[0][1]
	deltaBottom := textureRect[1][1] - charRect[1][1]
	heightDelta := deltaTop - deltaBottom

	return []float32{
		// Bottom left
		(-float32(textureWidth) + widthDelta) / 2,
		(-float32(textureHeight) + heightDelta) / 2,
		0.0,
		textureRect[0][0] / float32(textureSizeX),
		textureRect[1][1] / float32(textureSizeY),
		// Top right
		(float32(textureWidth) + widthDelta) / 2,
		(float32(textureHeight) + heightDelta) / 2,
		0.0,
		textureRect[1][0] / float32(textureSizeX),
		textureRect[0][1] / float32(textureSizeY),
		// Top left
		(-float32(textureWidth) + widthDelta) / 2,
		(float32(textureHeight) + heightDelta) / 2,
		0.0,
		textureRect[0][0] / float32(textureSizeX),
		textureRect[0][1] / float32(textureSizeY),
		// Bottom left
		(-float32(textureWidth) + widthDelta) / 2,
		(-float32(textureHeight) + heightDelta) / 2,
		0.0,
		textureRect[0][0] / float32(textureSizeX),
		textureRect[1][1] / float32(textureSizeY),
		// Bottom right
		(float32(textureWidth) + widthDelta) / 2,
		(-float32(textureHeight) + heightDelta) / 2,
		0.0,
		textureRect[1][0] / float32(textureSizeX),
		textureRect[1][1] / float32(textureSizeY),
		// Top right
		(float32(textureWidth) + widthDelta) / 2,
		(float32(textureHeight) + heightDelta) / 2,
		0.0,
		textureRect[1][0] / float32(textureSizeX),
		textureRect[0][1] / float32(textureSizeY),
	}
}

/////////////////////////////////////////////////////////////
// Explosion renderer config
//

type ExplosionRendererConfig struct {
	ParticleCount int
	MaxAge        float64
	MeshRenderers []MeshRendererConfig
}

func (config *ExplosionRendererConfig) Create() (render.Renderer, error) {
	var err error
	renderers := []render.Renderer{}
	uniformValues := map[int]interface{}{}

	uniform := []int{opengl.UniformAlpha, opengl.UniformTexture}

	for _, meshRendererConfig := range config.MeshRenderers {
		uniformValues[opengl.UniformTexture], err = opengl.CreateTexture(meshRendererConfig.Texture)
		if err != nil {
			return nil, err
		}

		meshRenderer, err := opengl.CreateMeshRenderer(
			meshRendererConfig.VertexShader,
			meshRendererConfig.FragmentShader,
			uniform,
			uniformValues,
			meshRendererConfig.MeshVertices)
		if err != nil {
			return nil, err
		}

		renderers = append(renderers, meshRenderer)
	}

	return &ExplosionRenderer{
		ParticleCount: config.ParticleCount,
		MaxAge:        config.MaxAge,
		Renderers:     renderers,
	}, nil
}
