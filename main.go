package main

import (
	"flag"
	"fmt"
	"time"

	"github.com/disintegration/imaging"
	"github.com/golang/glog"
)

const (
	TIME_BETWEEN_CHECKS = 3600 // seconds == 1 hour
	IMAGE_URL           = "https://i.groupme.com/373x506.png.3553b5442784435885844fc8ff026bba"
	IMAGE_FILENAME      = "./downloaded_image.png"
)

func configureAccount() {
	// 1. Check for existing account configuration

	// If non exists, go through the setup proc
	fmt.Println(generateCode())

}

func runService() {
	var signalOnline = make(chan bool)
	go waitForWifi(signalOnline)
	lastCheckTime := time.Time{}
	for {
		<-signalOnline

		if time.Now().Sub(lastCheckTime) > TIME_BETWEEN_CHECKS {
			glog.Info("Downloading Image")
			urlReader, err := getUrlReader(IMAGE_URL)
			if err != nil {
				glog.Info(err)
				glog.Info("Will try again next time.")
			} else {
				img, err := ioReaderToImage(urlReader)
				if err != nil {
					glog.Info("Error converting reader to image")
				}

				glog.Info("Saving Image")
				scaledImage := adjustImage(img)
				imaging.Save(scaledImage, IMAGE_FILENAME)

				// TODO: need to convert image to png
				// imageFile, err := os.Create(IMAGE_FILENAME)
				// if err != nil {
				// 	glog.Info(err)
				// } else {
				// 	_, err := io.Copy(imageFile, urlReader)
				// 	if err != nil {
				// 		glog.Info(err)
				// 	}
				// }

			}

		}
		break
	}

}

func main() {
	flag.Parse()
	glog.Info("Starting")
	configureAccount()
	runService()

}
