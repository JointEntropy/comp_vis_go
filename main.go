package main

import (
	"./src/image_wrapper"
	"fmt"
	"log"
	"math"
)

func main(){
	img1 := image_wrapper.LoadMImageWrapperFromString("images/cat1.jpg")
	img1.SaveImage("output", "png")
	imgMatr := img1.ToMatrix()

	binary := image_wrapper.Binarize(imgMatr, 150)
	//resMatr.DumpToFile("output/binarize.txt")

	//проверка валидности вычислений моментов
	eps := 1e-3
	m00 := image_wrapper.CentralTwoDimensionalMoment(&binary, 0,0)
	cm00 := image_wrapper.TwoDimensionalMoment(&binary, 0,0)
	if  math.Abs(m00 - cm00)>eps{
		log.Fatal(fmt.Sprintf("Не выполняется равенство нулевого и центрального моментов %f!=%f",
			m00, cm00))
	}
	cm10 := image_wrapper.CentralTwoDimensionalMoment(&binary, 1,0)
	cm01 := image_wrapper.CentralTwoDimensionalMoment(&binary, 0,1)
	if  (math.Abs(cm10 - 0)>eps)  || (math.Abs(cm01 - 0)>eps){
		log.Fatal(fmt.Sprintf("Не выполняется равенство центральных моментов первого порядка %f, %f: !=0",
			cm10, cm01))
	}

	// Вычисление моментов.
	fmt.Printf("Геометрический момент нулевого порядка %f\n",
				image_wrapper.TwoDimensionalMoment(&binary, 0,0))
	fmt.Printf("Геометрический момент первого порядка p=1, q=0: %f\n",
				image_wrapper.TwoDimensionalMoment(&binary, 1,0))
	fmt.Printf("Геометрический момент первого порядка p=0, q=1: %f\n",
		image_wrapper.TwoDimensionalMoment(&binary, 0,1))
	fmt.Printf("Центральный геометрический момент второго порядка %f\n",
				image_wrapper.CentralTwoDimensionalMoment(&binary, 1,1))
	fmt.Printf("Нормализованный центральный геометрический момент второго порядка %f\n",
		image_wrapper.NormalizedCentralTwoDimensionalMoment(&binary, 1,1))
	image_wrapper.FromMatrix(binary, "binarize", "jpg").SaveImage("output", "jpg")

	img1 = image_wrapper.LoadMImageWrapperFromString("images/synt.png")
	labeledMask := image_wrapper.LabelSectors4d(img1.ToMatrix())
	image_wrapper.FromMatrix(labeledMask, "sectors4d", "png").SaveImage("output", "png")

}

