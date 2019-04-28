package image_wrapper

import (
	"image"
	"image/draw"
	"log"
	"math"
)

func ApplyTransformer(tr Transformer, imgWrapper MImageWrapper) draw.Image{
	img := imgWrapper.data
	inverseTransform := tr.matrix.Inverse()
	imgBounds := img.Bounds()
	H, W :=  imgBounds .Max.Y, imgBounds.Max.X
	bounds := Matrix{4, 3, [][]float64{
		{0,0,1},
		{float64(W),0,1},
		{0,float64(H),1},
		{float64(W), float64(H),1},
	}}
	newBounds := DotMatrix(bounds, tr.matrix)

	newH := 0
	for i:=0;i<4;i++{
		newH = int(math.Max(newBounds.data[i][0], float64(newH)))
	}
	newW := 0
	for i:=0;i<4;i++{
		newW = int(math.Max(newBounds.data[i][1], float64(newW)))
	}
	newImg := image.NewGray(image.Rectangle{image.Point{0, 0},
			image.Point{newH, newW}})
	for y:=0; y<newH;y++{
		for x:=0; x<newW;x++{
			pnt := Matrix{1,3, [][]float64{{float64(x),float64(y),1}}}
			oldCoords := DotMatrix(pnt, inverseTransform).data
			oldX := int(oldCoords[0][0])
			oldY := int(oldCoords[0][1])
			if (oldY >(H-1)) || (oldY<0) || (oldX>(W-1)) || (oldX<0){
				continue
			}
			val := rgbaToGrey(img.At(oldX, oldY))
			newImg.Set(x,y, val)
		}
	}
	return newImg
}

func CreateRotateTransformer(angle float64, mode string) Transformer{
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