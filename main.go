package main

import (
	"./src/image_wrapper"
)

func main(){
	img1 := image_wrapper.LoadMImageWrapperFromString("images/girl.jpg")
	// Реализовать поворот изображения на произвольный угол с использованием бикубической интерполяции.
	rotateTransformer := image_wrapper.CreateRotateTransformer(20.0, "straight")
	resImg := image_wrapper.MImageWrapper{"", "jpg",
		image_wrapper.ApplyTransformer(rotateTransformer, )}
	resImg.SaveImage("output", "png")
}

