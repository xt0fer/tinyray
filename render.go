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
	width  = 768
	height = 768
)

func render(sphere vector.Sphere) {
	log.Printf("render size(%d)\n", width*height)
	framebuffer := [height * width]vector.Vector{}
	fov := math.Pi / 2.0

	for j := 0; j < height; j++ {
		for i := 0; i < width; i++ {
			// framebuffer[i*height+j] = vector.Vector{
			// 	X: float64(j) / float64(width),
			// 	Y: float64(i) / float64(height),
			// 	Z: 0,
			// }
			x := float64((2*(float64(i)+0.5)/width - 1) * math.Tan(fov/2.) * width / height)
			y := -(2*(float64(j)+0.5)/height - 1) * math.Tan(fov/2.)
			dir := vector.Vector{x, y, -1}.Normalize()
			framebuffer[i+j*width] = vector.CastRay(vector.Vector{0, 0, 0}, dir, sphere)
		}
	}

	scene := engine.NewScene(width, height)

	scene.EachPixel(func(x, y int) color.RGBA {
		vec := framebuffer[x*height+y]
		return color.RGBA{
			uint8(math.Round(vec.X * 255.0)),
			uint8(math.Round(vec.Y * 255.0)),
			uint8(math.Round(vec.Z * 255.0)),
			255,
		}
	})
	log.Println("saving")
	//scene.Save(fmt.Sprintf("./renders/%d.png", time.Now().Unix()))
	scene.Save(fmt.Sprintf("./renders/foo.png"))

}
