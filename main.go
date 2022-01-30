package main

import (
	"flag"
	"io"
	"strings"
	"time"

	"github.com/disintegration/imaging"
	"github.com/golang/glog"
)

const (
	TIME_BETWEEN_CHECKS = 3600 // seconds == 1 hour
	// This URL needs to return a the URL of the image to use as the suspend screen (plain text)
	IMAGE_REFERENCE_URL  = "https://uriel-rmcss.herokuapp.com/latest-nyt"
	DEBUG_IMAGE_FILENAME = "./downloaded_image.png"
	IMAGE_FILENAME       = "/usr/share/remarkable/suspended.png"
)

var IS_DEBUG bool = false

func runService() {
	var signalOnline = make(chan bool)
	go waitForWifi(signalOnline)
	lastCheckTime := time.Time{}
	for {
		<-signalOnline

		if time.Now().Sub(lastCheckTime).Seconds() > TIME_BETWEEN_CHECKS {
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
			glog.Info("Using Image: ", imageUrl)
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

			scaledImage := adjustImage(img)

			fileSaveLocation := IMAGE_FILENAME
			if IS_DEBUG {
				fileSaveLocation = DEBUG_IMAGE_FILENAME
			}
			glog.Info("Saving Image to ", fileSaveLocation)
			imaging.Save(scaledImage, fileSaveLocation)

			glog.Info("Done")
		}
	}

}

func main() {
	debugFlag := flag.Bool("debug", false, "Saves the image to "+DEBUG_IMAGE_FILENAME)
	flag.Parse()

	glog.Info("Starting")
	if *debugFlag {
		IS_DEBUG = true
		glog.Info("Debug")
	}

	runService()

}
