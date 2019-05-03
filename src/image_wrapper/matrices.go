package image_wrapper

import (
	"fmt"
	"os"
)

type Matrix struct{
	h int
	w int
	data [][]float64
}
type Transformer struct{
	matrix Matrix
}

func CreateMatrix(h int, w int, data [][]float64) Matrix{
	return Matrix{h,w,data}
}

func CreatePoint(x int, y int) Matrix{
	return Matrix{1,3, [][]float64{{float64(y), float64(x), 1}}}
}

func DotMatrix(a  * Matrix, b  * Matrix) Matrix{
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

func (matr * Matrix) DumpToFile(filepath string){
	w, h := matr.w, matr.h
	f, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i:=0;i<h;i++{
		for j:=0;j<w;j++{
			st := fmt.Sprintf(" %f", matr.data[i][j]) // s == "123.456000"
			f.WriteString(st)
		}
		f.WriteString("\n")
	}

	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
func(matr Matrix)Inverse() Matrix{
	return Matrix{}
}


func  MatrixLikeAnother(matr * Matrix) Matrix{
	data := make([][]float64, matr.h)
	for i:=0;i<matr.h;i++{
		data[i] = make([]float64, matr.w)
	}
	return CreateMatrix(matr.h, matr.w, data)

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