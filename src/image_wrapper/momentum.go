package image_wrapper

import (
	"log"
	"math"
)

// См. Лекция 11: Вычисление признаков объектов.
func TwoDimensionalMoment(mask * Matrix, p int, q int) float64{
	h, w := mask.h, mask.w
	total := 0.0
	for y:=0; y<h; y++{
		for x:=0; x<w; x++{
			total += math.Pow(float64(x), float64(p)) *  math.Pow(float64(y), float64(q)) * mask.data[y][x]
		}
	}
	return total
}

func CentralTwoDimensionalMoment(mask * Matrix, p int, q int) float64{
	h, w := mask.h, mask.w
	total := 0.0

	zeroMoment := TwoDimensionalMoment(mask, 0,0)
	xC := TwoDimensionalMoment(mask, 1,0) / zeroMoment
	yC := TwoDimensionalMoment(mask, 0,1) / zeroMoment
	for y:=0; y<h; y++{
		for x:=0; x<w; x++{
			total += math.Pow(float64(x) - xC, float64(p)) *  math.Pow(float64(y) - yC, float64(q)) * mask.data[y][x]
		}
	}
	return total
}

func NormalizedCentralTwoDimensionalMoment(mask * Matrix, p int, q int) float64{
	if p+q<2{
		log.Fatal("Не определено для p+q<2")
	}
	cm00 := CentralTwoDimensionalMoment(mask, 0, 0)
	normalizer := math.Pow(cm00, float64(p+q)/2.0 + 1.0)
	return CentralTwoDimensionalMoment(mask, p, q) / normalizer
}

func Binarize(img Matrix, threshold uint8) Matrix{
	h, w := img.h, img.w
	thresholdFl := float64(threshold)/255.0
	mask := MatrixLikeAnother(&img)
	for y:=0;y<h;y++{
		for x:=0;x<w;x++{
			if img.data[y][x]>thresholdFl{
				mask.data[y][x] = 1.0
			}else {
				mask.data[y][x] = 0.0
			}
		}
	}
	return mask
}


func fill(img * Matrix, labels * Matrix, x int, y int, L float64, W int, H int){
	if (labels.data[y][x] == 0) && (img.data[y][x] == 1) {
		labels.data[y][x] = L
		if y > 0 {
			fill(img, labels, x, y-1, L, W, H)
		}
		if y < H - 1{
			fill(img, labels, x, y+1, L, W, H)
		}
		if x > 0{
			fill(img, labels, x-1, y, L, W, H)
		}
		if x < W - 1 {
			fill(img, labels, x+1, y, L, W, H)

		}
	}
}

func LabelSectors4d(binary Matrix) Matrix{
	L := 1.
	h, w := binary.h, binary.w
	labeledArea := ZeroLikeAnother(&binary)
	for y:=0; y < h; y++{
		for x:=0; x < w; x++{
			fill(&binary, &labeledArea, x, y, L, w, h)
			L += 0.001

		}
	}

	for y:=0; y < h; y++{
		for x:=0; x < w; x++{
			labeledArea.data[y][x] = labeledArea.data[y][x] / L
		}
	}

	return labeledArea
}
