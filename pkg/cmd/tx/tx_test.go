package tx_test

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/0xsequence/ethkit-cli/internal"
	"github.com/0xsequence/ethkit-cli/pkg/cmd/tx"
	"github.com/0xsequence/ethkit/go-ethereum/core/types"
	"github.com/stretchr/testify/assert"
)

func execTxCmd(args string) (string, error) {
	cmd := tx.NewTxCmd()
	actual := new(bytes.Buffer)
	cmd.SetOut(actual)
	cmd.SetErr(actual)
	cmd.SetArgs(strings.Split(args, " "))
	if err := cmd.Execute(); err != nil {
		return "", err
	}

	return actual.String(), nil
}

func Test_TxCmd_ValidTransactionHash(t *testing.T) {
	res, err := execTxCmd("0xb11f3c2c5bd49ce4b19b61107dea54c7eecf49e0a0bec88374c066a12b808df8 --rpc-url https://nodes.sequence.app/mainnet")
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func Test_TxCmd_InvalidRpcUrl(t *testing.T) {
	res, err := execTxCmd("0xb11f3c2c5bd49ce4b19b61107dea54c7eecf49e0a0bec88374c066a12b808df8 --rpc-url nodes.sequence.app/mainnet")
	assert.Equal(t, err, internal.ErrInvalidRpcUrl)
	assert.Empty(t, res)
}

func Test_TxCmd_StatusSuccess(t *testing.T) {
	res, err := execTxCmd("0xb11f3c2c5bd49ce4b19b61107dea54c7eecf49e0a0bec88374c066a12b808df8 --rpc-url https://nodes.sequence.app/mainnet")
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res, tx.TxStatusSuccess)
}

func Test_TxCmd_StatusFail(t *testing.T) {
	res, err := execTxCmd("0xe59a336672e73c36f54f62bb248ad93d0b7af239f799723a54e35e1100d74499 --rpc-url https://nodes.sequence.app/mainnet")
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res, tx.TxStatusFail)
}

func Test_TxCmd_WithReceipt(t *testing.T) {
	res, err := execTxCmd("0xb11f3c2c5bd49ce4b19b61107dea54c7eecf49e0a0bec88374c066a12b808df8 --rpc-url https://nodes.sequence.app/mainnet --json --full")
	assert.Nil(t, err)
	assert.NotNil(t, res)
	var txobj *tx.Transaction
	var txrec *types.Receipt
	_ = json.Unmarshal([]byte(res), &txobj)
	_ = json.Unmarshal([]byte(res), &txrec)
	assert.NotNil(t, res, txobj.Receipt)
	assert.Equal(t, txobj.Receipt.TxHash.String(), "0xb11f3c2c5bd49ce4b19b61107dea54c7eecf49e0a0bec88374c066a12b808df8")
}

func Test_TxCmd_TxValidJSON(t *testing.T) {
	res, err := execTxCmd("0xb11f3c2c5bd49ce4b19b61107dea54c7eecf49e0a0bec88374c066a12b808df8 --rpc-url https://nodes.sequence.app/mainnet --json")
	assert.Nil(t, err)
	txobj := tx.Transaction{}
	var p internal.Printable
	_ = p.FromStruct(txobj)
	for k := range p {
		assert.Contains(t, res, k)
	}
}

func Test_TxCmd_BlockValidFieldHash(t *testing.T) {
	// validating also that -f is case-insensitive
	res, err := execTxCmd("0xb11f3c2c5bd49ce4b19b61107dea54c7eecf49e0a0bec88374c066a12b808df8 --rpc-url https://nodes.sequence.app/mainnet -f HASh")
	assert.Nil(t, err)
	assert.Equal(t, res, "0xb11f3c2c5bd49ce4b19b61107dea54c7eecf49e0a0bec88374c066a12b808df8\n")
}

func Test_TxCmd_BlockInvalidField(t *testing.T) {
	res, err := execTxCmd("0xb11f3c2c5bd49ce4b19b61107dea54c7eecf49e0a0bec88374c066a12b808df8 --rpc-url https://nodes.sequence.app/mainnet -f invalid")
	assert.Nil(t, err)
	assert.Equal(t, res, "<nil>\n")
}
