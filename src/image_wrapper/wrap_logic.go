package image_wrapper

import (
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
	"log"
	"math"
	"os"
	"strings"
)

type MImageWrapper struct{
	Name string
	Format string
	data draw.Image
}

func LoadMImageWrapperFromString(path string) MImageWrapper{
	img := MImageWrapper{}
	f,err := os.Open(path)
	if err !=nil{
		log.Fatal("Can't load an image by path" + path)
	}

	var tmp_img image.Image
	tmp_img, img.Format, err =  image.Decode(f)
	if err !=nil{
		log.Fatal("Can't decode an image by path" + path)
	}
	tmp_path := strings.Split(path, string(os.PathSeparator))
	fullname := tmp_path [len(tmp_path )-1]
	img.Name = strings.Split(fullname, ".")[0]

	img.data = toGrey(tmp_img)
	return img
}

func toGrey(imgSrc image.Image) draw.Image{
	bounds := imgSrc.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	grayScale := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			grayScale.Set(x, y, rgbaToGrey(imgSrc.At(x,y)))
		}
	}
	return grayScale
}

func (img MImageWrapper) Transpose() draw.Image{
	bounds := img.data.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	grayScale := image.NewGray(image.Rectangle{image.Point{0, 0},
		image.Point{h, w}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			grayScale.Set(y, x, img.data.At(x,y))
		}
	}
	return grayScale
}

func (img *MImageWrapper) UpdateData(new_img draw.Image){
	(*img).data = new_img
}

func (img  MImageWrapper) Mirror(axis uint8){
	bounds := img.data.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	//grayScale := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	if axis==1 {
		for x := 0; x < w; x++ {
			for y := 0; y < h/2; y++ {
				tmp1 := rgbaToGrey(img.data.At(x, y))
				tmp2 := rgbaToGrey(img.data.At(x, h-y))
				img.data.Set(x, h-y, tmp1)
				img.data.Set(x, y, tmp2)

			}
		}
	} else if axis == 0{
		for x := 0; x < w/2; x++ {
			for y := 0; y < h; y++ {
				tmp1 :=  rgbaToGrey(img.data.At(x,y))
				tmp2 := rgbaToGrey(img.data.At(w-x, y))
				img.data.Set(w-x, y, tmp1)
				img.data.Set(x, y,  tmp2)
			}
		}

	}
}


func (img  MImageWrapper) SaveImage(path string, format string) {
	allowed_formats := map[string]bool{"png":true,
		"jpeg":true,
		"jpg":true}
	save_format := img.Format
	if format != "" {
		save_format = format
	}
	if _, ok := allowed_formats[format]; !ok{
		log.Fatal("Can't use format " + format + " to save.")
	}

	fullpath := path + "/" + img.Name + "." + save_format
	f_out, err := os.Create(fullpath)

	if err != nil {
		log.Fatal("Error while creating the file to save on path" + path)
	}


	bounds := img.data.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	save_buffer := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			val := rgbaToGrey(img.data.At(x,y))
			save_buffer.Set(x, y, val)
		}
	}
	if save_format == "png"{
		err = png.Encode(f_out, save_buffer)
	} else if save_format == "jpg"{
		err = jpeg.Encode(f_out, save_buffer, nil)
	}
	if err!= nil{
		log.Fatal("Can't encode an image while saving")
	}

}

func rgbaToGrey(col color.Color) color.Gray{
	rr, gg, bb, _ := col.RGBA()
	r := math.Pow(float64(rr), 2.2)
	g := math.Pow(float64(gg), 2.2)
	b := math.Pow(float64(bb), 2.2)
	m := math.Pow(0.2125*r+0.7154*g+0.0721*b, 1/2.2)
	Y := uint16(m + 0.5)
	return color.Gray{uint8(Y >> 8)}
}


