package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"ecg1/metadata"
	"ecg1/proto"

	//"golang.org/x/net/context"
	"context"

	"google.golang.org/grpc"
)

// MACID is the unique identifier of the device
const MACID = "00-14-22-01-23-45"

var ecg1 proto.BlockchainClient

const deviceFile = metadata.LOCALECG1OUT
const freq = 500
const offSet = -freq

var chunk = 0
var dataSlice []byte
var start = 0

// Counter is just for showing block number when transferring data blocks
var Counter int

func main() {

	//conn, err := grpc.Dial("localhost:8000", grpc.WithInsecure())
	fmt.Println("Dialing 83.80.200.16:50000")
	conn, err := grpc.Dial("83.80.200.16:50000", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot dial server, %v", err)
	}
	defer conn.Close()

	// create the client; the first ecg device, ecg1
	ecg1 = proto.NewBlockchainClient(conn)

	// create the file for this device for writing
	if err = addFile(deviceFile); err != nil {
		log.Fatalf("could not create file from device, %s: ", MACID)
	}

	if err := PrepareSimulation(); err != nil {
		log.Fatalf("could not prepare data file for simulation, %v", err)
	}
	// comment out for PEBL testing
	startSpinOff() // Start hashing data file and save it to the postgreSQL or FHIR DB
}

// addFile will request the gRPC server to safe a temp file where the data will be captured
func addFile(fileName string) error {

	fileResponse, err := ecg1.AddFile(context.Background(), &proto.AddFileRequest{
		FileName: fileName,
	})
	if err != nil {
		log.Fatalf("unable to add file to server: %v", err)
	}
	log.Printf("Created file for device: %t", fileResponse.Response)

	return nil
}

// PrepareSimulation generates a random ECG file to simulate an ECG device
func PrepareSimulation() error {
	ECGfile := metadata.LOCALECGIN0
	fmt.Println("ecgfile0.csv is:", ECGfile)

	// count the number of lines and use this int to loop through the slice
	TotNoOfLines, err := LineCounter(ECGfile)
	if err != nil {
		log.Fatalf("error during line counting of ECGFILE: %s", err)
	}
	fmt.Println("Total number of lines in file ECGFILE is: ", TotNoOfLines)

	allLinesInFile, err := scanFileToLines(ECGfile)
	if err != nil {
		log.Fatalf("Could not count lines, %s ", err)
	}

	c := time.Tick(1 * time.Second)
	for range c {
		if chunk >= TotNoOfLines {
			break
		}

		chunk = chunk + freq
		_ = processLines(allLinesInFile, TotNoOfLines, chunk, MACID)

	}

	return nil
}

// chunks of data are send to the server where the data will be put together
func processLines(lines []string, fileSize, chunk int, macid string) error {
	var blockRequests []*proto.AddBlockRequest

	data := lines[offSet+chunk : chunk] // 1st round -500 + 500 : 500 == 0:500

	dataSliceToAdd := proto.AddBlockRequest{
		Data:     data,
		Macid:    macid,
		Filename: deviceFile,
	}

	blockRequests = append(blockRequests, &dataSliceToAdd)

	/* --------- Here is where the gRPC and Blockchain magic happens --------- */
	stream, err := ecg1.AddBlock(context.Background())
	if err != nil {
		log.Fatalf("%v.AddBlock(_) = _, %v", ecg1, err)
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
	resp, err := ecg1.StartSpinOff(context.Background(), &proto.SpinOffRequest{
		Macid: MACID, // Every device has its own unique MACID
	})
	if err != nil {
	}
	fmt.Println("response from client is: ", resp)
}

func scanFileToLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string

	// create a scanner and read in all the lines in the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

//LineCounter returns the number of lines found in the file
func LineCounter(f string) (int, error) {
	ecg, err := os.Open(f)
	if err != nil {
		log.Fatalf("could not open file to count lines, %s", err)
	}
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	defer ecg.Close()

	for {
		c, err := ecg.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}
