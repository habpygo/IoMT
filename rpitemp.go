/*
Copyright 2016 yryz Author. All Rights Reserved.
DappDevelopment.com
*/

package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"iota/autocheck/mamutils"
	"strconv"
	"strings"
	"time"
	//"github.com/yryz/ds18b20"
)

func main() {
	sensors, err := Sensors()
	if err != nil {
		panic(err)
	}

	fmt.Printf("sensor IDs: %v\n", sensors)

	c := time.Tick(50 * time.Second)
	for _ = range c {
		for _, sensor := range sensors {
			t, err := Temperature(sensor) // t is float64
			if err == nil {
				fmt.Printf("sensor: %s temperature: %.2f°C\n", sensor, t)
				tStr := strconv.FormatInt(t, 10)
				mamutils.CreateConnectionAndSendMessage(0, tStr)
			}
		}
	}
}

var ErrReadSensor = errors.New("failed to read sensor temperature")

// Sensors get all connected sensor IDs as array
func Sensors() ([]string, error) {
	data, err := ioutil.ReadFile("/sys/bus/w1/devices/w1_bus_master1/w1_master_slaves")
	if err != nil {
		return nil, err
	}

	sensors := strings.Split(string(data), "\n")
	if len(sensors) > 0 {
		sensors = sensors[:len(sensors)-1]
	}

	return sensors, nil
}

// Temperature get the temperature of a given sensor
func Temperature(sensor string) (int64, error) {
	data, err := ioutil.ReadFile("/sys/bus/w1/devices/" + sensor + "/w1_slave")
	if err != nil {
		return 0.0, ErrReadSensor
	}

	raw := string(data) // data is []byte

	i := strings.LastIndex(raw, "t=")
	if i == -1 {
		return 0.0, ErrReadSensor
	}

	c, err := strconv.ParseInt(raw[i+2:len(raw)-1], 0, 64)
	if err != nil {
		return 0.0, ErrReadSensor
	}

	return c, nil
}
