package image_wrapper

import (
	"image"
	"image/draw"
)

// Legacy. Трансфрмаци без использования афинных преобразвований.
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

//
//func DotMatrix(a draw.Image, b draw.Image) draw.Image {
//	a_bounds := a.Bounds()
//	b_bounds := b.Bounds()
//	h, w :=  a_bounds.Max.Y, b_bounds.Max.X
//	common_dim := a_bounds.Max.X // == b_bounds.Max.Y
//	c := image.NewGray(image.Rectangle{image.Point{0, 0},
//		image.Point{h, w}})
//	for i:=0; i<h; i++{
//		for j:=0; j<w; j++{
//			res := 0.0
//			for k:=0; k<common_dim; k++{
//				a_dot := float64(rgbaToGrey(a.At(i,k)).Y) / 255.0
//				b_dot := float64(rgbaToGrey(b.At(k, j)).Y) / 255.0
//				res +=  a_dot * b_dot
//			}
//			c.Set(i, j, color.Gray{uint8(res * 255)})
//		}
//	}
//	return c
//}

