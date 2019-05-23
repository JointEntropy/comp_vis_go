package image_wrapper

import (
	"log"
	"math"
)



func ApplyTransformer(tr Transformer, img Matrix) Matrix{
	inverseTransform := tr.matrix.Inverse()
	H, W :=  img.h, img.w
	bounds := Matrix{4, 3, [][]float64{
		{0.0,0.0,1.0},
		{float64(W),0.0,1.0},
		{0.0,float64(H),1.0},
		{float64(W), float64(H),1.0},
	}}
	newBounds := DotMatrix(&bounds, &tr.matrix)

	newH := 0
	for i:=0;i<4;i++{
		newH = int(math.Max(newBounds.data[i][1], float64(newH)))
	}
	newW := 0
	for i:=0;i<4;i++{
		newW = int(math.Max(newBounds.data[i][0], float64(newW)))
	}

	data := make([][]float64, newH)
	for i:=0;i<newH;i++{
		data[i] = make([]float64, newW)
	}
	newImg := CreateMatrix(newH, newW, data)
	for y:=0; y<newH;y++{
		for x:=0; x<newW;x++{
			pnt := Matrix{1,3, [][]float64{{float64(x),float64(y),1}}}
			oldCoords := DotMatrix(&pnt, &inverseTransform).data
			oldX := int(oldCoords[0][0])
			oldY := int(oldCoords[0][1])
			if (oldY >(H-1)) || (oldY<0) || (oldX>(W-1)) || (oldX<0){
				continue
			}
			newImg.data[y][x] = img.data[oldY][oldX]
		}
	}

	return newImg
}

func CreateRotateTransformer(angle float64, mode string) Transformer{
	angle = angle/180.0 * math.Pi
	if mode == "straight"{
		matrix := Matrix{3,3,
						[][]float64{
									{math.Cos(angle), math.Sin(angle),  0},
									{-math.Sin(angle), math.Cos(angle), 0},
									{0,					0,				1}},
		}
		return Transformer{matrix}
	} else if mode == "inverse"{
		matrix := Matrix{3,3,
						[][]float64{
						{math.Cos(angle), -math.Sin(angle), 0},
						{math.Sin(angle), math.Cos(angle), 0},
						{0,				  0,			   1}},
		}
		return Transformer{matrix}
	}
	log.Fatal("Invalid mode for rotate transformer")
	return Transformer{}
}

func CreateShearTransformer(angle float64) Transformer{
	matrix := Matrix{3,3,
		[][]float64{
			{1, 0, 0},
			{math.Tan(angle), 1, 0},
			{0, 0,	1}},
	}
	return Transformer{matrix}
}

func CreateScaleTransformer(x float64, y float64) Transformer{
	matrix := Matrix{3,3,
		[][]float64{
			{x, 0, 0},
			{0, y, 0},
			{0, 0,	1}},
	}
	return Transformer{matrix}
}