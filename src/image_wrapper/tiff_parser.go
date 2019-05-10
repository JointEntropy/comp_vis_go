package image_wrapper

import (
	"image"
	"io"
)

type Tiff struct{
	reader io.Reader
	data image.Image

}

func CreateTiffParser() Tiff{
	return Tiff{}
}

func (tiff * Tiff) Parse(path string) bool{

	data, err := parseTiff(path)
	if err != nil{
		return false
	}
	tiff.data = data
	return true
}

func (tiff * Tiff) GetData() image.Image{
	return tiff.data

}