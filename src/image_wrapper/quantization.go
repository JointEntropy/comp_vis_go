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
	ln := int(math.Round(math.Pow(2.0, float64(k))))
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
	return float64(closest)
}



func FloydSteinberg(image * MImageWrapper, k uint8) Matrix{
	pallet := NewPallet(k)
	imgMatrix := image.ToMatrix()
	w, h := imgMatrix.w, imgMatrix.h
	for y:=0;y<h;y++{
		for x:=0;x<w;x++{
			pixel := imgMatrix.data[y][x]
			value := pallet.FindClosest(pixel)
			err_ := pixel - value
			imgMatrix.data[y][x] = value



			if y==h-1{
				if x<w-1{
					imgMatrix.data[y][x + 1] += err_
				}
			}else{
				if x==w-1{
					imgMatrix.data[y + 1][x] += err_
					continue
				}
				imgMatrix.data[y + 0][x + 1] += err_ * 7.0/16.0
				if x>0{
					imgMatrix.data[y + 1][x - 1] += err_ * 3.0/16.0
				}
				imgMatrix.data[y + 1][x    ] += err_ * 5.0/16.0
				imgMatrix.data[y + 1][x + 1] += err_ * 1.0/16.0


			}

		}
	}

	//for i:=0;i<h;i++{
	//	res.data[i][w-1] = 1.0
	//}
	//for i:=0;i<w;i++{
	//	res.data[h-1][i] = 1.0
	//}
	return imgMatrix

}
