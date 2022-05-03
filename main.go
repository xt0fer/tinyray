package main

import "github.com/xt0fer/tinyray/vector"

func main() {
	spheres := make([]vector.Sphere, 4)
	lights := make([]vector.Light, 1)

	lights[0] = vector.Light{
		Position:  vector.Vector{
			X: -20, Y: 20, Z: 20},
		Intensity: 1.5,
	}

	ivory := vector.Material{
		DiffuseColor: vector.Vector{X: 0.4, Y: 0.4, Z: 0.3},
	}
	redrubber := vector.Material{
		DiffuseColor: vector.Vector{X: 0.3, Y: 0.1, Z: 0.1},
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
