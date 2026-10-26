package tbrl

import (
	"strings"
	tb "trenchbroom-map-go"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// Load all tetures from tbMap into a map of string to rl.Texture2D
func LoadTexturesFromMap(tbMap *tb.Map, folder string, extension string) map[string]rl.Texture2D {
	textures := make(map[string]rl.Texture2D)
	tbMap.TextureDimensions = make(map[string]tb.Vector2)
	if !strings.HasSuffix(folder, "/") {
		folder += "/"
	}
	if !strings.HasPrefix(extension, ".") {
		extension += "."
	}

	for _, tex := range tbMap.Textures {
		rltex := rl.LoadTexture(folder + tex + extension)
		textures[tex] = rltex
		tbMap.TextureDimensions[tex] = tb.Vector2{X: float64(rltex.Width), Y: float64(rltex.Height)}
	}
	return textures
}

// Unload multiple textures, use with defer for cleanliness
func UnloadTextureMap(textures map[string]rl.Texture2D) {
	for _, t := range textures {
		rl.UnloadTexture(t)
	}
}

// Unload multiple models, use with defer for cleanliness.
func UnloadModels(models []rl.Model) {
	for _, m := range models {
		UnloadModel(m)
	}
}

// Draw multiple models.
func DrawModels(models []rl.Model, position rl.Vector3, scale float32, tint rl.Color) {
	for _, m := range models {
		rl.DrawModel(m, position, scale, tint)
	}
}

// Draw multiple models with default parameters
func DrawModelsSimple(models []rl.Model) {
	for _, m := range models {
		rl.DrawModel(m, rl.Vector3Zero(), 1.0, rl.White)
	}
}

func UnloadModel(model rl.Model) {
	// Message from gen2brain, copied here, about Go slices:
	// Vertices, Normals and Texcoords of your CUSTOM mesh are Go slices.
	// UnloadModel calls UnloadMesh for every mesh and UnloadMesh tries
	// to free your Go slices. This will panic because it cannot free
	// Go slices. Free() is a C function and it expects to free C memory
	// and not a Go slice. So clear the slices manually like this.
	model.Meshes.Vertices = nil
	model.Meshes.Normals = nil
	model.Meshes.Texcoords = nil

	rl.UnloadModel(model)
}

// Takes an entity and the texture map and generates models. An entity can have multiple models to render.
func GenModelsFromEntity(entity tb.Entity, textures map[string]rl.Texture2D, tbMap tb.Map) []rl.Model {
	var models []rl.Model
	meshes := tb.BuildEntityMeshes(entity, tbMap)
	if len(meshes) == 0 {
		return []rl.Model{}
	}

	for _, tb := range meshes {

		verts := make([]float32, len(tb.Vertices))
		norms := make([]float32, len(tb.Normals))
		uvs := make([]float32, len(tb.Texcoords))
		for j, v := range tb.Vertices {
			verts[j] = float32(v)
		}
		for j, v := range tb.Normals {
			norms[j] = float32(v)
		}
		for j, v := range tb.Texcoords {
			uvs[j] = float32(v)
		}

		rlmesh := rl.Mesh{
			VertexCount:   int32(tb.VertexCount),
			TriangleCount: int32(tb.TriangleCount),
			Vertices:      unsafe.SliceData(verts),
			Normals:       unsafe.SliceData(norms),
			Texcoords:     unsafe.SliceData(uvs),
		}
		rl.UploadMesh(&rlmesh, false)
		model := rl.LoadModelFromMesh(rlmesh)
		models = append(models, model)

		// One material per texture
		tex := textures[tb.TextureName]
		if tex.ID != 0 {
			rl.SetMaterialTexture(model.Materials, rl.MapAlbedo, tex)
		}
	}

	return models
}
