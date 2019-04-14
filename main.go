package main

import (
	"./src/image_wrapper"
)

func main(){
	img1 := image_wrapper.LoadMImageWrapperFromString("images/lena.jpg")
	img2 := image_wrapper.LoadMImageWrapperFromString("images/cat1.jpg")
	img3:= image_wrapper.LoadMImageWrapperFromString("images/batman.jpg")

	// Blend mode testing
	blender := image_wrapper.NewBlender(img1, img2, img3)
	modes := []string{"normal", "multiply", "screen", "darken", "difference",
						"lighten", "colordodge", "colorburn", "softlight"}
	for _, mode := range modes{
		img := blender.Blend(mode+"_blend")
		img.SaveImage("output", "png")
	}
	// Image affine test
	img1.Mirror(1)
	img1.Name = "mirrored_ver"
	img1.SaveImage("output", "png")

	img2.Mirror(0)
	img2.Name = "mirrored_hor"
	img2.SaveImage("output", "png")

	img3.Transpose()
	img3.Name = "transposed"
	img3.SaveImage("output","png")
}

