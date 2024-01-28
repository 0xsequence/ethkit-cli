package block_test

import (
	"bytes"
	"strings"
	"testing"

	"github.com/0xsequence/ethkit-cli/pkg/cmd/block"

	"github.com/stretchr/testify/assert"
)

func execBlockNumberCmd(args string) (string, error) {
	cmd := block.NewBlockNumberCmd()
	actual := new(bytes.Buffer)
	cmd.SetOut(actual)
	cmd.SetErr(actual)
	cmd.SetArgs(strings.Split(args, " "))
	if err := cmd.Execute(); err != nil {
		return "", err
	}

	return actual.String(), nil
}

func Test_BlockNumberCmd(t *testing.T) {
	res, err := execBlockNumberCmd("--rpc-url https://nodes.sequence.app/sepolia")
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func Test_BlockNumberCmd_InvalidRPC(t *testing.T) {
	res, err := execBlockNumberCmd("--rpc-url nodes.sequence.app/sepolia")
	assert.NotNil(t, err)
	assert.Empty(t, res)
	assert.Contains(t, err.Error(), "please provide a valid rpc url")
}

func Test_BlockNumberCmd_OneOrMoreArgs(t *testing.T) {
	res, err := execBlockNumberCmd("test --rpc-url nodes.sequence.app/sepolia")
	assert.NotNil(t, err)
	assert.Empty(t, res)
}
