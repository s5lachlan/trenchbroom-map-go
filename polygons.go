package tb

import (
	"math"
	"sort"
)

const (
	DEDUP_EPS = 1e-3 // vertex deduplication
)

type Polygon = []Vector3

// Collect every corner of the brush (intersection of three planes that
// still lies inside every half-space).
func CollectBrushVertices(planes []Plane, eps float) []Vector3 {
	n := len(planes)
	var verts []Vector3
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			for k := j + 1; k < n; k++ {
				pt, ok := PlaneIntersection(planes[i], planes[j], planes[k])
				if !ok {
					continue
				}
				if PointInsideBrush(pt, planes, eps) {
					verts = append(verts, pt)
				}
			}
		}
	}
	return Deduplicate(verts, DEDUP_EPS)
}

// Assign corners to faces and order each face into a convex polygon.
func BuildBrushPolygonsFromVertices(planes []Plane, verts []Vector3, eps float) []Polygon {
	n := len(planes)
	faceVerts := make([]Polygon, n)

	for _, v := range verts {
		for i, pl := range planes {
			if PointOnPlane(v, pl, eps) {
				faceVerts[i] = append(faceVerts[i], v)
			}
		}
	}

	polys := make([]Polygon, n)
	for i, pts := range faceVerts {
		pts = Deduplicate(pts, DEDUP_EPS)
		if len(pts) < 3 {
			continue
		}
		polys[i] = OrderConvexPolygon(pts, planes[i])
	}
	return polys
}

func BuildBrushPolygonsAndPlanes(faces []BrushFace, eps float) ([]Polygon, []Plane) {
	n := len(faces)
	planes := make([]Plane, n)
	for i, f := range faces {
		planes[i] = ConvertBrushfaceToPlane(f)
	}

	// Try map winding as-is, outward normals
	verts := CollectBrushVertices(planes, eps)
	if len(verts) >= 4 {
		return BuildBrushPolygonsFromVertices(planes, verts, eps), planes
	}

	// Map winding was inward flip everything and try once
	planes = FlipAllPlanes(planes)
	verts = CollectBrushVertices(planes, eps)
	if len(verts) >= 4 {
		return BuildBrushPolygonsFromVertices(planes, verts, eps), planes
	}

	// Return nothing if brush is degenerate
	return make([]Polygon, n), planes
}

// Check if a point is inside a brush.
func PointInsideBrush(p Vector3, planes []Plane, eps float) bool {
	for _, pl := range planes {
		if Vector3DotProduct(pl.Normal, p)-pl.Distance > eps {
			return false
		}
	}
	return true
}

// Order a convex polygon
func OrderConvexPolygon(pts []Vector3, plane Plane) Polygon {
	if len(pts) < 3 {
		return pts
	}
	// Face centre
	c := Vector3{}
	for _, p := range pts {
		c = Vector3Add(c, p)
	}
	c = Vector3Scale(c, 1.0/float(len(pts)))
	// Build an orthonormal basis in the plane
	var ref Vector3
	if math.Abs(float64(plane.Normal.X)) < 0.9 {
		ref = Vector3{X: 1}
	} else {
		ref = Vector3{Y: 1}
	}
	u := Vector3Normalize(Vector3CrossProduct(ref, plane.Normal))
	v := Vector3CrossProduct(plane.Normal, u)
	// Project & sort by angle
	type AnglePoint struct {
		point Vector3
		angle float
	}
	tmp := make([]AnglePoint, len(pts))
	for i, p := range pts {
		d := Vector3Subtract(p, c)
		tmp[i] = AnglePoint{
			point: p,
			angle: math.Atan2(float64(Vector3DotProduct(d, v)),
				float64(Vector3DotProduct(d, u))),
		}
	}
	sort.Slice(tmp, func(i, j int) bool {
		return tmp[i].angle < tmp[j].angle
	})

	ordered := make([]Vector3, len(tmp))
	for i, a := range tmp {
		ordered[i] = a.point
	}
	return ordered
}

// Deduplicate removes vertices that are closer than `eps` to an earlier one.
// Order is preserved (first occurrence wins).
func Deduplicate(pts []Vector3, eps float) Polygon {
	if len(pts) <= 1 {
		return pts
	}

	eps2 := eps * eps // compare squared distances to avoid sqrt
	result := make([]Vector3, 0, len(pts))

	for _, p := range pts {
		tooClose := false
		for _, q := range result {
			dx := p.X - q.X
			dy := p.Y - q.Y
			dz := p.Z - q.Z
			if dx*dx+dy*dy+dz*dz < eps2 {
				tooClose = true
				break
			}
		}
		if !tooClose {
			result = append(result, p)
		}
	}
	return result
}
