/*
Copyright DappDevelopment. All Rights Reserved.
SPDX-License-Identifier: Apache-2.0
*/

// Package metadata contains the paths used by the sdk
package metadata

// Project is the Go project name relative to the Go Path
var Project = ""

// DATADIR is the directory where the data file is saved from the Chipox device
const DATADIR = "/c/Users/hboer/Development/IoMT/chipox/data"

const LOCALECG1IN = "/Users/harryboer/Developer/Dev_Projects/IoMT/ecg/filegenerator/ecgfile1.csv"
const LOCALECG2IN = "/Users/harryboer/Developer/Dev_Projects/IoMT/ecg/filegenerator/ecgfile2.csv"
const LOCALECG3IN = "/Users/harryboer/Developer/Dev_Projects/IoMT/ecg/filegenerator/ecgfile3.csv"
const LOCALBLOODPIN = "/Users/harryboer/Developer/Dev_Projects/IoMT/ecg/filegenerator/bloodp.csv"
const LOCALO2IN = "/Users/harryboer/Developer/Dev_Projects/IoMT/ecg/filegenerator/Chipox1.csv"

// LOCALECG1OUT is the path in the container; the path is set in the Dockerfile; use for containerized application
const LOCALECG1OUT = "/app/filevault/Chipox1out.csv"

// ------------------------------------------------------

const LOCALECG2OUT = "/app/filevault/ECG2out.csv"
const LOCALECG3OUT = "/app/filevault/ECG3out.csv"
const LOCALBLOODPOUT = "/app/filevault/BLOODPressure.csv"
const LOCALO2OUT = "/app/filevault/Chipox.csv"

const MAC1 = "00-14-22-01-23-45"
const MAC2 = "00-04-DC-20-09-11"
const MAC3 = "00-3C-77-14-88-D3"
const MAC4 = "00-55-96-21-43-E6"
const MAC5 = "00-1B-63-84-32-F8"
