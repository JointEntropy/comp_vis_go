package main

import (
	"./src/image_wrapper"
)

func main(){
	img1 := image_wrapper.LoadMImageWrapperFromString("images/cat1.jpg")
	img1.SaveImage( "output", "jpg")

	///res := image_wrapper.FromMatrix(img_matr, "test_matr_repr", "jpg")
	res := image_wrapper.FloydSteinberg(&img1, 2)
	res.DumpToFile("output/FloydSteinbergRes.txt")
	image_wrapper.FromMatrix(res, "FloydSteinbergRes", "jpg").SaveImage("output", "jpg")
}

