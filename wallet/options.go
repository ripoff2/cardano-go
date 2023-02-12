package wallet

import (
	"github.com/ripoff2/cardano-go"
	cardanocli "github.com/ripoff2/cardano-go/cardano-cli"
)

type Options struct {
	Node cardano.Node
	DB   DB
}

func (o *Options) init() {
	if o.Node == nil {
		o.Node = cardanocli.NewNode(cardano.Testnet)
	}
	if o.DB == nil {
		o.DB = newMemoryDB()
	}
}
