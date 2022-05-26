package cmd

import (
	"fmt"
	"strconv"

	"github.com/echovl/cardano-go/node"
	"github.com/echovl/cardano-go/types"
	"github.com/echovl/cardano-go/wallet"
	"github.com/spf13/cobra"
)

// TODO: Ask for password if present
// Experimental feature, only for testnet
var transferCmd = &cobra.Command{
	Use:   "transfer [wallet-id] [amount] [receiver-address]",
	Short: "Transfer an amount of lovelace to the given address",
	Args:  cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		useTestnet, _ := cmd.Flags().GetBool("testnet")
		network := types.Mainnet
		if useTestnet {
			network = types.Testnet
		}
		opts := &wallet.Options{
			Node: node.NewCli(network),
		}
		client := wallet.NewClient(opts)
		defer client.Close()
		senderId := args[0]
		receiver := types.Address(args[2])
		amount, err := strconv.ParseUint(args[1], 10, 64)
		if err != nil {
			return err
		}
		w, err := client.Wallet(senderId)
		if err != nil {
			return err
		}
		txHash, err := w.Transfer(receiver, types.Coin(amount))
		fmt.Println(txHash)

		return err
	},
}

func init() {
	rootCmd.AddCommand(transferCmd)
	transferCmd.Flags().Bool("testnet", false, "Use testnet network")
}
