package sameriver

import (
	"github.com/aquilax/go-perlin"
)

func GeneratePerlinTerrain(seed int, w, h int) [][]float64 {
	p := perlin.NewPerlin(2, 2, 3, int64(seed))
	terrain := make([][]float64, w)
	for i := range terrain {
		terrain[i] = make([]float64, h)
	}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			terrain[x][y] = p.Noise2D(float64(x)/10, float64(y)/10)
		}
	}
	return terrain
}
