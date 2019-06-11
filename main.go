package main

import (
	"./src/image_wrapper"
	"log"
)

func main(){
	parsedData, err := image_wrapper.ParseTiff("images/girl.tif")
	//ok := parser.Parse("images/multipage_tiff_example.tif")
	if err != nil{
		log.Fatal("Error parsing file")
	}

	img := image_wrapper.CreateImageWrapperFromStatic("tiff_parse_res", "jpg", parsedData)
	img.UpdateFromStatic(parsedData)
	matr := img.ToMatrix()

	flRes := image_wrapper.FloydSteinbergMatr(matr, 1)

	flResImg := image_wrapper.FromMatrix(flRes, "parsed_data", "png")
	flResImg.SaveImage("output", "jpg")
}

