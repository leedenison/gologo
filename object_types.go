package gologo

import (
    "github.com/pkg/errors"
)

type ObjectType struct {
    Name string
    Primitive Primitive
    Renderer Renderer
}

func InitObjectTypes(configs map[string]*ObjectTypeConfig) (
        map[string]*ObjectType,
        error) {
    result := map[string]*ObjectType {}

    for _, config := range configs {
        var err error

        var renderer Renderer
        switch rendererConfig := config.RendererConfig.(type) {
        case GLMeshRendererConfig:
            renderer, err = InitGLMeshRenderer(&rendererConfig)
        case SpriteMeshRendererConfig:
            renderer, err = InitSpriteMeshRenderer(&rendererConfig)
        default:
            return nil, errors.Errorf("Unhandled RenderType: %v\n", config.RendererType)
        }

        if err != nil {
            return nil, errors.Wrapf(
                err,
                "Invalid ObjectTypeConfig[%s]: RendererType[%s]",
                config.Name,
                config.RendererType)
        }

        result[config.Name] = &ObjectType {
            Name: config.Name,
            Renderer: renderer,
        }

        var primitive Primitive
        switch config.PhysicsPrimitiveType {
        case SPRITE_CIRCLE:
            if meshRenderer, ok := renderer.(*GLMeshRenderer); ok {
                primitive = InitCircleFromMesh(meshRenderer.MeshVertices)
            } else {
                return nil, errors.Errorf(
                    "Cannot use SPRITE_CIRCLE primitive with RendererType: %t\n", renderer)
            }
        default:
            return nil, errors.Errorf("Unhandled PhysicsPrimitiveType: %v\n",
                config.PhysicsPrimitiveType)
        }

        if err != nil {
            return nil, errors.Wrapf(
                err,
                "Invalid ObjectTypeConfig[%s]: PhysicsPrimitiveType[%s]",
                config.Name,
                config.PhysicsPrimitiveType)
        }

        result[config.Name].Primitive = primitive
    }

    return result, nil
}

