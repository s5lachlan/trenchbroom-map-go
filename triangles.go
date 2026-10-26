package tb

// A triangle is comprised of 3 points in 3D space.
type Triangle [3]Vector3

// Orient a triangle
func OrientTriangle(tri Triangle, normal Vector3) Triangle {
	edge1 := Vector3Subtract(tri[1], tri[0])
	edge2 := Vector3Subtract(tri[2], tri[0])

	if Vector3DotProduct(Vector3CrossProduct(edge1, edge2), normal) < 0 {
		tri[1], tri[2] = tri[2], tri[1]
	}

	return tri
}

// Triangulate polygons
func Triangulate(poly []Vector3) []Triangle {
	if len(poly) < 3 {
		return nil
	}

	triangles := make([]Triangle, 0, len(poly)-2)

	for i := 1; i < len(poly)-1; i++ {
		triangles = append(triangles, Triangle{
			poly[0],
			poly[i],
			poly[i+1],
		})
	}

	return triangles
}
