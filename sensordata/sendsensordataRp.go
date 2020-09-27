package main

import "iota/autocheck/mamutils"

func main() {
	mamutils.CheckForFileOnRp()
}

// dataplicity@raspberrypi:/home/pi/projects/RMAM$ ./sendsensordataRp
// panic: could not read file as it should, EOF

// goroutine 1 [running]:
// iota/autocheck/mamutils.ReadFile(0x25be92, 0x2c)
//         /Users/harryboer/Developer/DappDevelopment/Blockchain_Projects/Gowork/src/iota/autocheck/mamutils/sendsensordata.go:64 +0x2d4
// iota/autocheck/mamutils.CheckForFileOnRp()
//         /Users/harryboer/Developer/DappDevelopment/Blockchain_Projects/Gowork/src/iota/autocheck/mamutils/sendsensordata.go:41 +0x5c
// main.main()
//         /Users/harryboer/Developer/DappDevelopment/Blockchain_Projects/Gowork/src/iota/autocheck/sensordata/sendsensordataRp.go:6 +0x14
// dataplicity@raspberrypi:/home/pi/projects/RMAM$
