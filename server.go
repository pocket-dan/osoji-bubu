package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/pocket-dan/osoji-bubu/camera"
	"github.com/pocket-dan/osoji-bubu/ml"
	"gocv.io/x/gocv"
	"gopkg.in/yaml.v2"
)

var countVR, countML int

var (
	VRCameraID      int
	VRCamera        *camera.Camera
	VRCamImageCache gocv.Mat
)

var (
	MLCameraID       int
	MLCamera         *camera.Camera
	MLCamResultCache MLCamResult
	MLModel          ml.Classifier
)

type MLCamResult struct {
	image    gocv.Mat
	score    int
	detected bool
}

type TrainingSettings struct {
	Train           bool   `yaml:"train"`
	NumberOfSamples int    `yaml:"num_samples"`
	SnapShotPath    string `yaml:"snapshot_path"`
}

type Config struct {
	VRCamera camera.CamSettings `yaml:"vr_camera"`
	MLCamera camera.CamSettings `yaml:"ml_camera"`
	Training TrainingSettings   `yaml:"training"`
}

func loadConfig(configFile string) (*Config, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	config := Config{}
	err = yaml.Unmarshal([]byte(data), &config)
	if err != nil {
		return nil, fmt.Errorf("unmarshal yaml: %w", err)
	}

	return &config, nil
}

var logger = log.New(os.Stderr, "", log.LstdFlags)

func initCameras(config *Config) error {
	var err error

	// camera for head mount display
	logger.Println("VR Camera: ", config.VRCamera)
	camConfig := config.VRCamera
	VRCamera, err = camera.NewCamera(camConfig)
	if err != nil {
		return fmt.Errorf("new camera: %w", err)
	}
	VRCamImageCache = gocv.NewMatWithSize(camConfig.Height, camConfig.Width, gocv.MatTypeCV8U)

	// camera for machine learning
	logger.Println("ML Camera: ", config.MLCamera)
	camConfig = config.MLCamera
	MLCamera, err = camera.NewCamera(config.MLCamera)
	if err != nil {
		return fmt.Errorf("new camera: %w", err)
	}
	MLCamResultCache = MLCamResult{
		image:    gocv.NewMatWithSize(camConfig.Height, camConfig.Width, gocv.MatTypeCV8U),
		score:    0,
		detected: false,
	}

	return nil
}

func captureVRCamLoop() {
	// VR Camera capturing loop
	var err error
	img := gocv.NewMatWithSize(VRCamera.Height(), VRCamera.Width(), gocv.MatTypeCV8U)
	for {
		err = VRCamera.GetImage(&img)
		if err != nil {
			logger.Printf("[WARN]: %s\n", fmt.Errorf("capture image: %w", err))
			continue
		}
		VRCamImageCache = img
		countVR++
	}
}

func captureMLCamLoop() {
	// ML Camera capturing and ML inference loop
	var score int
	var err error
	img := gocv.NewMatWithSize(MLCamera.Height(), MLCamera.Width(), gocv.MatTypeCV8U)
	for {
		err = MLCamera.GetImage(&img)
		if err != nil {
			logger.Printf("[WARN]: %s\n", fmt.Errorf("capture image: %w", err))
			continue
		}

		x, err := ml.ExtractFeature(&img)
		if err != nil {
			logger.Printf("[WARN]: %s\n", fmt.Errorf("extract feature from image: %w", err))
			continue
		}

		detected := MLModel.Predict(x) == -1
		if detected {
			score += 10
		}

		MLCamResultCache = MLCamResult{
			image:    img,
			detected: detected,
			score:    score,
		}
		countML++
	}
}

