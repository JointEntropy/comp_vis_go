package main

import (
	"./src/image_wrapper"
)

func main(){
	img1 := image_wrapper.LoadMImageWrapperFromString("images/cat1.jpg")
	img1.SaveImage("output", "jpg")
	imgMatr := img1.ToMatrix()

	// сумма равна 0, но в ядре есть и отрицатеьные и положительные элеменьы
	highPassFilter := image_wrapper.CreateMatrix(3,3,[][]float64{
		{-0.125,-0.125, -0.125},
		{-0.125, 1.0, -0.125},
		{-0.125,-0.125,-0.125},
	})
	// сумма равна единице, но все элемены неотрицательны
	lowPassFilter := image_wrapper.CreateMatrix(3,3,[][]float64{
		{1/9.,1/9.,1/9.},
		{1/9.,1/9.,1/9.},
		{1/9.,1/9.,1/9.},
	})

	// сумма ядра равна 1
	sharpeningFilter := image_wrapper.CreateMatrix(3,3,[][]float64{
		{0.0,		-1.0/6.0, 	0.0},
		{-1.0/6.0, 	10.0/6.0, 		-1.0/6.0},
		{0.0,		-1.0/6.0,	0.0},
	})
	convolveRes := image_wrapper.Convolve2d(imgMatr, highPassFilter)
	convolveRes.DumpToFile("output/highPassFilterRes.txt")
	image_wrapper.FromMatrix(convolveRes, "highPassRes", "jpg").SaveImage("output", "jpg")

	convolveRes = image_wrapper.Convolve2d(imgMatr, lowPassFilter)
	convolveRes.DumpToFile("output/lowPassFilterRes.txt")
	image_wrapper.FromMatrix(convolveRes, "lowPassRes", "jpg").SaveImage("output", "jpg")

	convolveRes = image_wrapper.Convolve2d(imgMatr, sharpeningFilter)
	convolveRes.DumpToFile("output/sharpeningFilterRes.txt")
	image_wrapper.FromMatrix(convolveRes, "sharpeningRes", "jpg").SaveImage("output", "jpg")
}

