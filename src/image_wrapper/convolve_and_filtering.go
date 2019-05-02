package image_wrapper

import "sort"

func Convolve2d(a Matrix, kernel Matrix) Matrix{
	h, w := a.h, a.w
	kH, kW := kernel.h, kernel.w

	// Инициализируем результирующий массивчик.
	res := make([][]float64, h-kH+1)
	for i:=0;i<h-kH+1;i++{
		res[i] =  make([]float64, w-kW+1)
	}
	for y:=0; y<h-kH+1; y++{
		for x:=0; x<w-kW+1; x++{
			productRes := 0.0
			for i:=0; i < kH; i ++{
				for j:=0; j<kW; j++{
					productRes += a.data[y+i][x+j] * kernel.data[i][j]
				}
			}
			res[y][x] = productRes
		}
	}
	return CreateMatrix(h-kH+1, w-kW+1, res)
}

func RangeFilter(a Matrix, kernel Matrix, k int) Matrix{
	h, w := a.h, a.w
	kH, kW := kernel.h, kernel.w

	// Инициализируем результирующий массивчик.
	res := make([][]float64, h-kH+1)
	for i:=0;i<h-kH+1;i++{
		res[i] =  make([]float64, w-kW+1)
	}
	filter := make([]float64, int(kH * kW))
	for y:=0; y<h-kH+1; y++{
		for x:=0; x<w-kW+1; x++{
			idx := 0
			for i:=0; i < kH; i ++{
				for j:=0; j<kW; j++{
					if kernel.data[i][j] == 1{
						filter[idx] = a.data[y+i][x+j]
						idx += 1
					}
 				}
			}
			sort.Float64s(filter)
			res[y][x] = filter[k]
		}
	}
	return CreateMatrix(h-kH+1, w-kW+1, res)
}

