package main

import "github.com/xt0fer/tinyray/vector"

func main() {
	spheres := make([]vector.Sphere, 4)
	lights := make([]vector.Light, 3)

	lights[2] = vector.Light{
		Position: vector.Vector{
			X: -20, Y: 20, Z: 20},
		Intensity: 1.5,
	}

	lights[1] = vector.Light{
		Position: vector.Vector{
			X: 30, Y: 50, Z: -25},
		Intensity: 1.8,
	}

	lights[0] = vector.Light{
		Position: vector.Vector{
			X: 30, Y: 20, Z: 30},
		Intensity: 1.7,
	}

	ivory := vector.Material{
		Albedo:       [2]float64{0.6, 0.3},
		DiffuseColor: vector.Vector{X: 0.4, Y: 0.4, Z: 0.3},
		SpecularExp:  50.0,
	}
	redrubber := vector.Material{
		Albedo:       [2]float64{0.9, 0.1},
		DiffuseColor: vector.Vector{X: 0.3, Y: 0.1, Z: 0.1},
		SpecularExp:  10.0,
	}

	// spheres.push_back(Sphere(Vec3f(-3,    0,   -16), 2,      ivory));
	spheres[3] = vector.Sphere{
		Center:   vector.Vector{X: -3.0, Y: 0.0, Z: -16.0},
		Radius:   2,
		Material: ivory,
	}
	// spheres.push_back(Sphere(Vec3f(-1.0, -1.5, -12), 2, red_rubber));
	spheres[2] = vector.Sphere{
		Center:   vector.Vector{X: -1.0, Y: -1.5, Z: -12.0},
		Radius:   2,
		Material: redrubber,
	}
	// spheres.push_back(Sphere(Vec3f( 1.5, -0.5, -18), 3, red_rubber));
	spheres[1] = vector.Sphere{
		Center:   vector.Vector{X: -1.5, Y: -0.5, Z: -18.0},
		Radius:   3,
		Material: redrubber,
	}

	// spheres.push_back(Sphere(Vec3f( 7,    5,   -18), 4,      ivory));
	spheres[0] = vector.Sphere{
		Center:   vector.Vector{X: 7.0, Y: 5.0, Z: -18.0},
		Radius:   4,
		Material: ivory,
	}

	render(spheres, lights)
}
