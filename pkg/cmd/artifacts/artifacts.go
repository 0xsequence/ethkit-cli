package artifacts

import (
	"errors"
	"fmt"

	"github.com/0xsequence/ethkit/ethartifact"
	"github.com/spf13/cobra"
)

func NewArtifactsCmd() *cobra.Command {
	artifacts := &artifacts{}
	cmd := &cobra.Command{
		Use:   "artifacts",
		Short: "Print the contract abi or bytecode from a truffle artifacts file",
		RunE:  artifacts.Run,
	}

	cmd.Flags().String("file", "", "path to truffle contract artifacts file (required)")
	cmd.Flags().Bool("abi", false, "abi")
	cmd.Flags().Bool("bytecode", false, "bytecode")

	return cmd
}

type artifacts struct {
}

func (c *artifacts) Run(cmd *cobra.Command, args []string) error {
	fFile, _ := cmd.Flags().GetString("file")
	fAbi, _ := cmd.Flags().GetBool("abi")
	fBytecode, _ := cmd.Flags().GetBool("bytecode")

	if fFile == "" {
		return errors.New("error: please pass --file")
	}
	if !fAbi && !fBytecode {
		return errors.New("error: please pass either --abi or --bytecode")
	}
	if fAbi && fBytecode {
		return errors.New("error: please pass either --abi or --bytecode, not both")
	}

	artifacts, err := ethartifact.ParseArtifactFile(fFile)
	if err != nil {
		return err
	}

	if fAbi {
		fmt.Println(string(artifacts.ABI))
	}

	if fBytecode {
		fmt.Println(artifacts.Bytecode)
	}

	return nil
}
