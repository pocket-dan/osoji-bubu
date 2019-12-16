package ml

import (
	"fmt"

	"gocv.io/x/gocv"
)

type Classifier interface {
	Train([][]float32, []int) error
	Predict([]float32) int
	Save(string) error
	Load(string) error
}

func ExtractFeature(img *gocv.Mat) ([]float32, error) {
	// convert into mono
	if img.Channels() == 3 {
		_img := gocv.NewMatWithSize(img.Rows(), img.Cols(), img.Type())
		gocv.CvtColor(*img, &_img, gocv.ColorBGRToGray)
		img = &_img
	}

	// calculate histgram
	hist, mask := gocv.NewMat(), gocv.NewMat()
	gocv.CalcHist([]gocv.Mat{*img}, []int{0}, mask, &hist, []int{255}, []float64{0, 255}, false)

	// convert into []float32
	histValues, err := hist.DataPtrFloat32()
	if err != nil {
		return []float32{}, fmt.Errorf("retrieve data array from gocv.Mat: %w", err)
	}

	// use only white part
	t := 80
	histValues = histValues[t:]

	// scale to [0. 1]
	size := float32(img.Cols() * img.Rows())
	for i := 0; i < len(histValues); i++ {
		histValues[i] = histValues[i] / size
	}

	return histValues, nil
}
