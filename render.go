package main

import (
	"fmt"
	"image/color"
	"log"
	"time"

	"github.com/xt0fer/tinyray/engine"
	"github.com/xt0fer/tinyray/vector"
)

const (
	width  = 1024
	height = 768
)

func render() {
	log.Printf("render size(%d)\n", width*height)
	framebuffer := [width * height]vector.Vector{}

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			framebuffer[i*height+j] = vector.Vector{
				X: float64(j) / float64(height),
				Y: float64(i) / float64(width),
				Z: 0,
			}
		}
	}

	scene := engine.NewScene(width, height)
	scene.EachPixel(func(x, y int) color.RGBA {
		return color.RGBA{
			uint8(x * 255 / width),
			uint8(y * 255 / height),
			100,
			255,
		}
	})
	log.Println("saving")
	scene.Save(fmt.Sprintf("./renders/%d.png", time.Now().Unix()))

}
