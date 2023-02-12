package cmd

import (
	"github.com/ripoff2/cardano-go"
	"github.com/ripoff2/cardano-go/blockfrost"
	"github.com/ripoff2/cardano-go/wallet"
	"github.com/spf13/cobra"
)

// deleteWalletCmd represents the deleteWallet command
var deleteWalletCmd = &cobra.Command{
	Use:     "delete-wallet [wallet]",
	Short:   "Delete a wallet",
	Aliases: []string{"delw"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		useTestnet, _ := cmd.Flags().GetBool("testnet")
		network := cardano.Mainnet
		if useTestnet {
			network = cardano.Testnet
		}

		node := blockfrost.NewNode(network, cfg.BlockfrostProjectID)
		opts := &wallet.Options{Node: node}
		client := wallet.NewClient(opts)
		id := args[0]
		return client.DeleteWallet(id)
	},
}

func init() {
	rootCmd.AddCommand(deleteWalletCmd)
	deleteWalletCmd.Flags().Bool("testnet", false, "Use testnet network")
}
