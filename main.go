package main

import (
	"fmt"
	"image/color"
	"log"
	"math"

	"github.com/xt0fer/tinyray/engine"
	"github.com/xt0fer/tinyray/vector"
)

const (
	width  = 1024
	height = 768
)

func render(spheres []vector.Sphere, lights []vector.Light) {
	log.Printf("render size(%d)\n", width*height)
	framebuffer := [width * height]vector.Vector{}
	fov := math.Pi / 2.0

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			// framebuffer[i+j*width] = vector.Vector{
			// 	X: float64(j) / float64(height),
			// 	Y: float64(i) / float64(width),
			// 	Z: 0.0,
			// }
			// x := float64((2*(float64(i)+0.5)/width - 1) * math.Tan(fov/2.) * width / height)
			// y := -(2*(float64(j)+0.5)/height - 1) * math.Tan(fov/2.)
			// dir := vector.Vector{X: x, Y: y, Z: -1}.Normalize()
			// framebuffer[i+j*width] = vector.CastRay(vector.Vector{X: 0, Y: 0, Z: 0}, dir, spheres, lights, 0)
			dir_x := (float64(i) + 0.5) - width/2.
			dir_y := -(float64(j) + 0.5) + height/2. // this flips the image at the same time
			dir_z := -height / (2. * math.Tan(fov/2.))
			framebuffer[i+j*width] = vector.CastRay(vector.Vector{X: 0, Y: 0, Z: 0},
				vector.Vector{X: dir_x, Y: dir_y, Z: dir_z}.Normalize(), spheres, lights, 0)
		}
	}

	scene := engine.NewScene(width, height)

	scene.EachPixel(func(x, y int) color.RGBA {
		vec := framebuffer[x+y*width]
		max := math.Max(vec.X, math.Max(vec.Y, vec.Z))
		if max > 1 {
			vec = vec.MulS(1.0 / max)
		}
		return color.RGBA{
			uint8(math.Round(vec.X * 255.0)),
			uint8(math.Round(vec.Y * 255.0)),
			uint8(math.Round(vec.Z * 255.0)),
			255,
		}
	})
	//scene.Save(fmt.Sprintf("./renders/%d.png", time.Now().Unix()))
	scene.Save(fmt.Sprintf("./renders/foo.png"))

}

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
		RefractIdx:   1.0,
		Albedo:       vector.V4{X: 0.6, Y: 0.3, Z: 0.1, A: 0.0},
		DiffuseColor: vector.Vector{X: 0.4, Y: 0.4, Z: 0.3},
		SpecularExp:  50.0,
	}
	redrubber := vector.Material{
		RefractIdx:   1.0,
		Albedo:       vector.V4{X: 0.9, Y: 0.1, Z: 0.0, A: 0.0},
		DiffuseColor: vector.Vector{X: 0.3, Y: 0.1, Z: 0.1},
		SpecularExp:  10.0,
	}
	// mirror(Vec3f(0.0, 10.0, 0.8), Vec3f(1.0, 1.0, 1.0), 1425.);
	mirror := vector.Material{
		RefractIdx:   1.0,
		Albedo:       vector.V4{X: 0.0, Y: 10.0, Z: 0.8, A: 0.0},
		DiffuseColor: vector.Vector{X: 1.0, Y: 1.0, Z: 1.0},
		SpecularExp:  1425.0,
	}
	// glass(1.5, Vec4f(0.0,  0.5, 0.1, 0.8), Vec3f(0.6, 0.7, 0.8),  125.);
	glass := vector.Material{
		RefractIdx:   1.5,
		Albedo:       vector.V4{X: 0.0, Y: 0.5, Z: 0.1, A: 0.8},
		DiffuseColor: vector.Vector{X: 0.6, Y: 0.7, Z: 0.8},
		SpecularExp:  150.0, //1425.0,
	}

	// spheres.push_back(Sphere(Vec3f(-3,    0,   -16), 2,      ivory));
	spheres[3] = vector.Sphere{
		Center:   vector.Vector{X: -3.0, Y: 0.0, Z: -16.0},
		Radius:   2,
		Material: ivory,
	}
	// spheres.push_back(Sphere(Vec3f(-1.0, -2.0, -12), 2, glass));
	spheres[2] = vector.Sphere{
		Center:   vector.Vector{X: -1.0, Y: -1.5, Z: -12.0},
		Radius:   2,
		Material: glass,
	}
	// spheres.push_back(Sphere(Vec3f( 1.5, -0.5, -18), 3, red_rubber));
	spheres[1] = vector.Sphere{
		Center:   vector.Vector{X: 1.5, Y: -0.5, Z: -18.0},
		Radius:   3,
		Material: redrubber,
	}

	// spheres.push_back(Sphere(Vec3f( 7,    5,   -18), 4,      mirror));
	spheres[0] = vector.Sphere{
		Center:   vector.Vector{X: 7.0, Y: 5.0, Z: -18.0},
		Radius:   4,
		Material: mirror,
	}

	render(spheres, lights)
}
