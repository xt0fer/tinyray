package main

import "github.com/xt0fer/tinyray/vector"

func main() {
	sphere := vector.Sphere{
		Center: vector.Vector{X: -3.0, Y: 0.0, Z: -16.0},
		Radius: 2}
	render(sphere)
}
