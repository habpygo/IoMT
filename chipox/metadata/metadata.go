/*
Copyright DappDevelopment. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package metadata contains the paths used by the sdk
package metadata

// Project is the Go project name relative to the Go Path
var Project = ""

// LOCALECG1OUT is the path in the container; the path is set in the Dockerfile; use for containerized application
const LOCALECG1OUT = "/app/filevault/Chipox1out.csv"

// LOCALECGIN0 sends the data. The file is the location of the simulation files locally
// const LOCALECGIN0 = "/Users/harryboer/Developer/Dev_Projects/IoMT/ecg1/chipox1.csv"

// DATAFILE is the file where the data are saved from the Chipox
const DATAFILE = "\Users\hboer\Development\data\"