package vector

import "math"

// Vector - struct holding X Y Z values of a 3D vector
type Vector struct {
	X, Y, Z float64
}
type V4 struct {
	X, Y, Z, A float64
}

type Material struct {
	RefractIdx   float64
	Albedo       V4
	DiffuseColor Vector
	SpecularExp  float64
}

type Sphere struct {
	Center   Vector
	Radius   float64
	Material Material
}

type Light struct {
	Position  Vector
	Intensity float64
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
func (a Vector) Neg() Vector {
	return Vector{
		X: -a.X,
		Y: -a.Y,
		Z: -a.Z,
	}
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

func (s *Sphere) RayIntersect(orig Vector, dir Vector, t0 *float64) bool {
	L := s.Center.Sub(orig)
	tca := L.Dot(dir)
	d2 := L.Dot(L) - tca*tca
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
		distI := 0.0
		if spheres[i].RayIntersect(orig, dir, &distI) && distI < spheres_dist {
			spheres_dist = distI
			*hit = orig.Add(dir.MulS(distI))
			*N = (hit.Sub(spheres[i].Center)).Normalize()
			*material = spheres[i].Material
		}
	}

	checkerboard_dist := math.MaxFloat64
	if math.Abs(dir.Y) > 1e-3 {
		d := -(orig.Y + 4) / dir.Y // the checkerboard plane has equation y = -4
		pt := orig.Add(dir.MulS(d))
		if d > 0 && math.Abs(pt.X) < 10.0 && pt.Z < -10.0 && pt.Z > -30.0 && d < spheres_dist {
			checkerboard_dist = d
			*hit = pt
			*N = Vector{X: 0, Y: 1, Z: 0}
			if (int(0.5*hit.X+1000)+int(0.5*hit.Z))&1 == 1 {
				material.DiffuseColor = Vector{X: 0.3, Y: 0.3, Z: 0.3}
			} else {
				material.DiffuseColor = Vector{X: 0.3, Y: 0.2, Z: 0.1}
			}
		}
	}
	return math.Min(spheres_dist, checkerboard_dist) < 1000

}

func CastRay(orig Vector, dir Vector, spheres []Sphere, lights []Light, depth int) Vector {
	point := Vector{0, 0, 0}
	N := Vector{0, 0, 0}

	material := Material{
		RefractIdx:   1.0,
		Albedo:       V4{X: 2.0, Y: 0, Z: 0, A: 0},
		DiffuseColor: Vector{X: 0.0, Y: 0, Z: 0},
		SpecularExp:  0.0,
	}

	if depth > 4 || !SceneIntersect(orig, dir, spheres, &point, &N, &material) {
		return Vector{X: 0.2, Y: 0.7, Z: 0.8} // background color
	}

	reflect_dir := Reflect(dir, N).Normalize()
	refract_dir := Refract(dir, N, material.RefractIdx, 1.0).Normalize()

	reflect_orig := Vector{0, 0, 0}
	if reflect_dir.Dot(N) < 0 {
		reflect_orig = point.Sub(N.MulS(1e-3))
	} else {
		reflect_orig = point.Add(N.MulS(1e-3))
	}
	refract_orig := Vector{0, 0, 0}
	if reflect_dir.Dot(N) < 0 {
		refract_orig = point.Sub(N.MulS(1e-3))
	} else {
		refract_orig = point.Add(N.MulS(1e-3))
	}
	reflect_color := CastRay(reflect_orig, reflect_dir, spheres, lights, depth+1)
	refract_color := CastRay(refract_orig, refract_dir, spheres, lights, depth+1)

	var diffuseLightIntensity float64 = 0.0
	var specularLightIntensity float64 = 0.0

	for i := 0; i < len(lights); i++ {
		light_dir := lights[i].Position.Sub(point).Normalize()
		light_distance := lights[i].Position.Sub(point).Norm()

		shadow_orig := Vector{0, 0, 0}
		if light_dir.Dot(N) < 0 {
			shadow_orig = point.Sub(N.MulS(1e-3))
		} else {
			shadow_orig = point.Add(N.MulS(1e-3))
		}
		// checking if the point lies in the shadow of the lights[i]

		shadow_pt := Vector{0, 0, 0}
		shadow_N := Vector{0, 0, 0}

		tmpmaterial := Material{
			RefractIdx:   1.0,
			Albedo:       V4{X: 2.0, Y: 0, Z: 0, A: 0},
			DiffuseColor: Vector{X: 0.0, Y: 0, Z: 0},
			SpecularExp:  0.0,
		}
		if SceneIntersect(shadow_orig, light_dir, spheres, &shadow_pt, &shadow_N, &tmpmaterial) &&
			(shadow_pt.Sub(shadow_orig)).Norm() < light_distance {
			continue
		}
		diffuseLightIntensity += lights[i].Intensity * math.Max(0, light_dir.Dot(N))
		mLightDir := light_dir.MulS(-1.0)
		specularLightIntensity += math.Pow(math.Max(0.0, Reflect(mLightDir, N).MulS(-1.0).Dot(dir)),
			material.SpecularExp) * lights[i].Intensity
	}

	tcolor := material.DiffuseColor.MulS(diffuseLightIntensity).MulS(material.Albedo.X)
	tcolor = tcolor.Add(Vector{X: 1, Y: 1, Z: 1}.MulS(specularLightIntensity * material.Albedo.Y))
	tcolor = tcolor.Add(reflect_color.MulS(material.Albedo.Z))
	tcolor = tcolor.Add(refract_color.MulS(material.Albedo.A))
	return tcolor
	// return material.DiffuseColor.MulS(diffuseLightIntensity).MulS(material.Albedo.X).Add(
	// 	Vector{X: 1, Y: 1, Z: 1}.MulS(specularLightIntensity * material.Albedo.Y).Add(
	// 		reflect_color.MulS(material.Albedo.Z)))
}

func Reflect(I Vector, N Vector) Vector {
	return I.Sub(N.MulS(2.0).MulS(I.Dot(N)))
}

func Refract(I Vector, N Vector, eta_t float64, eta_i float64) Vector { // Snell's law
	cosi := -math.Max(-1.0, math.Min(1.0, I.Dot(N)))

	if cosi < 0 {
		return Refract(I, N.Neg(), eta_i, eta_t)
	} // if the ray comes from the inside the object, swap the air and the media
	eta := eta_i / eta_t

	k := 1 - eta*eta*(1-cosi*cosi)
	tmpk := Vector{1, 0, 0}
	if !(k < 0) {
		tmpk = I.MulS(eta).Add(N.MulS((eta*cosi - math.Sqrt(k))))
	}
	return tmpk
}
