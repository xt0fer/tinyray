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
	Albedo       [2]float64
	DiffuseColor Vector
	SpecularExp  float64
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
	var specularLightIntensity float64 = 0.0

	for i := 0; i < len(lights); i++ {
		light_dir := lights[i].Position.Sub(point).Normalize()
		light_distance := lights[i].Position.Sub(point).Norm()

		shadow_orig := Vector{0, 0, 0}
		if light_dir.Mul(N) < 0 {
			shadow_orig = point.Sub(N.MulS(1e-3))
		} else {
			shadow_orig = point.Add(N.MulS(1e-3))
		}
		// checking if the point lies in the shadow of the lights[i]

		shadow_pt := Vector{0, 0, 0}
		shadow_N := Vector{0, 0, 0}

		tmpmaterial := Material{}
		if SceneIntersect(shadow_orig, light_dir, spheres, &shadow_pt, &shadow_N, &tmpmaterial) && (shadow_pt.Sub(shadow_orig)).Norm() < light_distance {
			continue
		}
		diffuseLightIntensity += lights[i].Intensity * math.Max(0, light_dir.Mul(N))
		// specular_light_intensity += powf(std::max(0.f, -reflect(-light_dir, N)*dir), material.specular_exponent)*lights[i].intensity;
		specularLightIntensity += math.Pow(math.Max(0.0, Reflect(light_dir, N).Mul(dir)),
			material.SpecularExp) * lights[i].Intensity
	}
	// * material.albedo[0] + Vec3f(1., 1., 1.)*specular_light_intensity * material.albedo[1];
	return material.DiffuseColor.MulS(diffuseLightIntensity).MulS(material.Albedo[0]).Add(
		Vector{X: 1, Y: 1, Z: 1}.MulS(specularLightIntensity * material.Albedo[1]))
}

type Light struct {
	Position  Vector
	Intensity float64
}

func Reflect(I Vector, N Vector) Vector {
	return I.Sub(N.MulS(2.0 * I.Mul(N)))
}
