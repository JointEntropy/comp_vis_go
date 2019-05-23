package main

import (
	"./src/image_wrapper"
	"fmt"
)

func main(){
	img1 := image_wrapper.LoadMImageWrapperFromString("images/lena_small.jpg")
	img1.SaveImage("output", "png")
	imgMatr := img1.ToMatrix()

	// Реализовать поворот изображения на произвольный угол с использованием бикубической интерполяции.
	rotateTransformer := image_wrapper.CreateRotateTransformer(20.0, "straight")
	rotateRes:= image_wrapper.ApplyTransformer(rotateTransformer, imgMatr)
	res := image_wrapper.FromMatrix(rotateRes, "rotate_result", "png")
	res.SaveImage("output", "png")

	//// Реализация интерполции
	//resMatr := image_wrapper.BicubicInterpolation(&imgMatr)
	//res = image_wrapper.FromMatrix(resMatr, "interpolation_result", "png")
	//res.SaveImage("output", "jpg")
	fmt.Println("Fine")

}

