package main

import (
	"bufio"
	"fmt"
	"iota/autocheck/mamutils"
	"os"
)

func main() {
	fmt.Println("BEFORE YOU CONTINUE, CHECK AVAILABILITY OF END-POINT FIRST, BY SENDING AND READING MESSAGES.")
	fmt.Println("1. Send a message")
	fmt.Println("2. Read all messages from the IOTA tangle")
	fmt.Println("3. Automatically check for sensor messages every 5 seconds and send it to the IOTA Tangle")
	fmt.Println("----------------------------------------------------")
	fmt.Println("Test the messaging on the web: https://devnet.thetangle.org/ and paste address below.")
	fmt.Println("HZIQ9FHDEQEQLMKFVXLMLEWXBUZQORXGWIWFMSOCWA9KMIXRJ9HTTUPTQMUDFWBCUVUZSACHECQQGPHNC")
	fmt.Println("----------------------------------------------------")
	fmt.Println("Your choice: Press 1, 2 or 3")
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		fmt.Println(err)
	}

	switch char {
	case '1':
		message := mamutils.CaptureMessage()
		mamutils.CreateConnectionAndSendMessage(message)
	case '2':
		mamutils.CreateConnectionAndReceiveMessages()
		break
	case '3':
		fmt.Println("Check for messages...")
		mamutils.CheckForFileOnRp()
	default:
		fmt.Println("Wrong input. Only option 1, 2 or 3 is allowed. Try again!")
		return
	}

	fmt.Println("Successful connection and data transfer")
}
