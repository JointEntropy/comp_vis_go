package main

import (
	"./src/image_wrapper"
)

func main(){
	img1 := image_wrapper.LoadMImageWrapperFromString("images/girl.jpg")
	//img_matr := img1.ToMatrix()

	///res := image_wrapper.FromMatrix(img_matr, "test_matr_repr", "jpg")
	res := image_wrapper.FloydSteinberg(img1, 1)
	res.DumpToFile("result.txt")
	//res.SaveImage("output", "jpg")
}

