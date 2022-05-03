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

func render() {
	log.Printf("render size(%d)\n", width*height)
	framebuffer := [height * width]vector.Vector{}

	for i := 0; i < width; i++ {
		for j := 0; j < height; j++ {
			framebuffer[i*height+j] = vector.Vector{
				X: float64(j) / float64(width),
				Y: float64(i) / float64(height),
				Z: 0,
			}
			// if framebuffer[i*height+j].X > 1.0 ||  framebuffer[i*height+j].Y > 1.0{
			// 	log.Println("GREATER THAN 1", i, j)
			// }
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
