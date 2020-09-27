package connections

import (
	"iota/autocheck/metadata"

	"github.com/iotaledger/giota"
)

type Connection struct {
	api      *giota.API
	seed     string
	security int
	mwm      int64
}

//NewConnection establishes a connection with the given provider and the seed
func NewConnection(provider, seed string) (*Connection, error) {
	return &Connection{
		api:      giota.NewAPI(provider, nil),
		seed:     seed,
		security: metadata.SecurityLevel,
		mwm:      metadata.MWM,
	}, nil
}

func (c *Connection) SendToApi(trs []giota.Transfer) (giota.Bundle, error) {
	seed, err := giota.ToTrytes(c.seed)
	if err != nil {
		return nil, err
	}
	_, bestPow := giota.GetBestPoW()
	return giota.Send(c.api, seed, c.security, trs, c.mwm, bestPow)
}

/*
The way the library works: it makes use of FindTransactions to
figure out if there is any address in the Tangle that has already been used.
If FindTransactions returns associated transactions, the key index is
simply incremented, a new address generated and FindTransactions called
again until it returns null.
*/
func (c *Connection) FindTransactions(req giota.FindTransactionsRequest) ([]giota.Transaction, error) {
	found, err := c.api.FindTransactions(&req)
	if err != nil {
		return nil, err
	}
	return c.ReadTransactions(found.Hashes)
}

func (c *Connection) ReadTransactions(tIDs []giota.Trytes) ([]giota.Transaction, error) {
	found, err := c.api.GetTrytes(tIDs)
	return found.Trytes, err
}
