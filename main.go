package main

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
	name string
	format string
	data image.Image
}

func LoadMImageWrapperFromString(path string, convertToGrey bool) MImageWrapper{
	img := MImageWrapper{}
	f,err := os.Open(path)
	if err !=nil{
		log.Fatal("Can't load an image by path" + path)
	}
	img.data, img.format, err =  image.Decode(f)
	if err !=nil{
		log.Fatal("Can't decode an image by path" + path)
	}
	tmp_path := strings.Split(path, string(os.PathSeparator))
	fullname := tmp_path [len(tmp_path )-1]
	img.name = strings.Split(fullname, ".")[0]

	if convertToGrey {
		tmp_data := toGrey(img.data)
		img.data = tmp_data
	}
	return img
}

func toGrey(imgSrc image.Image) image.Image{
	bounds := imgSrc.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	grayScale := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			grayColor := rgbaToGrey(imgSrc.At(x,y))
			grayScale.Set(x, y, grayColor)
		}
	}
	return grayScale
}

func (img MImageWrapper) transpose() image.Image{
	bounds := img.data.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	grayScale := image.NewGray(image.Rectangle{image.Point{0, 0},
											image.Point{h, w}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			grayColor := rgbaToGrey(img.data.At(x,y))
			grayScale.Set(y, x, grayColor)
		}
	}
	return grayScale
}


func (img MImageWrapper) mirror(axis uint8) image.Image{
	bounds := img.data.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	grayScale := image.NewGray(image.Rectangle{image.Point{0, 0}, image.Point{w, h}})
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			grayColor := rgbaToGrey(img.data.At(x,y))
			if axis == 0{
				grayScale.Set(w-x, y, grayColor)
			} else if axis == 1{
				grayScale.Set(x, h-y, grayColor)
			}
		}
	}
	return grayScale
}


func (img MImageWrapper) saveImage(path string, format string) {
	allowed_formats := map[string]bool{"png":true,
								"jpeg":true,
								"jpg":true}
	save_format := img.format
	if format != "" {
		save_format = format
	}
	if _, ok := allowed_formats[format]; !ok{
		log.Fatal("Can't use format " + format + " to save.")
	}

	fullpath := path + "/" + img.name + "." + save_format
	f_out, err := os.Create(fullpath)

	if err != nil {
		log.Fatal("Error while creating the file to save on path" + path)
	}

	if save_format == "png"{
		err = png.Encode(f_out, img.data)
	} else if save_format == "jpg"{
		err = jpeg.Encode(f_out, img.data, nil)
	}
	if err!= nil{
		log.Fatal("Can't encode an image while saving")
	}

}
func (img MImageWrapper) toMatrix(){
	b := img.data.Bounds()
	m := image.NewRGBA(image.Rect(0, 0, b.Dx(), b.Dy()))
	draw.Draw(m, m.Bounds(), img.data, b.Min, draw.Src)
}

type Blender struct{
	C_s MImageWrapper
	C_b  MImageWrapper
	C_alpha MImageWrapper
}


// Как получить матрицу чиселок из прочитанного файла. См.:
// https://stackoverflow.com/questions/33186783/get-a-pixel-array-from-from-golang-image-image

