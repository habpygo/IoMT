package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"time"

	"chipox/metadata"
	"chipox/proto"

	//"golang.org/x/net/context"
	"context"

	"google.golang.org/grpc"
)

// MACID is the unique identifier of the device
const MACID = "00-14-22-01-23-45"

var chipox1 proto.BlockchainClient

// DeviceDataOutputFile contains the data from the Chipox device
const DeviceDataOutputFile = metadata.LOCALECG1OUT

var BeginValueOfSlice = 0

//const freq = 500
//const offSet = -freq

//var chunk = 0
var dataSlice []byte

//var start = 0

// Counter is just for showing block number when transferring data blocks
var Counter int

// ErrReadSensor is the error message when reading data fails
var ErrReadSensor = errors.New("failed to read sensor temperature")
var O2File string

func main() {

	//conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	fmt.Println("Dialing 83.80.200.16:50000")
	conn, err := grpc.Dial("83.80.200.16:50000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot dial server, %v", err)
	}
	defer conn.Close()

	// create the client; the first chipox device, chipox1
	chipox1 = proto.NewBlockchainClient(conn)

	// read the file given by Chipox; NOTE: clean this file out after every run
	for {
		time.Sleep(10 * time.Second)
		files, err := ioutil.ReadDir(metadata.DATADIR)
		if err != nil {
			log.Fatalf("Error while reading files in data directory, %s", err)
		}
		if len(files) == 0 {
			fmt.Println("File not yet opened by Chipox")
			continue
		} else {
			O2File := files[0].Name()
			fmt.Println("Chipox1out.csv is:", O2File)
			// create the file for this device for writing
			if err = addFile(DeviceDataOutputFile); err != nil {
				log.Fatalf("could not create file from device, %s: ", err)
			}
			if err := ReadDataFromFile(O2File); err != nil {
				log.Fatalf("could not read data file from Chipox, %v", err)
			}
		}
	}

	// comment out for PEBL testing
	// startSpinOff() // Start hashing data file and save it to the postgreSQL or FHIR DB
}

// addFile will request the gRPC server to safe a temp file where the data will be captured
// fileName is arbitrarily set in the metadata file
func addFile(fileName string) error {

	fileResponse, err := chipox1.AddFile(context.Background(), &proto.AddFileRequest{
		FileName: fileName,
	})
	if err != nil {
		log.Fatalf("unable to add file to server: %v", err)
	}
	log.Printf("Created file for device: %t", fileResponse.Response)

	return nil
}

// ReadDataFromFile generates a random ECG file to simulate an ECG device
// chunks of data are send to the server where the data will be put together
func ReadDataFromFile(filename string) error {
	fmt.Println("Entering the ReadDataFromFile function")
	path := metadata.DATADIR + filename
	//data, err := ioutil.ReadFile(O2File)
	//if err != nil {
	//log.Fatalf("error is: %s", err)
	//}
	//input := string(data)
	//strings.Replace(input, "\r?\n?", "@", -1)
	// data := lines[offSet+chunk : chunk] // 1st round (-500 + 500) : 500 == 0:500

	c := time.Tick(1 * time.Second)
	for range c {
		TotNoOfLines := LineCounter(path)
		chunkToBeProcessed := TotNoOfLines - BeginValueOfSlice
		fmt.Println("Total number of lines in file ECGFILE is: ", TotNoOfLines)
		allLinesInFile, err := scanFileToLines(path)
		if err != nil {
			log.Fatalf("Could not count lines, %s ", err)
		}

		_ = ProcessLines(allLinesInFile, TotNoOfLines, BeginValueOfSlice, chunkToBeProcessed, MACID)
		BeginValueOfSlice = BeginValueOfSlice + chunkToBeProcessed

	}

	return nil
}
func scanFileToLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string

	// create a scanner and read in all the lines in of file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

//LineCounter returns the number of lines found in the file
func LineCounter(f string) int {
	chipox, err := os.Open(f)
	if err != nil {
		log.Fatalf("could not open file to count lines, %s", err)
	}
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	defer chipox.Close()

	for {
		c, err := chipox.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count

		case err != nil:
			return count
		}
	}
}

// ProcessLines function is used to read a ready file and is not used for life updates
func ProcessLines(lines []string, fileSize, beginValueOfSlice, chunkToBeProcessed int, macid string) error {
	var blockRequests []*proto.AddBlockRequest

	data := lines[beginValueOfSlice:chunkToBeProcessed] // 1st round (-500 + 500) : 500 == 0:500

	dataSliceToAdd := proto.AddBlockRequest{
		Data:     data,
		Macid:    macid,
		Filename: DeviceDataOutputFile,
	}

	blockRequests = append(blockRequests, &dataSliceToAdd)

	/* --------- Here is where the gRPC and Blockchain magic happens --------- */
	stream, err := chipox1.AddBlock(context.Background())
	if err != nil {
		log.Fatalf("%v.AddBlock(_) = _, %v", chipox1, err)
	}

	for _, block := range blockRequests {
		if err := stream.Send(block); err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("%v.Send(%v) = %v", stream, block, err)
		}
	}
	return nil
}

func startSpinOff() {
	resp, err := chipox1.StartSpinOff(context.Background(), &proto.SpinOffRequest{
		Macid: MACID, // Every device has its own unique MACID
	})
	if err != nil {
	}
	fmt.Println("response from client is: ", resp)
}

// func PrepareSimulation() error {
// 	O2File := metadata.LOCALO2IN
// 	fmt.Println("ecgfile0.csv is:", O2File)

// 	// count the number of lines and use this int to loop through the slice
// 	TotNoOfLines, err := LineCounter(O2File)
// 	if err != nil {
// 		log.Fatalf("error during line counting of ECGFILE: %s", err)
// 	}
// 	fmt.Println("Total number of lines in file ECGFILE is: ", TotNoOfLines)

// 	allLinesInFile, err := scanFileToLines(O2File)
// 	if err != nil {
// 		log.Fatalf("Could not count lines, %s ", err)
// 	}

// 	c := time.Tick(1 * time.Second)
// 	for range c {
// 		if chunk >= TotNoOfLines {
// 			break
// 		}

// 		chunk = chunk + freq
// 		_ = processLines(allLinesInFile, TotNoOfLines, chunk, MACID)

// 	}

// 	return nil
// }
