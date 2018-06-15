package main

const WORLD_HEIGHT = 1024
const WORLD_WIDTH = 1024

const WINDOW_HEIGHT = 640
const WINDOW_WIDTH = 640
const FONTSZ = 24

const FPS = 60
const POINTSZ = 8
const MOVESPEED = 3
const VECLENGTH = 64

const OBSTACLESZ = 32

const DIFFUSION_DIM = 100

var DIFFUSION_CELL_W = float64(WORLD_WIDTH) / float64(DIFFUSION_DIM)
var DIFFUSION_CELL_H = float64(WORLD_HEIGHT) / float64(DIFFUSION_DIM)

var KERN2 = [][]float64{
	[]float64{0.0030, 0.0133, 0.0219, 0.0133, 0.0030},
	[]float64{0.0133, 0.0596, 0.0983, 0.0596, 0.0133},
	[]float64{0.0219, 0.0983, 0.1621, 0.0983, 0.0219},
	[]float64{0.0133, 0.0596, 0.0983, 0.0596, 0.0133},
	[]float64{0.0030, 0.0133, 0.0219, 0.0133, 0.0030},
}
