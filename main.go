package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/golang/glog"
)

const (
	TIME_BETWEEN_CHECKS = 3600 // seconds == 1 hour
	// IMAGE_URL            = "https://i.groupme.com/373x506.png.3553b5442784435885844fc8ff026bba"
	IMAGE_REFERENCE_URL  = "http://192.168.4.75:3000/get-image"
	DEBUG_IMAGE_FILENAME = "./downloaded_image.png"
	IMAGE_FILENAME       = "/usr/share/remarkable/suspended.png"
)

var IS_DEBUG bool = true

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
			urlReader, err := getUrlReader(IMAGE_REFERENCE_URL)
			if err != nil {
				glog.Info(err)
				glog.Info("Will try again next time.")
				break
			}

			referenceUrlBuffer := new(strings.Builder)
			_, err = io.Copy(referenceUrlBuffer, urlReader)
			if err != nil {
				glog.Info("Error getting image reference url from reader")
			}
			imageUrl := referenceUrlBuffer.String()
			urlReader, err = getUrlReader(imageUrl)
			if err != nil {
				glog.Info(err)
				glog.Info("Will try again next time.")
				break
			}

			img, err := ioReaderToImage(urlReader)
			if err != nil {
				glog.Info("Error converting reader to image")
			}

			glog.Info("Saving Image")
			scaledImage := adjustImage(img)

			fileSaveLocation := IMAGE_FILENAME
			if IS_DEBUG {
				fileSaveLocation = DEBUG_IMAGE_FILENAME
			}
			imaging.Save(scaledImage, fileSaveLocation)

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
		break
	}

}

func main() {
	flag.Parse()

	if env_debug := os.Getenv("DEBUG"); env_debug != "" {
		IS_DEBUG = "1" == os.Getenv("DEBUG")
		glog.Info(IS_DEBUG)
	}

	glog.Info("Starting")
	configureAccount()
	runService()

}
