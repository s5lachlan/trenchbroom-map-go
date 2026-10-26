package tb

import (
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path"
	"strings"
)

// Update texture dimensions for accurate UV mapping.
// Extension must be match the extension you use for your textures.
// Returns an error if a texture cannot be found or if config cannot be decoded
func UpdateTextureDimensions(tbMap *Map, folder string, extension string) error {
	tbMap.TextureDimensions = make(map[string]Vector2)

	if !strings.HasSuffix(folder, "/") {
		folder += "/"
	}
	if !strings.HasPrefix(extension, ".") {
		extension += "."
	}

	for _, texture := range tbMap.Textures {
		f, err := os.Open(path.Join(folder, texture+extension))
		if err != nil {
			return err
		}
		defer f.Close()

		config, _, err := image.DecodeConfig(f)
		if err != nil {
			return err
		}

		tbMap.TextureDimensions[texture] = Vector2{(float(config.Width)), float(config.Height)}
	}
	return nil
}

// UV Mapping function, returns U, V.
// Intended for Valve 220 format.
func TextureUV(point Vector3, face BrushFace, textureWidth, textureHeight float) (float, float) {
	textureSize := Vector2{textureWidth, textureHeight}

	uv_x := Vector3{face.UV_S.X, face.UV_S.Y, face.UV_S.Z}
	uv_y := Vector3{face.UV_T.X, face.UV_T.Y, face.UV_T.Z}

	uAxisDot := Vector3DotProduct(point, uv_x)
	vAxisDot := Vector3DotProduct(point, uv_y)

	u := (uAxisDot / (textureSize.X * face.U_Scale)) + (face.UV_S.W / textureSize.X)
	v := (vAxisDot / (textureSize.Y * face.V_Scale)) + (face.UV_T.W / textureSize.Y)

	return u, v
}
