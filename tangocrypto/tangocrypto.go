package tangocrypto

import (
	"context"
	"github.com/ripoff2/cardano-go"
	"github.com/ripoff2/tangocrypto-go"
)

type TangoCryptoNode struct {
	client  tangocrypto_go.APIClient
	appID   string
	apiKey  string
	network cardano.Network
}

func NewNode(network cardano.Network, appID string, apiKey string) cardano.Node {
	server := tangocrypto_go.CardanoMainNet

	if network == cardano.Testnet {
		server = tangocrypto_go.CardanoTestNet
	}

	return &TangoCryptoNode{
		network: network,
		appID:   appID,
		apiKey:  apiKey,
		client: tangocrypto_go.NewAPIClient(tangocrypto_go.APIClientOptions{
			AppID:  appID,
			ApiKey: apiKey,
			Server: server,
		}),
	}
}

func (t TangoCryptoNode) UTxOs(address cardano.Address) ([]cardano.UTxO, error) {
	//TODO implement me
	panic("implement me")
}

func (t TangoCryptoNode) Tip() (*cardano.NodeTip, error) {
	//TODO implement me
	panic("implement me")
}

func (t TangoCryptoNode) SubmitTx(tx *cardano.Tx) (*cardano.Hash32, error) {
	hash, err := t.client.TransactionSubmit(context.Background(), tx.Bytes())
	if err != nil {
		return nil, err
	}

	txHash, err := cardano.NewHash32(hash)
	if err != nil {
		return nil, err
	}

	return &txHash, nil
}

func (t TangoCryptoNode) ProtocolParams() (*cardano.ProtocolParams, error) {
	//TODO implement me
	panic("implement me")
}

func (t TangoCryptoNode) Network() cardano.Network {
	//TODO implement me
	panic("implement me")
}
