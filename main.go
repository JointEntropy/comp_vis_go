package main

import (
	"./src/image_wrapper"
	"fmt"
)

func main(){
	img1 := image_wrapper.LoadMImageWrapperFromString("images/Arch2.JPG")
	analyzer := image_wrapper.NewAnalyzer(img1)
	fmt.Println("Mean", image_wrapper.CalcMean(img1.GetData(),analyzer.ProbsData))
	fmt.Println("Var", image_wrapper.CalcVar(img1.GetData(),analyzer.ProbsData))
	fmt.Println("Energy", image_wrapper.CalcEnergy(img1.GetData(),analyzer.ProbsData))
	fmt.Println("Entropy", image_wrapper.CalcEntropy(img1.GetData(),analyzer.ProbsData))
	fmt.Println("Kurtosis", image_wrapper.CalcKurtosis(img1.GetData(),analyzer.ProbsData))
	fmt.Println("Skewness", image_wrapper.CalcSkewness(img1.GetData(),analyzer.ProbsData))
	res := image_wrapper.CoocurenceMatrix(img1.GetData(), 100, 100)
	img1.UpdateFromStatic(res)
	img1.SaveImage("output", "png")
}

