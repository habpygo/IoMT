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
	"path/filepath"
	"time"

	"chipox/metadata"
	"chipox/proto"

	//"golang.org/x/net/context"
	"context"

	"google.golang.org/grpc"
)

// MACID is the unique identifier of the device
const MACID = "00-1B-63-84-32-F8"

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
var ErrReadSensor = errors.New("failed to read sensor output")

// O2File is the file created by Chipox
var O2File string

func main() {

	fmt.Println("Dialing ")
	conn, err := grpc.Dial("", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("cannot dial server, %v", err)
	}
	defer conn.Close()

	// create the client; the first chipox device, chipox1
	chipox1 = proto.NewBlockchainClient(conn)

	for {
		time.Sleep(10 * time.Second)
		fmt.Println("file path is: ", filepath.Dir(metadata.DATADIR))
		files, err := ioutil.ReadDir(filepath.Dir(metadata.DATADIR))
		if err != nil {
			log.Fatalf("Error while reading files in data directory, %s", err)
		}
		if len(files) == 0 {
			fmt.Println("File not yet opened by Chipox")
			continue
		} else {
			O2File = files[0].Name()
			fmt.Println("File found: ", O2File)
			break
		}
	}
	fmt.Println("Chipox1out.csv is:", O2File)
	// create the file for this device for writing
	if err = addFile(DeviceDataOutputFile); err != nil {
		log.Fatalf("could not create file from device, %s: ", err)
	}
	fmt.Println("File added: ", DeviceDataOutputFile)
	if err := ReadDataFromFile(O2File); err != nil {
		log.Fatalf("could not read data file from Chipox, %v", err)
	}
	startSpinOff()
	// Clean out directory
	err = RemoveContents(filepath.Dir(metadata.DATADIR))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
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
	prevTotLines := 0

	c := time.Tick(1 * time.Second)
	for range c {
		TotNoOfLines := LineCounter(path)
		if prevTotLines == TotNoOfLines {
			break
		}
		prevTotLines = TotNoOfLines
		chunkToBeProcessed := TotNoOfLines - BeginValueOfSlice
		fmt.Println("Total number of lines in file PO2 is: ", TotNoOfLines)

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

	data := lines[beginValueOfSlice : beginValueOfSlice+chunkToBeProcessed] // 1st round (-500 + 500) : 500 == 0:500

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
	fmt.Println("Starting SpinOff")
	resp, err := chipox1.StartSpinOff(context.Background(), &proto.SpinOffRequest{
		Macid: MACID, // Every device has its own unique MACID
	})
	if err != nil {
		fmt.Println("Error in startSpinOff: ", err)
	}
	fmt.Println("response from client is: ", resp)
}

// RemoveContent cleans out the directory of chipox once the data have been send
func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}

// read the file given by Chipox; NOTE: clean this file out after every run
// for {
// 	time.Sleep(10 * time.Second)
// 	fmt.Println("file path is: ", filepath.Dir(metadata.DATADIR))
// 	files, err := ioutil.ReadDir(filepath.Dir(metadata.DATADIR))
// 	if err != nil {
// 		log.Fatalf("Error while reading files in data directory, %s", err)
// 	}
// 	if len(files) == 0 {
// 		fmt.Println("File not yet opened by Chipox")
// 		continue
// 	} else {
// 		O2File := files[0].Name()
// 		fmt.Println("Chipox1out.csv is:", O2File)
// 		// create the file for this device for writing
// 		if err = addFile(DeviceDataOutputFile); err != nil {
// 			log.Fatalf("could not create file from device, %s: ", err)
// 		}
// 		if err := ReadDataFromFile(O2File); err != nil {
// 			log.Fatalf("could not read data file from Chipox, %v", err)
// 		}
// 		if err == io.EOF {
// 			startSpinOff()
// 			// Clean out directory
// 			err = RemoveContents(filepath.Dir(metadata.DATADIR))
// 			if err != nil {
// 				fmt.Println(err)
// 				os.Exit(1)
// 			}
// 		}

// 	}
// }
