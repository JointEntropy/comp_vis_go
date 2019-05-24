package main

import (
	"./src/image_wrapper"
	"log"
)

func main(){
	parser := image_wrapper.CreateTiffParser()
	ok := parser.Parse("images/girl.tif")
	//ok := parser.Parse("images/multipage_tiff_example.tif")
	if !ok{
		log.Fatal("Error parsing file")
	}
	data := parser.GetData()

	img := image_wrapper.CreateImageWrapperFromStatic("tiff_parse_res", "jpg", data)
	img.UpdateFromStatic(data)
	//img := image_wrapper.FromMatrix(data, "parsed_data", "png")
	img.SaveImage("output", "jpg")
}

