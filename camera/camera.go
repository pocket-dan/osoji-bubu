package camera

import (
	"errors"
	"fmt"

	"gocv.io/x/gocv"
)

type CamSettings struct {
	ID     int `yaml:"device"`
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
	FPS    int `yaml:"fps"`
}

type Camera struct {
	webcam   *gocv.VideoCapture
	settings *CamSettings
}

func NewCamera(config CamSettings) (*Camera, error) {
	webcam, err := gocv.OpenVideoCapture(config.ID)
	if err != nil {
		return nil, err
	}
	webcam.Set(gocv.VideoCaptureFrameWidth, float64(config.Width))
	webcam.Set(gocv.VideoCaptureFrameHeight, float64(config.Height))
	webcam.Set(gocv.VideoCaptureFPS, float64(config.FPS))

	return &Camera{webcam, &config}, nil
}

func (c *Camera) GetImage(img *gocv.Mat) error {
	if ok := c.webcam.Read(img); !ok {
		msg := fmt.Sprintf("cannot read device %v\n", c.settings.ID)
		return errors.New(msg)
	}

	if img.Empty() {
		return errors.New("image is empty")
	}

	return nil
}

func (c *Camera) DeviceID() int {
	return c.settings.ID
}

func (c *Camera) Width() int {
	return c.settings.Width
}
func (c *Camera) Height() int {
	return c.settings.Height
}
