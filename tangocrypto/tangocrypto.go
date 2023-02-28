package tangocrypto

import (
	"context"
	"encoding/hex"
	"fmt"
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
	butxos, err := t.client.AddressUTXOs(context.Background(), address.Bech32())
	if err != nil {
		if err, ok := err.(*tangocrypto_go.APIError); ok {
			if _, ok := err.Response.(tangocrypto_go.NotFound); ok {
				return []cardano.UTxO{}, nil
			}
		}
		return nil, err
	}

	utxos := make([]cardano.UTxO, len(butxos.Data))

	for i, butxo := range butxos.Data {
		txHash, err := cardano.NewHash32(butxo.Hash)
		if err != nil {
			return nil, err
		}

		amount := cardano.NewValue(0)

		if butxo.Assets != nil && len(butxo.Assets) > 0 {
			for _, asset := range butxo.Assets {
				h, err := hex.DecodeString(asset.PolicyID)
				if err != nil {
					return nil, err
				}
				policyID := cardano.NewPolicyIDFromHash(h[:28])
				currentAssets := amount.MultiAsset.Get(policyID)
				if currentAssets != nil {
					currentAssets.Set(
						cardano.NewAssetName(asset.AssetName),
						cardano.BigNum(asset.Quantity),
					)
				} else {
					amount.MultiAsset.Set(
						policyID,
						cardano.NewAssets().Set(
							cardano.NewAssetName(asset.AssetName),
							cardano.BigNum(asset.Quantity),
						),
					)
				}
			}
		} else {
			amount.Coin += cardano.Coin(butxo.Value)
		}

		utxos[i] = cardano.UTxO{
			Spender: address,
			TxHash:  txHash,
			Amount:  amount,
			Index:   uint64(butxo.Index),
		}
	}

	return utxos, nil
}

func (t TangoCryptoNode) Tip() (*cardano.NodeTip, error) {
	block, err := t.client.LatestBlock(context.Background())
	if err != nil {
		return nil, err
	}

	return &cardano.NodeTip{
		Block: uint64(block.BlockNo),
		Epoch: uint64(block.EpochNo),
		Slot:  uint64(block.SlotNo),
	}, nil
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
	b, err := t.client.LatestBlock(context.Background())
	if err != nil {
		return nil, err
	}

	eParams, err := t.client.ProtocolParameters(context.Background(), fmt.Sprintf("%d", b.EpochNo))
	if err != nil {
		return nil, err
	}

	params := cardano.ProtocolParams{
		MinFeeA:            cardano.Coin(eParams.MinFeeA),
		MinFeeB:            cardano.Coin(eParams.MinFeeB),
		MaxBlockBodySize:   uint(eParams.MaxBlockSize),
		MaxTxSize:          uint(eParams.MaxTxSize),
		MaxBlockHeaderSize: uint(eParams.MaxBlockHeaderSize),
		KeyDeposit:         cardano.Coin(eParams.KeyDeposit),
		PoolDeposit:        cardano.Coin(eParams.PoolDeposit),
		MaxEpoch:           uint(eParams.MaxEpoch),
		NOpt:               uint(eParams.OptimalPoolCount),
		CoinsPerUTXOWord:   cardano.Coin(eParams.MinUtxo),
	}

	return &params, nil
}

func (t TangoCryptoNode) Network() cardano.Network {
	return t.network
}
