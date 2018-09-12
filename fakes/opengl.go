package fakes

import "github.com/leedenison/gologo/opengl"

// CreateVerticesOnlyMeshRendererImpl : Creates a MeshRenderer that only
// contains vertices and a vertexCount. All other values are nil or 0
func CreateVerticesOnlyMeshRendererImpl(
	vertexShader string,
	fragmentShader string,
	uniforms []int,
	uniformValues map[int]interface{},
	meshVertices []float32) (*opengl.MeshRenderer, error) {
	return &opengl.MeshRenderer{
		Shader:       nil,
		Mesh:         0,
		Uniforms:     nil,
		MeshVertices: meshVertices,
		VertexCount:  int32(len(meshVertices) / opengl.GlMeshStride),
	}, nil
}
