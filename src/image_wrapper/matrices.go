package image_wrapper

type Matrix struct{
	h int
	w int
	data [][]float64
}
type Transformer struct{
	matrix Matrix
}

func CreatePoint(x int, y int) Matrix{
	return Matrix{1,3, [][]float64{{float64(y), float64(x), 1}}}
}

func DotMatrix(a Matrix, b Matrix) Matrix{
	h, w :=  a.h, b.w
	commonDim := a.w // == b.h
	// готовим матрицу с результатом
	c := make([][]float64, h)
	for i := 0; i < h; i++ {
		c[i] = make([]float64, w)
	}
	for i:=0; i<h; i++{
		for j:=0; j<w; j++{
			res := 0.0
			for k:=0; k<commonDim; k++{
				res +=  a.data[i][k] * b.data[k][j]
			}
			c[i][j] = res
		}
	}
	return Matrix{h,w,c}
}

func(matr Matrix)Inverse() Matrix{
	return Matrix{}
}
//func(matr Matrix) Determinant() float64{
//	H, W := matr.h, matr.w
//	sign := 1
//	total := 0
//	for x:=0; x<W;x++{
//		total += sign * matr.data[0][x] * Matrix{2,2,
//			[][]float64{
//				{matr.data[x+1][1], matr.data[x][1]},
//				{}
//			}}
//		sign = sign * (-1)
//	}
//}
//func(matr Matrix) Det2d() float64{
//
//}