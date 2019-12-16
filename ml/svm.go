package ml

import (
	"fmt"

	"github.com/danieldk/gosvm"
)

type OneClassSVM struct {
	model  *gosvm.Model
	params *gosvm.Parameters
}

func (c *OneClassSVM) Train(x [][]float32, y []int) error {
	var s gosvm.TrainingInstance
	_x := make([]float64, len(x[0]))
	trainingSamples := gosvm.NewProblem()
	for i := 0; i < len(x); i++ {
		for j := 0; j < len(_x); j++ {
			_x[j] = float64(x[i][j])
		}
		s = gosvm.TrainingInstance{float64(y[i]), gosvm.FromDenseVector(_x)}
		trainingSamples.Add(s)
	}

	var err error
	c.model, err = gosvm.TrainModel(*c.params, trainingSamples)
	if err != nil {
		return fmt.Errorf("model training: %w", err)
	}

	return nil
}

func (c *OneClassSVM) Predict(x []float32) int {
	_x := make([]float64, len(x))
	for i := 0; i < len(x); i++ {
		_x[i] = float64(x[i])
	}
	pred := c.model.Predict(gosvm.FromDenseVector(_x))

	return int(pred)
}

func (c *OneClassSVM) Save(path string) error {
	return c.model.Save(path)
}

func (c *OneClassSVM) Load(path string) error {
	model, err := gosvm.LoadModel(path)
	if err != nil {
		return fmt.Errorf("load model: %w", err)
	}

	c.model = model
	return nil
}

func NewClassifier() *OneClassSVM {
	return &OneClassSVM{
		&gosvm.Model{},
		&gosvm.Parameters{
			SVMType: gosvm.NewOneClass(0.5),
			// SVMType: gosvm.NewCSVC(1.0),
			// Kernel: gosvm.NewLinearKernel(),
			Kernel: gosvm.NewRBFKernel(1.0 / 256.0),
			// Kernel: gosvm.NewPolynomialKernel(1.0/256.0, 0.0, 4),
			// Kernel: gosvm.NewSigmoidKernel(1.0/256.0, 0.0),
			CacheSize:   200,
			Epsilon:     0.001,
			Shrinking:   false,
			Probability: false,
		},
	}
}
