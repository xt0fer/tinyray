package main

import (
	"fmt"

	"github.com/xt0fer/tinyray/vector"
)

const (
	width  = 1024
	height = 768
)

func render() {
	var framebuffer [width * height]vector.Vector

	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			framebuffer[i+j*width] = vector.Vector{
				X: float64(j) / float64(height),
				Y: float64(i) / float64(width),
				Z: 0,
			}
		}
	}
	fmt.Println("render")
}