// api
type MLResponse struct {
	Image    string `json:"image"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Score    int    `json:"score"`
	Detected bool   `json:"detected"`
}

type VRResponse struct {
	Image    string `json:"image"`
	Width    int    `json:"width"`
	Height   int    `json:"height"`
	Score    int    `json:"score"`
	Detected bool   `json:"detected"`
}

func getVRImage(c echo.Context) error {
	imageJpeg, err := gocv.IMEncodeWithParams(gocv.JPEGFileExt, VRCamImageCache, []int{gocv.IMWriteJpegQuality, 85})
	if err != nil {
		return fmt.Errorf("encode image into jpeg: %w", err)
	}

	return c.JSON(http.StatusOK, &VRResponse{
		Image:    base64.StdEncoding.EncodeToString(imageJpeg),
		Width:    VRCamera.Width(),
		Height:   VRCamera.Height(),
		Score:    MLCamResultCache.score,
		Detected: MLCamResultCache.detected,
	})
}

func getMLImage(c echo.Context) error {
	imageJpeg, err := gocv.IMEncodeWithParams(gocv.JPEGFileExt, MLCamResultCache.image, []int{gocv.IMWriteJpegQuality, 85})
	if err != nil {
		return fmt.Errorf("encode image into jpeg: %w", err)
	}

	return c.JSON(http.StatusOK, &MLResponse{
		Image:  base64.StdEncoding.EncodeToString(imageJpeg),
		Width:  MLCamera.Width(),
		Height: MLCamera.Height(),
	})
}

func main() {
	runtime.GOMAXPROCS(2)

	var err error

	logger.Println("Load configurations...")
	var configFile string
	if len(os.Args) == 1 {
		configFile = "config/macbook-pro.yaml"
	} else {
		configFile = os.Args[1]
	}

	// load config
	config, err := loadConfig(configFile)
	if err != nil {
		logger.Println(err)
		os.Exit(1)
	}

	// initialize cameras
	logger.Println("Initialize cameras...")
	err = initCameras(config)
	if err != nil {
		logger.Println(err)
		os.Exit(1)
	}

	// new ml model
	MLModel = ml.NewGaussianClassifier()

	// train or load pretrained ML model
	if config.Training.Train {
		logger.Println("Start capturing to collect training samples...")
		XTrain := make([][]float32, config.Training.NumberOfSamples)
		YTrain := make([]int, config.Training.NumberOfSamples)
		img := gocv.NewMat()

		bar := pb.StartNew(config.Training.NumberOfSamples)
		for i := 0; i < config.Training.NumberOfSamples; i++ {
			err := MLCamera.GetImage(&img)
			if err != nil {
				logger.Println(err)
				os.Exit(1)
			}

			x, err := ml.ExtractFeature(&img)
			if err != nil {
				logger.Println(err)
				os.Exit(1)
			}
			XTrain[i] = x
			YTrain[i] = 1 // normal: 1, abnormal: -1
			bar.Increment()
		}
		bar.Finish()

		logger.Println("Start to train ML model...")
		err = MLModel.Train(XTrain, YTrain)
		if err != nil {
			logger.Println(err)
			os.Exit(1)
		}

		logger.Printf("Trained ML Model: %t\n", MLModel)
	} else {
		err = MLModel.Load(config.Training.SnapShotPath)
		if err != nil {
			logger.Println(err)
			os.Exit(1)
		}
	}

	// start capturing in background
	logger.Println("Start VR camera capture loop.")
	go captureVRCamLoop()

	logger.Println("Start ML camera capture loop.")
	go captureMLCamLoop()

	timer := time.NewTicker(time.Second * 5)
	go func() {
		for {
			<-timer.C
			logger.Printf("[VR] fps: %d\n", countVR/5)
			logger.Printf("[ML] fps: %d\n", countML/5)
			countVR = 0
			countML = 0
		}
	}()

	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet},
	}))

	// endpoints
	e.GET("/api/vr/capture", getVRImage) // vr camera image endpoint
	e.GET("/api/ml/capture", getMLImage) // ml camera image endpoint
	e.Static("/", "static")              // static files

	// start server
	e.Logger.Fatal(e.Start(":1323"))
}
