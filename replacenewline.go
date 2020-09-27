package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"regexp"
)

//const COR12FILE = "COR12/static_datafiles/500Hz_5min_NSR90"
const COR12FILE = "testtext.rtf"

func main() {
	data, err := ioutil.ReadFile(COR12FILE)
	if err != nil {
		log.Fatalf("error is: ", err)
	}
	input := string(data)
	//fmt.Println(input)
	// strings.Replace(input, "\r?\n?", "@", -1)
	//fmt.Println(input)
	re := regexp.MustCompile(`\n`)
	input = re.ReplaceAllString(input, "xxx")
	fmt.Println(input)

}
