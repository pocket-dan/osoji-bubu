local:
	env GO111MODULE=on go run server.go

build:
	# env GOOS=linux GOARCH=arm GOARM=7 go build
	env GO111MODULE=on go build

run:
	env GO111MODULE=on ./osoji-bubu config/raspberry-pi.yaml

list-cameras:
	v4l2-ctl --list-devices
