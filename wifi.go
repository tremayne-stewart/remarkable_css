package main

// dbus filter for getting WiFi state change from wpa_supplicant
// see `dbus-monitor --system` and/or run that command on the device.
/* OUTPUT EXAMPLE:
	signal time=1642751912.229222 sender=:1.1 -> destination=(null destination) serial=984 path=/org/freedesktop/network1/link/_351; interface=org.freedesktop.DBus.Properties; member=PropertiesChanged
   string "org.freedesktop.network1.Link"
   array [
      dict entry(
         string "AddressState"
         variant             string "routable"
      )
      dict entry(
         string "OperationalState"
         variant             string "routable"
      )
   ]
   array [
   ]
*/

import (
	"os"
	"strings"

	"github.com/godbus/dbus"
	"github.com/golang/glog"
)

func waitForWifi(signalOnline chan bool) {
	glog.Info("Running waitForWifi")

	if IS_DEBUG {
		glog.Info("Bypassing DBUS Check for Wifi")
		signalOnline <- true
		return
	}

	// Use Dbus to wait for wpa_supplicant to emit message that we're connected
	// example: https://github.com/godbus/dbus/blob/c88335c0b1d28a30e7fc76d526a06154b85e5d97/_examples/monitor.go#L26
	busCon, err := dbus.SystemBus()
	if err != nil {
		glog.Fatalln(os.Stderr, "Failed to connect to system bus.", err)
	}
	defer busCon.Close()
	glog.Info("Connected to system bus.")

	var rules []string
	var flags uint = 0

	// First 2 arguments satisfy the library, second 2 satisfy the dbus command
	// More information (https://dbus.freedesktop.org/doc/dbus-specification.html)
	call := busCon.BusObject().Call("org.freedesktop.DBus.Monitoring.BecomeMonitor", 0, rules, flags)
	if call.Err != nil {
		glog.Fatalln(os.Stderr, "Failed to become monitor:", call.Err)
	}
	glog.Info("BecomeMonitor call success.")

	var wifiStep1Complete bool = false
	var messageChannel = make(chan *dbus.Message, 10)
	busCon.Eavesdrop(messageChannel)
	for message := range messageChannel {
		if !wifiStep1Complete && strings.Contains(message.String(), "routable") {
			glog.Info("Wifi connected.")
			wifiStep1Complete = true
		} else if wifiStep1Complete && strings.Contains(message.String(), "isSyncing") {
			glog.Info("Internet active")
			signalOnline <- true
		}
	}
	glog.Info("Done")
}
