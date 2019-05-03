package main

import (
	"./src/image_wrapper"
	"fmt"
)

func main(){
	img1 := image_wrapper.LoadMImageWrapperFromString("images/girl.jpg")
	//// Реализовать поворот изображения на произвольный угол с использованием бикубической интерполяции.
	rotateTransformer := image_wrapper.CreateRotateTransformer(20.0, "straight")
	resMatr := image_wrapper.ApplyTransformer(rotateTransformer, img1)
	img1.UpdateFromStatic(resMatr)
	img1.SaveImage("output", "png")
	fmt.Println("Fine")

}

