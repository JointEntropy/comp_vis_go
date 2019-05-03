package image_wrapper

import (
	"fmt"
	"log"
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
func(aMatr * Matrix)Inverse() Matrix{
	det := aMatr.Determinant()
	if det == 0{
		log.Fatal("Определитель равен нулю. Сорян :-(")
	}
	res := MatrixLikeAnother(aMatr)
	r := res.data
	a := aMatr.data
	// Ищем матрицу алгебраических дополнений
	r[0][0] = a[1][1]*a[2][2] - a[2][1]*a[1][2]
	r[0][1] = - (a[1][0]*a[2][2] - a[2][0]*a[1][2])
	r[0][2] = a[1][0]*a[2][1] - a[2][0]*a[1][1]

	r[1][0] = - (a[0][1]*a[2][2] - a[2][1]*a[0][2])
	r[1][1] = a[0][0]*a[2][2] - a[2][0]*a[0][2]
	r[1][2] = - (a[0][0]*a[2][1] - a[2][0]*a[0][1])

	r[2][0] = a[0][1]*a[1][2] - a[1][1]*a[0][2]
	r[2][1] = - (a[0][0]*a[1][2] - a[1][0]*a[0][2])
	r[2][2] = a[0][0]*a[1][1] - a[1][0]*a[0][1]

	res.data = r
	resTransposed := res.Transpose()
	result := resTransposed.ScalarMultiply(1.0/det)
	return result
}
func (a * Matrix) ScalarMultiply(val float64) Matrix{
	res := MatrixLikeAnother(a)
	for y:=0;y<res.h;y++{
		for x:=0;x<res.w; x++{
			res.data[y][x] = val * a.data[y][x]
		}
	}
	return res
}

func (a * Matrix)Transpose() Matrix{
	resData := make([][]float64, a.w)
	for i:=0;i<a.w;i++{
		resData[i] = make([]float64, a.h)
	}
	res := CreateMatrix(a.w,a.h, resData)
	for y:=0;y<a.h; y++{
		for x:=0;x<a.w;x++{
			res.data[x][y] = a.data[y][x]
		}
	}
	return res
}


func  MatrixLikeAnother(matr * Matrix) Matrix{
	data := make([][]float64, matr.h)
	for i:=0;i<matr.h;i++{
		data[i] = make([]float64, matr.w)
	}
	return CreateMatrix(matr.h, matr.w, data)

}
func(matr * Matrix) Determinant() float64{
	a := matr.data
	total := a[0][0] * (a[1][1]*a[2][2] - a[2][1]*a[1][2])+
		    -a[0][1] * (a[1][0]*a[2][2] - a[2][0]*a[1][2])+
		     a[0][2] * (a[1][0]*a[2][1] - a[2][0]*a[1][1])
	return total
}



func testTranspose(){
	testMatr2 :=CreateMatrix(3,2,[][]float64{
		{1,2},
		{1,4},
		{5,6},
	})
	fmt.Println(testMatr2.Transpose())

}

func testDeterminan(){
	testMatr := CreateMatrix(3,3,[][]float64{
		{1,2,3},
		{1,4,5},
		{5,6,7},
	})
	if testMatr.Determinant() != 8{
		log.Fatal("Неверно вычисляется определитель!")
	}
}

func testInverse(){

	testMatr := CreateMatrix(3,3,[][]float64{
		{1.0,2.0,3.0},
		{1.0,4.0,5.0},
		{5.0,6.0,7.0},
	})
	// Проверяем, что обратная ищется правильно, с помощью матричного умножения.
	invMatr := testMatr.Inverse()
	fmt.Println(DotMatrix(&testMatr, &invMatr))

}