package tb

type Mesh struct {
	// The name of this mesh's texture
	TextureName string
	// Number of vertices stored in arrays
	VertexCount int
	// Number of triangles stored (indexed or not)
	TriangleCount int
	// Vertex position (XYZ - 3 components per vertex) (shader-location = 0)
	Vertices []float
	// Vertex normals (XYZ - 3 components per vertex) (shader-location = 2)
	Normals []float
	// Texture coords
	Texcoords []float
}

// Builds meshes for an entity. An entity can have multiple meshes.
func BuildEntityMeshes(entity Entity, tbMap Map) []Mesh {
	// texture name → accumulating mesh
	textureNameToMeshMap := map[string]*Mesh{}

	for _, brush := range entity.Brushes {
		if len(brush.BrushFaces) < 4 {
			continue
		}

		polys, planes := BuildBrushPolygonsAndPlanes(brush.BrushFaces, tbMap.Epsilon)

		for i, poly := range polys {
			if len(poly) < 3 {
				continue
			}
			face := brush.BrushFaces[i]
			texName := face.Texture
			if texName == "" {
				texName = "__TB_empty"
			}

			m, ok := textureNameToMeshMap[texName]
			if !ok {
				m = &Mesh{}
				textureNameToMeshMap[texName] = m
			}

			normal := planes[i].Normal
			for _, tri := range Triangulate(poly) {
				tri = OrientTriangle(tri, normal)

				e1 := Vector3Subtract(tri[1], tri[0])
				e2 := Vector3Subtract(tri[2], tri[0])
				if Vector3Length(Vector3CrossProduct(e1, e2)) < 1e-4 {
					continue
				}

				for _, p := range tri {
					m.Vertices = Vector3Append(m.Vertices, p)
					textureWidth := tbMap.TextureDimensions[texName].X
					textureHeight := tbMap.TextureDimensions[texName].Y
					u, v := TextureUV(p, face, textureWidth, textureHeight)
					m.Texcoords = append(m.Texcoords, u, v)
					m.Normals = Vector3Append(m.Normals, normal)
				}
				m.TriangleCount++
			}
		}
	}

	out := make([]Mesh, 0, len(textureNameToMeshMap))
	for name, m := range textureNameToMeshMap {
		m.VertexCount = m.TriangleCount * 3
		if m.VertexCount == 0 {
			continue
		}

		m.TextureName = name
		ConvertMeshToXYZ(m)
		out = append(out, *m)
	}
	return out
}

// Converts a mesh to left-handed coordinate space usually used for Raylib and others.
func ConvertMeshToXYZ(m *Mesh) {
	for i := 0; i+2 < len(m.Vertices); i += 3 {
		x, y, z := m.Vertices[i], m.Vertices[i+1], m.Vertices[i+2]
		m.Vertices[i] = x
		m.Vertices[i+1] = z
		m.Vertices[i+2] = -y
	}

	// Normals (same rotation)
	for i := 0; i+2 < len(m.Normals); i += 3 {
		x, y, z := m.Normals[i], m.Normals[i+1], m.Normals[i+2]
		m.Normals[i] = x
		m.Normals[i+1] = z
		m.Normals[i+2] = -y
	}
}
