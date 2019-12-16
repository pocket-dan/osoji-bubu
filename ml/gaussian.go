package ml

import (
	"math"
)

type GaussianParameters struct {
	mean float32
	std  float32
}

type GaussianModel struct {
	params *GaussianParameters
}

func sum(d []float32) float32 {
	var s float32
	for i := 0; i < len(d); i++ {
		s += d[i]
	}
	return s
}

func mean(d []float32) float32 {
	return sum(d) / float32(len(d))
}

func std(d []float32) float32 {
	mu := mean(d)

	var s float32
	for i := 0; i < len(d); i++ {
		s += float32(math.Pow(float64(d[i]-mu), 2))
	}

	return float32(math.Sqrt(float64(s / float32(len(d)))))
}

func (g *GaussianModel) Train(x [][]float32, y []int) error {
	_x := make([]float32, len(x))
	for i := 0; i < len(x); i++ {
		_x[i] = sum(x[i])
	}

	// use mean and std
	g.params = &GaussianParameters{
		mean: mean(_x),
		std:  std(_x),
	}

	return nil
}

func (g *GaussianModel) Predict(x []float32) int {
	d := sum(x) - g.params.mean
	if d < 0 {
		d = -d
	}

	ret := 1
	if d > 3*g.params.std {
		ret = -1
	}

	return ret
}

func (g *GaussianModel) Save(path string) error {
	return nil
}

func (g *GaussianModel) Load(path string) error {
	return nil
}

func NewGaussianClassifier() *GaussianModel {
	return &GaussianModel{}
}
