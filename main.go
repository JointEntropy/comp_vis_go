package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"image/png"
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
	fmt.Println("Load image from path ", path)
	img := MImageWrapper{}
	f,_ := os.Open(path)
	img.data, img.format, _ =  image.Decode(f)

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
	save_format := img.format
	if format != "" {
		save_format = format
	}
	fullpath := path + "/" + img.name + "." + save_format
	f_out, _ := os.Create(fullpath)

	if save_format == "png"{
		png.Encode(f_out, img.data)
	} else if save_format == "jpg"{
		jpeg.Encode(f_out, img.data, nil)
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
	a_b := float64(1)
	for x := 0; x < xdim; x++ {
		for y := 0; y < ydim; y++ {
			c_s := float64(rgbaToGrey(bl.C_s.data.At(x, y)).Y) / 255.0
			c_b := float64(rgbaToGrey(bl.C_b.data.At(x, y)).Y) / 255.0
			a_s:= float64(rgbaToGrey(bl.C_alpha.data.At(x, y)).Y) / 255.0
			val := apply_blend(c_b, c_s, a_s, a_b, mode) * 255.0
			//blend_res := math.Min(c_b, c_s)
			//val := (((1-a_s)*a_b * c_b) + ((1-a_b)*a_s*c_s) + (a_b * a_s * blend_res)) * 255.0
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

func apply_blend(C_b float64, C_s float64, a_s float64, a_b float64, blend_method string) float64{
	blend_res := 1.0
	switch blend_method{
		case "normal_blend":
			blend_res = normal_blend(C_b, C_s)

		case "multiply_blend":
			blend_res = multiply_blend(C_b, C_s)

		case "screen_blend":
			blend_res = screen_blend(C_b, C_s)

		case "darken_blend":
			blend_res = darken_blend(C_b, C_s)

		case "lighten_blend":
			blend_res = lighten_blend(C_b, C_s)

		case "colordodge_blend":
			blend_res = colordodge_blend(C_b, C_s)

		case "colorburn_blend":
			blend_res = colorburn_blend(C_b, C_s)

		case "softlight_blend":
			blend_res = softlight_blend(C_b, C_s)
	}
	return (1.0-a_s)*a_b*C_b + (1.0-a_b)*a_s*C_s + a_s*a_b*blend_res
}



func main(){
	img1 := LoadMImageWrapperFromString("/home/grigory/PycharmProjects/comp_vis/images/cat1.jpg",
										true)
	fmt.Println(img1.name)
	img2 := LoadMImageWrapperFromString("/home/grigory/PycharmProjects/comp_vis/images/cat2.jpg",
										true)
	img3:= LoadMImageWrapperFromString("/home/grigory/PycharmProjects/comp_vis/images/batman.jpg",
		true)

	blender := NewBlender(img1, img2, img3)
	//// blender
	modes := []string{"normal", "multiply", "screen", "darken",
						"lighten", "colordodge", "colorburn", "softlight"}
	for _, mode := range modes{
		img := blender.blend(mode+"_blend")
		img.saveImage("output", "png")
	}
}

