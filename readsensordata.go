package main

//Reading files requires checking most calls for errors. This helper will streamline our error checks below.

import (
	"iota/autocheck/mamutils"
	"time"
)

// CheckForFile reads the file written by sensor on SBC and sends it to the Tangle
func CheckForData() {
	c := time.Tick(5 * time.Second)
	for _ = range c {
		mamutils.CreateConnectionAndReceiveSensorMessages()
	}
}

func main() {
	CheckForData()
}