func NewBlender(img_a MImageWrapper, img_b MImageWrapper, img_alpha MImageWrapper) Blender{
	blend := Blender{}
	blend.C_s = img_a
	blend.C_b = img_b
	blend.C_alpha = img_alpha
	// align sizes
	return blend
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


func (bl Blender) blend(mode string) MImageWrapper{
	b := bl.C_b.data.Bounds()
	s := bl.C_s.data.Bounds()

	xdim := int(math.Min(float64(b.Dx()), float64(s.Dx())))
	ydim := int(math.Min(float64(b.Dy()), float64(s.Dy())))

	upLeft := image.Point{0, 0}
	lowRight := image.Point{xdim, ydim}
	blend_result := image.NewGray(image.Rectangle{upLeft, lowRight})
	// Set color for each pixel.
	blend_method := select_blend(mode)


	a_b := float64(1)
	for x := 0; x < xdim; x++ {
		for y := 0; y < ydim; y++ {
			c_s := float64(rgbaToGrey(bl.C_s.data.At(x, y)).Y) / 255.0
			c_b := float64(rgbaToGrey(bl.C_b.data.At(x, y)).Y) / 255.0
			a_s:= float64(rgbaToGrey(bl.C_alpha.data.At(x, y)).Y) / 255.0
			val := blend_method(c_b, c_s)
			val = (((1-a_s)*a_b * c_b) + ((1-a_b)*a_s*c_s) + (a_b * a_s * val)) * 255.0
			c_res := color.Gray{uint8(val)}
			blend_result.Set(x, y, c_res)
		}
	}
	result := MImageWrapper{}
	result.name = "blend" + "_" + mode
	result.data = blend_result
	result.format = "png"
	return result
}


func normal_blend(C_b float64, C_s float64) float64{
	return C_s
}

func multiply_blend(C_b float64, C_s float64) float64{
	return C_b*C_s
}

func screen_blend(C_b float64, C_s float64) float64{
	return 1 - (1 - C_b)*(1 - C_s)
}

func darken_blend(C_b float64, C_s float64) float64{
	return math.Min(C_b, C_s)

}
func lighten_blend(C_b float64, C_s float64) float64{
	return math.Max(C_b, C_s)
}

func difference_blend(C_b float64, C_s float64) float64{
	return math.Abs(C_b - C_s)
}

func colordodge_blend(C_b float64, C_s float64) float64{
	res := 1.0
	if C_s < 1 {
		res = math.Min(1.0, C_b/(1.0-C_s))
	}
	return res
}

func colorburn_blend(C_b float64, C_s float64) float64{
	res := 0.0
	if C_s>0 {
		res = (1.0 - math.Min(1.0, (1.0-C_b)/C_s))
	}
	return res
}

func softlight_blend(C_b float64, C_s float64) float64{
	res := C_b - (1-2*C_s)*C_b*(1-C_b)
	if C_s > 0.5{
		res = (C_b + (2*C_s - 1) * (D_x(C_b) - C_b))
	}
	return res
}
func D_x(x float64) float64{
	res := ((16*x - 12)*x + 4)*x
	if x>0.25{
		res =  math.Pow(x,0.5)
	}
	return res
}

func select_blend(blend_method string) func(float64, float64) float64{
	switch blend_method{
		case "normal_blend":
			return normal_blend
		case "multiply_blend":
			return multiply_blend
		case "screen_blend":
			return screen_blend
		case "darken_blend":
			return darken_blend
		case "lighten_blend":
			return lighten_blend
		case "colordodge_blend":
			return colordodge_blend
		case "colorburn_blend":
			return colorburn_blend
		case "softlight_blend":
			return softlight_blend
		default:
			log.Fatal("Invalid blend mode "+blend_method)
	}
	return normal_blend
}

func main(){
	img1 := LoadMImageWrapperFromString("images/lena.jpg",
										true)
	img2 := LoadMImageWrapperFromString("images/cat1.jpg",
										true)
	img3:= LoadMImageWrapperFromString("images/batman.jpg",
		true)

	// Blend mode testing
	blender := NewBlender(img1, img2, img3)
	modes := []string{"normal", "multiply", "screen", "darken",
						"lighten", "colordodge", "colorburn", "softlight"}
	for _, mode := range modes{
		img := blender.blend(mode+"_blend")
		img.saveImage("output", "png")
	}
	// Image affine test
	img1.data = img1.mirror(0)
	img1.name = "mirrored_hor"
	img1.saveImage("output", "png")

	img2.data = img2.mirror(1)
	img2.name = "mirrored_ver"
	img2.saveImage("output", "png")

	img3.data = img3.transpose()
	img3.name = "transposed"
	img3.saveImage("output","png")
}

