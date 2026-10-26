package tb

import "math"

type Plane struct {
	Normal   Vector3
	Distance float
}

// Converts a Brushface to a Plane
func ConvertBrushfaceToPlane(face BrushFace) Plane {
	edge1 := Vector3Subtract(face.Pos2, face.Pos1)
	edge2 := Vector3Subtract(face.Pos3, face.Pos1)
	crossproduct := Vector3CrossProduct(edge1, edge2)
	normal := Vector3Normalize(crossproduct)
	distance := Vector3DotProduct(normal, face.Pos1)
	return Plane{normal, distance}
}

func FlipPlane(p Plane) Plane {
	return Plane{Normal: Vector3Negate(p.Normal), Distance: -p.Distance}
}

func FlipAllPlanes(planes []Plane) []Plane {
	out := make([]Plane, len(planes))
	for i, p := range planes {
		out[i] = FlipPlane(p)
	}
	return out
}

func PointOnPlane(p Vector3, pl Plane, eps float) bool {
	return math.Abs(float64(Vector3DotProduct(pl.Normal, p)-pl.Distance)) <= float64(eps)
}

// Orients a plane.
func OrientPlane(p Plane, centre Vector3) Plane {
	if Vector3DotProduct(p.Normal, centre)-p.Distance > 0 {
		p.Normal = Vector3Negate(p.Normal)
		p.Distance = -p.Distance
	}
	return p
}

// Intersection of three planes. Returns (point, ok).
func PlaneIntersection(a, b, c Plane) (Vector3, bool) {
	// Solve the 3×3 system N·X = D
	n1, n2, n3 := a.Normal, b.Normal, c.Normal
	denom := Vector3DotProduct(n1, Vector3CrossProduct(n2, n3))
	if math.Abs(float64(denom)) < 1e-8 {
		return Vector3{}, false // parallel / coincident
	}

	// Cramer’s rule / vector form
	// X = (D1(N2×N3) + D2(N3×N1) + D3(N1×N2)) / denom
	x := Vector3Scale(Vector3CrossProduct(n2, n3), a.Distance)
	x = Vector3Add(x, Vector3Scale(Vector3CrossProduct(n3, n1), b.Distance))
	x = Vector3Add(x, Vector3Scale(Vector3CrossProduct(n1, n2), c.Distance))
	x = Vector3Scale(x, 1.0/denom)
	return x, true
}
