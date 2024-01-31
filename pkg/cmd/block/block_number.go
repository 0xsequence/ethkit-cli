package block

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/spf13/cobra"

	"github.com/0xsequence/ethkit/ethrpc"
)

const (
	flagBlockNumberRpcUrl = "rpc-url"
)

func NewBlockNumberCmd() *cobra.Command {
	c := &blockNumber{}
	cmd := &cobra.Command{
		Use:     "block-number",
		Short:   "Get the latest block number for a given blockchain network",
		Aliases: []string{"bn"},
		Args:    cobra.NoArgs,
		RunE:    c.Run,
	}

	cmd.Flags().StringP(flagBlockNumberRpcUrl, "r", "", "The RPC endpoint to the blockchain node to interact with")

	return cmd
}

type blockNumber struct {
}

func (c *blockNumber) Run(cmd *cobra.Command, args []string) error {
	fRpc, err := cmd.Flags().GetString(flagBlockNumberRpcUrl)
	if err != nil {
		return err
	}

	if _, err = url.ParseRequestURI(fRpc); err != nil {
		return errors.New("error: please provide a valid rpc url (e.g. https://nodes.sequence.app/mainnet)")
	}

	provider, err := ethrpc.NewProvider(fRpc)
	if err != nil {
		return err
	}

	bh, err := provider.BlockNumber(context.Background())
	if err != nil {
		return err
	}

	fmt.Fprintln(cmd.OutOrStdout(), bh)

	return nil
}
