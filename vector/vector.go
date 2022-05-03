package vector

import "math"

// Vector - struct holding X Y Z values of a 3D vector
type Vector struct {
	X, Y, Z float64
}

// float operator*(const vec3& v) const { return x*v.x + y*v.y + z*v.z; }
// vec3  operator-()              const { return {-x, -y, -z};          }
// float norm() const { return std::sqrt(x*x+y*y+z*z); }

// vec3  operator+(const vec3& v) const { return {x+v.x, y+v.y, z+v.z}; }

func (a Vector) Add(b Vector) Vector {
	return Vector{
		X: a.X + b.X,
		Y: a.Y + b.Y,
		Z: a.Z + b.Z,
	}
}

// vec3  operator-(const vec3& v) const { return {x-v.x, y-v.y, z-v.z}; }

func (a Vector) Sub(b Vector) Vector {
	return Vector{
		X: a.X - b.X,
		Y: a.Y - b.Y,
		Z: a.Z - b.Z,
	}
}

// vec3  operator*(const float v) const { return {x*v, y*v, z*v};       }

func (a Vector) MulS(s float64) Vector {
	return Vector{
		X: a.X * s,
		Y: a.Y * s,
		Z: a.Z * s,
	}
}

// float operator*(const vec3& v) const { return x*v.x + y*v.y + z*v.z; }
func (a Vector) Mul(b Vector) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a Vector) Dot(b Vector) float64 {
	return a.X*b.X + a.Y*b.Y + a.Z*b.Z
}

func (a Vector) Length() float64 {
	return math.Sqrt(a.Dot(a))
}

// float norm() const { return std::sqrt(x*x+y*y+z*z); }
func (a Vector) Norm() float64 {
	return a.Dot(a)
}

func (a Vector) Cross(b Vector) Vector {
	return Vector{
		X: a.Y*b.Z - a.Z*b.Y,
		Y: a.Z*b.X - a.X*b.Z,
		Z: a.X*b.Y - a.Y*b.X,
	}
}

// vec3 normalized() const { return (*this)*(1.f/norm()); }

func (a Vector) Normalize() Vector {
	return a.MulS(1. / a.Length())
}

type Material struct {
	DiffuseColor Vector
}

type Sphere struct {
	Center   Vector
	Radius   float64
	Material Material
}

func (s *Sphere) RayIntersect(orig Vector, dir Vector, t0 *float64) bool {
	L := s.Center.Sub(orig)
	tca := L.Mul(dir)
	d2 := L.Mul(L) - tca*tca
	if d2 > s.Radius*s.Radius {
		return false
	}
	thc := math.Sqrt(s.Radius*s.Radius - d2)
	*t0 = tca - thc
	t1 := tca + thc
	if *t0 < 0 {
		*t0 = t1
	}
	if *t0 < 0 {
		return false
	}
	return true
}

func SceneIntersect(orig Vector, dir Vector, spheres []Sphere, hit *Vector, N *Vector, material *Material) bool {
	spheres_dist := math.MaxFloat64
	for i := 0; i < len(spheres); i++ {
		var distI float64 = 0.0
		if spheres[i].RayIntersect(orig, dir, &distI) && distI < spheres_dist {
			spheres_dist = distI
			*hit = orig.Add(dir).MulS(distI)
			*N = (hit.Sub(spheres[i].Center)).Normalize()
			*material = spheres[i].Material
		}
	}
	return spheres_dist < 1000
}

func CastRay(orig Vector, dir Vector, spheres []Sphere, lights []Light) Vector {
	point := Vector{0, 0, 0}
	N := Vector{0, 0, 0}
	material := Material{}

	if !SceneIntersect(orig, dir, spheres, &point, &N, &material) {
		return Vector{X: 0.2, Y: 0.7, Z: 0.8} // background color
	}
	//return material.DiffuseColor
	var diffuseLightIntensity float64 = 0.0
	for i := 0; i < len(lights); i++ {
		light_dir := lights[i].Position.Sub(point).Normalize()
		diffuseLightIntensity += lights[i].Intensity * math.Max(0, light_dir.Mul(N))
	}
	return material.DiffuseColor.MulS(diffuseLightIntensity)
}

type Light struct {
	Position  Vector
	Intensity float64
}
