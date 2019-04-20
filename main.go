package main

import (
	"./src/image_wrapper"
	"fmt"
)

func main(){
	img1 := image_wrapper.LoadMImageWrapperFromString("images/grey_sanity_check.jpg")
	res := image_wrapper.CoocurenceMatrix(img1.GetData(), 100, 100)
	img1.UpdateFromStatic(res)
	analyzer := image_wrapper.NewAnalyzer(img1)

	fmt.Println("HistData", analyzer.HistData)
	fmt.Println("ProbsData", analyzer.ProbsData[120])
	fmt.Println("Mean", image_wrapper.CalcMean(img1.GetData(),analyzer.ProbsData))
	fmt.Println("Var", image_wrapper.CalcVar(img1.GetData(),analyzer.ProbsData))
	fmt.Println("Energy", image_wrapper.CalcEnergy(img1.GetData(),analyzer.ProbsData))
	fmt.Println("Entropy", image_wrapper.CalcEntropy(img1.GetData(),analyzer.ProbsData))
	fmt.Println("Kurtosis", image_wrapper.CalcKurtosis(img1.GetData(),analyzer.ProbsData))
	fmt.Println("Skewness", image_wrapper.CalcSkewness(img1.GetData(),analyzer.ProbsData))
	img1.SaveImage("output", "png")
}

