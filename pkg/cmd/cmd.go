package cmd

import (
	"fmt"
	"os"

	"github.com/0xsequence/ethkit-cli/pkg/cmd/abigen"
	"github.com/0xsequence/ethkit-cli/pkg/cmd/artifacts"
	"github.com/0xsequence/ethkit-cli/pkg/cmd/balance"
	"github.com/0xsequence/ethkit-cli/pkg/cmd/block"
	"github.com/0xsequence/ethkit-cli/pkg/cmd/tx"
	"github.com/0xsequence/ethkit-cli/pkg/cmd/wallet"

	"github.com/spf13/cobra"
)

var (
	VERSION       = "dev"
	GITBRANCH     = "branch"
	GITCOMMIT     = "last commit"
	GITCOMMITDATE = "last change"

	rootCmd = &cobra.Command{
		Use:     "ethkit",
		Aliases: []string{"ek"},
		Short:   "ethkit - Ethereum dev toolkit",
		Long:    banner(),
		Args:    cobra.MinimumNArgs(1),
		CompletionOptions: cobra.CompletionOptions{
			HiddenDefaultCmd: true,
		},
	}

	versionCmd = &cobra.Command{
		Use:   "version",
		Short: "print the version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("ethkit", version())
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(abigen.NewAbigenCmd())
	rootCmd.AddCommand(artifacts.NewArtifactsCmd())
	rootCmd.AddCommand(balance.NewBalanceCmd())
	rootCmd.AddCommand(block.NewBlockCmd())
	rootCmd.AddCommand(block.NewBlockNumberCmd())
	rootCmd.AddCommand(tx.NewTxCmd())
	rootCmd.AddCommand(wallet.NewWalletCmd())
}

func help(cmd *cobra.Command) {
	fmt.Printf("\n--\n\n")
	cmd.Help()
	os.Exit(0)
}

func version() string {
	if GITBRANCH == "master" {
		return fmt.Sprintf("%s (commit:%s %s)", VERSION, GITCOMMIT, GITCOMMITDATE)
	}
	return fmt.Sprintf("%s (commit:%s %s %s)", VERSION, GITCOMMIT, GITCOMMITDATE, GITBRANCH)
}

func banner() string {
	s := ""
	s += `=====================================================================================` + "\n"
	s += `____________________________/\\\_____________________________________________________` + "\n"
	s += `____________________________\/\\\___________/\\\_____________________________________` + "\n"
	s += `__________________/\\\_______\/\\\__________\/\\\___________/\\\______/\\\___________` + "\n"
	s += `____/\\\\\\\\___/\\\\\\\\\\\__\/\\\__________\/\\\___ /\\___\///____/\\\\\\\\\\\_____` + "\n"
	s += `___/\\\         \////\\\////___\/\\\\\\\\\\___\/\\\_ /\\\___________\////\\\////_____` + "\n"
	s += `___/\\\\\\\\\\\_____\/\\\_______\/\\\/////\\\__\/\\\\\\_______\/\\\_____\/\\\________` + "\n"
	s += `___\//\\\            \/\\\_______\/\\\___\/\\\__\/\\\__\/\\\___\/\\\_____\/\\\_______` + "\n"
	s += `____\//\\\\\\\\\\_____\//\\\\\____\/\\\___\/\\\__\/\\\__\/\\\___\/\\\_____\//\\\\\___` + "\n"
	s += `_____\///////////______\//////_____\///____\///___\///___\///____\///______\/////____` + "\n"
	s += "\n"
	s += "==================================== we <3 Ethereum =================================\n"
	return s
}
