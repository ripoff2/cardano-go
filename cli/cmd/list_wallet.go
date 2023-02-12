package cmd

import (
	"fmt"

	"github.com/ripoff2/cardano-go"
	"github.com/ripoff2/cardano-go/blockfrost"
	"github.com/ripoff2/cardano-go/wallet"
	"github.com/spf13/cobra"
)

// listWalletCmd represents the listWallet command
var listWalletCmd = &cobra.Command{
	Use:     "list-wallet",
	Short:   "Print a list of known wallets",
	Aliases: []string{"lsw"},
	RunE: func(cmd *cobra.Command, args []string) error {
		useTestnet, _ := cmd.Flags().GetBool("testnet")
		network := cardano.Mainnet
		if useTestnet {
			network = cardano.Testnet
		}

		node := blockfrost.NewNode(network, cfg.BlockfrostProjectID)
		opts := &wallet.Options{Node: node}
		client := wallet.NewClient(opts)
		defer client.Close()
		wallets, err := client.Wallets()
		if err != nil {
			return err
		}
		fmt.Printf("%-18v %-9v %-9v\n", "ID", "NAME", "ADDRESS")
		for _, w := range wallets {
			addresses, err := w.Addresses()
			if err != nil {
				return err
			}
			fmt.Printf("%-18v %-9v %-9v\n", w.ID, w.Name, len(addresses))
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(listWalletCmd)
	listWalletCmd.Flags().Bool("testnet", false, "Use testnet network")
}
