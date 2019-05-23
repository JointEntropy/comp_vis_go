package image_wrapper

func validateBorders(x int, y int, arr * Matrix) bool{
	w, h := arr.w, arr.h
	if (x!=0) && (x!=w-1) && (y!=0) && (y!=h-1){
		return true
	}
	return false
}

func partialDeriveX(x int, y int, arr * Matrix) float64{
	if !validateBorders(x,y, arr){
		return 0.0
	}
	return arr.data[y][x+1] - arr.data[y][x-1]
}

func partialDeriveY(x int, y int, arr * Matrix) float64{
	if !validateBorders(x,y, arr){
		return 0.0
	}
	return arr.data[y+1][x] - arr.data[y-1][x]
}

func partialDeriveXY(x int, y int, arr * Matrix) float64{
	if !validateBorders(x,y, arr){
		return 0.0
	}
	return arr.data[y+1][x+1] - arr.data[y-1][x-1]
}

type partialDeriveInterface func(int, int, *Matrix) float64

func applyPartialDerive(arr * Matrix, method partialDeriveInterface) Matrix{
	w, h := arr.w, arr.h
	res := MatrixLikeAnother(arr)
	for y:=0; y<h; y++{
		for x:=0; x<w; x++{
			res.data[y][x] = method(x, y, arr)
		}
	}
	return res
}

func BicubicInterpolation(img * Matrix) Matrix{
	f := img.data
	fX := applyPartialDerive(img, partialDeriveX).data
	fY := applyPartialDerive(img, partialDeriveX).data
	fXY := applyPartialDerive(img, partialDeriveXY).data

	aLeft := CreateMatrix(4,4, [][]float64{
		{1.0, 0.0, 0, 0.0},
		{0.0, 0.0, 1, 0.0},
		{-3.0, 3.0, -2.0, -1.0},
		{2.0, -2.0, 1.0, 1.0},
	})
	aRight := aLeft.Transpose()

	h, w  := img.h, img.w
	res := MatrixLikeAnother(img)
	for y:=0; y<h-1; y++{
		for x:=0;x<w-1;x++{
			fDer := CreateMatrix(4,4, [][]float64{
				{f[0+y][0+x], f[0+y][1+x], fY[0+y][0+x], fY[0+y][1+x]},
				{f[1+y][0+x], f[1+y][1+x], fY[1+y][0+x], fY[1+y][1+x]},
				{fX[0+y][0+x], fX[0+y][1+x], fXY[0+y][0+x], fXY[0+y][1+x]},
				{fX[1+y][0+x], fX[1+y][1+x], fXY[1+y][0+x], fXY[1+y][1+x]},
			})
			tmp := DotMatrix(&aLeft, &fDer)
			a := DotMatrix(&tmp, &aRight)
			x_ := float64(x)/float64(w)
			y_ := float64(y)/float64(h)
			left := CreateMatrix(1,4, [][]float64{
				{1, x_, x_*x_, x_*x_*x_},
			})
			right := CreateMatrix(4,1, [][]float64{
				{1},
				{y_},
				{y_*y_},
				{y_*y_*y_},
			})
			tmp = DotMatrix(&left, &a)
			res.data[y][x] = DotMatrix(&tmp, &right).data[0][0]
		}
	}
	for x:=0;x<w;x++{
		res.data[0][x] = img.data[0][x]
		res.data[h-1][x] = img.data[h-1][x]
	}
	for y:=0;y<h;y++{
		res.data[y][0] = img.data[y][0]
		res.data[y][w-1] = img.data[y][w-1]
	}
	return  res
}