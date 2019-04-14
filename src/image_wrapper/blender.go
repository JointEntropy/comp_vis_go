package image_wrapper

import (
	"image"
	"image/color"
	"log"
	"math"
)

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


func (bl Blender) Blend(mode string) MImageWrapper{
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
	result.Name = "blend" + "_" + mode
	result.data = blend_result
	result.Format = "png"
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
		res = math.Min(1.0, C_b/(1.0-math.Min(C_s, 0.999)))
	}
	return res
}

func colorburn_blend(C_b float64, C_s float64) float64{
	res := 0.0
	if C_s>0.0 {
		res = (1.0 - math.Min(1.0, (1.0-C_b)/math.Max(C_s, 0.001)))
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
	case "difference_blend":
		return difference_blend
	default:
		log.Fatal("Invalid blend mode "+blend_method)
	}
	return normal_blend
}