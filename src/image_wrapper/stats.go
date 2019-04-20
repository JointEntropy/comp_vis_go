package image_wrapper

import (
	"image"
	"image/color"
	"math"
)

const eps float64 = 1e-4

type Analyzer struct {
	image MImageWrapper
	HistData map[uint8]uint32
	ProbsData map[uint8]float64
}

func NewAnalyzer(img MImageWrapper) Analyzer{
	hist := GetHist(img.data)
	probs := GetProbs(img.data, hist)
	return Analyzer{img, hist, probs}
}

func GetHist(img image.Image) map[uint8]uint32{
	res := make(map[uint8]uint32)
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			orig := img.At(x,y)
			val := rgbaToGrey(orig).Y
			if count, ok := res[val]; ok {
				res[val] = count + 1
			}else{
				res[val] = 1
			}
		}
	}
	return res
}

func GetProbs(img image.Image, hist map[uint8] uint32) map[uint8]float64{
	res := make(map[uint8]float64)
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	N := float64(w*h)
	for k, v := range hist {
		res[k] = float64(v)/N  //res.Set(k.b, k.a, color.Gray{v})
	}
	return res
}

func CalcMean(img image.Image, probs map[uint8]float64) float64{
	var sum float64 = 0
	for i, prob := range probs {
		sum += float64(i)*(prob + eps)
	}
	return sum
}

func CalcVar(img image.Image, probs map[uint8]float64)float64{
	E := CalcMean(img, probs)
	var sigma2 float64 = 0
	for i, prob := range probs {
		sigma2 += math.Pow(float64(i) - E, 2)*prob
	}
	return sigma2
}

func CalcKurtosis(img image.Image, probs map[uint8]float64)float64{
	V := CalcVar(img, probs)
	E := CalcMean(img, probs)
	kurt := (1/(math.Pow(V,2))+eps)
	var sum float64 = 0
	for i, prob := range probs {
		sum += math.Pow(float64(i) - E, 4)*prob
	}
	kurt = kurt*sum - 3
	return kurt
}
func CalcSkewness(img image.Image, probs map[uint8]float64)float64{
	V := CalcVar(img, probs)
	E := CalcMean(img, probs)
	kurt :=  (1/(math.Pow(V,1.5)))
	var sum float64 = 0
	for i, prob := range probs {
		sum += math.Pow(float64(i) - E, 3)*prob
	}
	kurt = kurt*sum
	return kurt
}

func CalcEnergy(img image.Image, probs map[uint8]float64)float64{
	var energy float64 = 0
	for _, prob := range probs {
		energy += math.Pow(prob, 2)
	}
	return energy
}
func CalcEntropy(img image.Image, probs map[uint8]float64)float64{
	var entropy float64 = 0
	for _, prob := range probs {
		entropy += prob *  math.Log2(prob+eps)
	}
	return -entropy
}

type pair struct {
	a int
	b int
}
func CoocurenceMatrix(img image.Image, d_r int, d_c int) image.Image{
	pr := make(map[pair]int)
	bounds := img.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y

	var max_val int = 0
	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			if (x+d_r >= w) || (y+d_c >=h){
				continue
			}
			val1 := int(rgbaToGrey(img.At(x,y)).Y)
			val2 :=int(rgbaToGrey(img.At(x+d_r, y+d_c)).Y)
			key := pair{val1, val2}
			if val, ok := pr[key]; ok {
				pr[key] = val + 1
			}else{
				pr[key] = 1
			}
			max_val = int(math.Max(float64(max_val), math.Max(float64(val1), float64(val2))))
		}
	}
	res := image.NewGray(image.Rectangle{image.Point{0, 0},
						image.Point{max_val + 1, max_val + 1}})
	for k, v := range pr {
		res.Set(k.b, k.a, color.Gray{uint8(v)})
	}
	return res
}