package connections

import (
	"fmt"
	"iota/merklemamgiota/mamutils"

	"github.com/iotaledger/giota"
)

type ApiSender interface {
	SendToApi([]giota.Transfer) (giota.Bundle, error)
}

func Send(recipient string, value int64, message string, sender ApiSender) (string, error) {
	address, err := giota.ToAddress(recipient)
	if err != nil {
		return "", err
	}

	// Mask messasge
	encodedMessage, err := mamutils.ToMAMTrytes(message)
	if err != nil {
		panic(fmt.Errorf("could not encode message to trytes, %v ", err))
	}

	trs := []giota.Transfer{
		giota.Transfer{
			Address: address,
			Value:   value,
			Message: encodedMessage,
			Tag:     "",
		},
	}

	mamBundle, sendErr := sender.SendToApi(trs)
	if sendErr != nil {
		return "", sendErr
	}

	return string(mamBundle[0].Hash()), nil
}
