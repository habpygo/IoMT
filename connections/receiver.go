package connections

import (
	"fmt"
	"iota/merklemamgiota/mamutils"
	"sort"
	"time"

	"github.com/iotaledger/giota"
)

type Transaction struct {
	Message   string
	Value     int64
	Timestamp time.Time
	Recipient string
}

type ApiTransactionsFinder interface {
	FindTransactions(giota.FindTransactionsRequest) ([]giota.Transaction, error)
}

// ReadTransactions reads historic transaction data from a particular address
func ReadTransactions(address string, f ApiTransactionsFinder) ([]Transaction, error) {
	fmt.Println("address is: ", address)
	iotaAddress, err := giota.ToAddress(address)
	if err != nil {
		return nil, err
	}

	req := giota.FindTransactionsRequest{
		Addresses: []giota.Address{iotaAddress},
	}
	//fmt.Println("req.Addresses is: ", req.Addresses)

	foundTx, err := f.FindTransactions(req)
	if err != nil {
		return nil, err
	}

	sort.Slice(foundTx, func(i, j int) bool {
		return !(foundTx[i].Timestamp.Unix() > foundTx[j].Timestamp.Unix())
	})

	transactions := make([]Transaction, len(foundTx))
	fmt.Println("length of foundTx is: ", len(foundTx))
	for i, t := range foundTx {
		message, err := mamutils.FromMAMTrytes(t.SignatureMessageFragment)
		if err != nil {
			return nil, err
		}
		transactions[i] = Transaction{
			Message:   message,
			Value:     t.Value,
			Timestamp: t.Timestamp,
			Recipient: string(t.Address),
		}
		fmt.Println("Message is: ", transactions[i].Message)
	}
	return transactions, nil
}

// ReadSensorTransactions reads last x historic data and passes them to the plotter
func ReadSensorTransactions(address string, f ApiTransactionsFinder) ([]Transaction, error) {
	iotaAddress, err := giota.ToAddress(address)
	if err != nil {
		return nil, err
	}

	req := giota.FindTransactionsRequest{
		Addresses: []giota.Address{iotaAddress},
	}

	foundTx, err := f.FindTransactions(req)
	if err != nil {
		return nil, err
	}

	sort.Slice(foundTx, func(i, j int) bool {
		return !(foundTx[i].Timestamp.Unix() > foundTx[j].Timestamp.Unix())
	})

	transactions := make([]Transaction, len(foundTx))
	for i, t := range foundTx {
		message, err := mamutils.FromMAMTrytes(t.SignatureMessageFragment)
		if err != nil {
			return nil, err
		}
		transactions[i] = Transaction{
			Message:   message,
			Value:     t.Value,
			Timestamp: t.Timestamp,
			Recipient: string(t.Address),
		}
	}
	return transactions, nil
}

type ApiTransactionsReader interface {
	ReadTransactions([]giota.Trytes) ([]giota.Transaction, error)
}

// ReadTransaction reads a particular transaction
func ReadTransaction(transactionID string, r ApiTransactionsReader) (Transaction, error) {
	tID, err := giota.ToTrytes(transactionID)
	if err != nil {
		return Transaction{}, err
	}

	txs, err := r.ReadTransactions([]giota.Trytes{tID})
	if len(txs) != 1 {
		return Transaction{}, fmt.Errorf("Requested 1 Transaction but got %d", len(txs))
	}
	if err != nil {
		return Transaction{}, err
	}

	tx := txs[0]
	message, err := mamutils.FromMAMTrytes(tx.SignatureMessageFragment)
	if err != nil {
		return Transaction{}, err
	}
	transaction := Transaction{
		Message:   message,
		Value:     tx.Value,
		Timestamp: tx.Timestamp,
		Recipient: string(tx.Address),
	}

	return transaction, nil
}
