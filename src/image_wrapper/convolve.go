package image_wrapper


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



