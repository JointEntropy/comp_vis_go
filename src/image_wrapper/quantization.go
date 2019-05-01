package image_wrapper

import (
	"log"
	"math"
)

type Pallet struct {
	k uint8
	pallet []float64
}

func NewPallet(k uint8) Pallet{
	ln := int(math.Round(math.Pow(2, float64(k))))
	var pallet []float64 = make([]float64, ln)
	for i:=0;i<ln;i++{
		pallet[i] = float64(i)/float64(ln-1)
	}
	return Pallet{k, pallet}
}

func (pal Pallet)FindClosest(val float64) float64 {
	minDelta := 1.0
	closest := -1.0
	for _, element := range pal.pallet{
		delta := math.Abs(val - element)
		if delta < minDelta{
			minDelta = delta
			closest = element
		}
	}
	if closest == -1.0{
		log.Fatal("Huinya kakaya to")
	}
	return closest
}



func FloydSteinberg(image MImageWrapper, k uint8) Matrix{
	pallet := NewPallet(k)
	imgMatrix := image.ToMatrix()
	w, h := imgMatrix.w, imgMatrix.h
	for y:=0;y<h-1;y++{
		for x:=1;x<w-1;x++{
			pixel := imgMatrix.data[y][x]
			value := pallet.FindClosest(pixel)
			err := pixel - value
			imgMatrix.data[y + 0][x + 1] += err * 7/16
			imgMatrix.data[y + 1][x - 1] += err * 3/16
			imgMatrix.data[y + 1][x    ] += err * 5/16
			imgMatrix.data[y + 1][x + 1] += err * 1/16
		}
	}

	for i:=0;i<h;i++{
		imgMatrix.data[i][w-1] = 1
	}
	for i:=0;i<w;i++{
		imgMatrix.data[h-1][i] = 1
	}
	return imgMatrix

}
