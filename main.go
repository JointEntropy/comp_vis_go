package main

import (
	"./src/image_wrapper"
)

func main(){
	img1 := image_wrapper.LoadMImageWrapperFromString("images/parrot_noise.jpg")
	img1.SaveImage("output", "jpg")
	imgMatr := img1.ToMatrix()

	// сумма равна 0, но в ядре есть и отрицатеьные и положительные элеменьы
	medianFilter := image_wrapper.CreateMatrix(3,3,[][]float64{
		{0, 1, 0},
		{1, 1, 1},
		{0, 1, 0},
	})
	k := 3
	convolveRes := image_wrapper.RangFilter(imgMatr, medianFilter, k)
	convolveRes.DumpToFile("output/medianFilter.txt")
	image_wrapper.FromMatrix(convolveRes, "medianFilter", "jpg").SaveImage("output", "jpg")

}

