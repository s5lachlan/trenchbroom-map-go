package tb

import "math"

type float = float64

type Vector2 struct {
	X float
	Y float
}

type Vector3 struct {
	X float
	Y float
	Z float
}

type Vector4 struct {
	X float
	Y float
	Z float
	W float
}

func NewVector3(x, y, z float) Vector3 {
	return Vector3{x, y, z}
}

func NewVector4(x, y, z, w float) Vector4 {
	return Vector4{x, y, z, w}
}

func Vector3Zero() Vector3 {
	return Vector3{X: 0.0, Y: 0.0, Z: 0.0}
}

func Vector3Length(v Vector3) float {
	return float(
		math.Sqrt(
			float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z),
		),
	)
}

func Vector3CrossProduct(v1 Vector3, v2 Vector3) Vector3 {
	return Vector3{
		v1.Y*v2.Z - v1.Z*v2.Y,
		v1.Z*v2.X - v1.X*v2.Z,
		v1.X*v2.Y - v1.Y*v2.X,
	}
}

func Vector3Scale(v Vector3, scalar float64) Vector3 {
	return Vector3{
		v.X * scalar,
		v.Y * scalar,
		v.Z * scalar,
	}
}

func Vector3DotProduct(v1 Vector3, v2 Vector3) float {
	return float(v1.X*v2.X + v1.Y*v2.Y + v1.Z*v2.Z)
}

func Vector3Append(dst []float64, v Vector3) []float64 {
	return append(dst, v.X, v.Y, v.Z)
}

func Vector3Add(vectors ...Vector3) Vector3 {
	output := Vector3Zero()
	for _, vec3 := range vectors {
		output.X += vec3.X
		output.Y += vec3.Y
		output.Z += vec3.Z
	}
	return output
}

func Vector3Subtract(v1 Vector3, v2 Vector3) Vector3 {
	return Vector3{
		v1.X - v2.X,
		v1.Y - v2.Y,
		v1.Z - v2.Z,
	}
}

func Vector3Negate(v Vector3) Vector3 {
	return Vector3{-v.X, -v.Y, -v.Z}
}

func Vector3Normalize(v Vector3) Vector3 {
	result := v

	length := math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
	if length != 0.0 {
		ilength := 1.0 / length

		result.X *= ilength
		result.Y *= ilength
		result.Z *= ilength
	}

	return result
}
