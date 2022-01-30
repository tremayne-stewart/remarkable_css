# Remarkable Changing Suspend Screen Service
reMarkable service to automatically download and replace your suspend screen. No cloud needed.

![](demo.jpg)


---

This was originallly a part of a larger project to control and update the suspend screen on the reMarkable 2 tablet from a web browser via a SAAS offering. 

The idea was as followed.
1. You downloaded the installer and plug in your reMarkable.
1. You would run the installer and create an account (via a webframe) which would link your device to that account.
1. The installer would install the necessary binaries (this repo) on your device.
1. You would use the web interface to upload images you wanted for you Suspend Screen, or select from the many premade ones.
1. Your device would periodically download one of your uploaded images and set it as the Suspend Screen.

I started working on the whole thing but decided it wasn't fully worth the effort since the general idea garnered [lukewarm support](https://www.reddit.com/r/RemarkableTablet/comments/sdk0b9/remarkable_suspend_screen_as_a_service/).

If you want any of the code for the other pieces, let me know and I'll dig the code out for you.

Special thanks to [remarkable_news](https://github.com/Evidlo/remarkable_news) for the springboard.


# Where does the service pull the images from?
Every time you connect to WiFi, it will try to grab the image specified by the `IMAGE_REFERENCE_URL`; which is currently pointed to a relay that grabs the latest front page of the New York Times.

Whatever address `IMAGE_REFERENCE_URL` points to needs to return a url string to an image. In the flow explained above, this endpoint was used to allow for multiple users to hit the same URL and have the backend process which user received which image url.




# Quickstart (Max/Linux)
Plug in you reMarkable and run the following.
_make sure to have your reMarkable's ssh password handy_
```
git clone git remote add origin https://github.com/tremayne-stewart/remarkable_css.git && cd remarkable_css
make install
```

This will install the remarkable_css service onto your reMarkable. 

# Quickstart (Windows)
I don't have a windows computer to test this on but I assume however you can run this with whatever setup you normally use to execute shell commands on your computer and on your reMarkable.

# Stopping
To stop service running on your reMarkable run
```
# Run in the remarkable_css directory
make stop
```


# Debugging
## Locally on the host machine
Run the following commmand which executes the program and saves the image would become the new suspend screen on the tablet to `downloaded_image.png`.
```
# Run in the remarkable_css directory
go run . -logtostderr=true -debug
```

## The Service running on the reMarkable
With the remarkable connected to the local network or via USB, execute the following command on the host machine (your computer). This will pull up a stream of the logs produced by the remarkable change suspend screen service.

```
# Run in the remarkable_css directory
make debug
```

Then disconnect and reconnect WiFi to trigger a download.


# Notes about the image loading
There was an update to the mainline that blocks hot swapping the suspended.png. To get around this you can either
1. Use `ddvk-hacks` which re-enables the hotswapping
1. Run `systemctl restart xochitl` when you want to pull in the new image
1. Restart the device.

Without `ddvk-hacks`, `xochitl` needs to be relaunched in order for the new suspended.png to be loaded. Some folks have sprung up ideas about somehow tieing the relaunch to the pressing the sleep/wake button but ... eh.. 

---

Dictated, not read.

