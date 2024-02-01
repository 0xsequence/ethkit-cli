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

var (
	// ethereum mainnet
	successTxHash = "0xb11f3c2c5bd49ce4b19b61107dea54c7eecf49e0a0bec88374c066a12b808df8"
	failTxHash    = "0xe59a336672e73c36f54f62bb248ad93d0b7af239f799723a54e35e1100d74499"
	blockHash     = "0x98def9b38c71135124e4e39879580449ee03a62b0740229ecab7fe43af17d7fc"
	blockNumber   = "19121859"
	txIndex       = "81"
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

func Test_TxCmd_ValidTxHash(t *testing.T) {
	res, err := execTxCmd(successTxHash + " --rpc-url https://nodes.sequence.app/mainnet")
	assert.Nil(t, err)
	assert.NotNil(t, res)
}

func Test_TxCmd_ValidBlockHashAndTxIndex(t *testing.T) {
	res, err := execTxCmd(blockHash + " " + txIndex + " --rpc-url https://nodes.sequence.app/mainnet --json")
	var txobj *tx.Transaction
	_ = json.Unmarshal([]byte(res), &txobj)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, txobj.Hash.String(), successTxHash)
}

func Test_TxCmd_ValidBlockHashAndTxIndexInverted(t *testing.T) {
	res, err := execTxCmd(txIndex + " " + blockHash + " --rpc-url https://nodes.sequence.app/mainnet --json")
	var txobj *tx.Transaction
	_ = json.Unmarshal([]byte(res), &txobj)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, txobj.Hash.String(), successTxHash)
}

func Test_TxCmd_ValidBlockNumberAndTxIndex(t *testing.T) {
	res, err := execTxCmd(blockNumber + " " + txIndex + " --rpc-url https://nodes.sequence.app/mainnet --json")
	var txobj *tx.Transaction
	_ = json.Unmarshal([]byte(res), &txobj)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, txobj.Hash.String(), successTxHash)
}

func Test_TxCmd_ValidBlockNumberAndTxIndexInverted(t *testing.T) {
	res, err := execTxCmd(txIndex + " " + blockNumber + " --rpc-url https://nodes.sequence.app/mainnet --json")
	var txobj *tx.Transaction
	_ = json.Unmarshal([]byte(res), &txobj)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Equal(t, txobj.Hash.String(), successTxHash)
}

func Test_TxCmd_AllQueryResultsAreEqual(t *testing.T) {
	byHash, _ := execTxCmd(successTxHash + " --rpc-url https://nodes.sequence.app/mainnet --json")
	byBlockHashAndIndex, _ := execTxCmd(blockHash + " " + txIndex + " --rpc-url https://nodes.sequence.app/mainnet --json")
	byIndexAndBlockHash, _ := execTxCmd(txIndex + " " + blockHash + " --rpc-url https://nodes.sequence.app/mainnet --json")
	byBlockNumberAndIndex, _ := execTxCmd(blockNumber + " " + txIndex + " --rpc-url https://nodes.sequence.app/mainnet --json")
	byIndexAndBlockNumber, _ := execTxCmd(txIndex + " " + blockNumber + " --rpc-url https://nodes.sequence.app/mainnet --json")
	assert.NotNil(t, byHash)
	assert.NotNil(t, byBlockHashAndIndex)
	assert.NotNil(t, byIndexAndBlockHash)
	assert.NotNil(t, byBlockNumberAndIndex)
	assert.NotNil(t, byIndexAndBlockNumber)
	assert.Equal(t, byHash, byBlockHashAndIndex)
	assert.Equal(t, byHash, byIndexAndBlockHash)
	assert.Equal(t, byHash, byBlockNumberAndIndex)
	assert.Equal(t, byHash, byIndexAndBlockNumber)
}

func Test_TxCmd_InvalidTxHash(t *testing.T) {
	res, err := execTxCmd(successTxHash[:len(successTxHash)-1] + " --rpc-url https://nodes.sequence.app/mainnet")
	assert.Equal(t, err, internal.ErrInvalidHash)
	assert.Empty(t, res)
}

func Test_TxCmd_InvalidArgsByInvalidHash(t *testing.T) {
	res, err := execTxCmd(blockHash[:len(blockHash)-1] + " " + txIndex + " --rpc-url https://nodes.sequence.app/mainnet --json")
	assert.Equal(t, err, internal.ErrCmdInvalidArgs)
	assert.Empty(t, res)
}

func Test_TxCmd_WrongTxByInvalidTxIndex(t *testing.T) {
	res, err := execTxCmd(blockNumber + " " + txIndex[:len(txIndex)-1] + " --rpc-url https://nodes.sequence.app/mainnet --json")
	var txobj *tx.Transaction
	_ = json.Unmarshal([]byte(res), &txobj)
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.NotEqual(t, txobj.Hash.String(), successTxHash)
}

func Test_TxCmd_TxNotFoundByWrongBlockNumber(t *testing.T) {
	res, err := execTxCmd(blockNumber[:len(blockNumber)-1] + " " + txIndex + " --rpc-url https://nodes.sequence.app/mainnet --json")
	assert.Equal(t, err, internal.ErrTxNotFound)
	assert.Empty(t, res)
}

func Test_TxCmd_TxNotFoundByWrongHash(t *testing.T) {
	res, err := execTxCmd(blockHash + " --rpc-url https://nodes.sequence.app/mainnet")
	assert.Equal(t, err, internal.ErrTxNotFound)
	assert.Empty(t, res)
}

func Test_TxCmd_InvalidRpcUrl(t *testing.T) {
	res, err := execTxCmd(successTxHash + " --rpc-url nodes.sequence.app/mainnet")
	assert.Equal(t, err, internal.ErrInvalidRpcUrl)
	assert.Empty(t, res)
}

func Test_TxCmd_StatusSuccess(t *testing.T) {
	res, err := execTxCmd(successTxHash + " --rpc-url https://nodes.sequence.app/mainnet")
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res, tx.TxStatusSuccess)
}

func Test_TxCmd_StatusFail(t *testing.T) {
	res, err := execTxCmd(failTxHash + " --rpc-url https://nodes.sequence.app/mainnet")
	assert.Nil(t, err)
	assert.NotNil(t, res)
	assert.Contains(t, res, tx.TxStatusFail)
}

func Test_TxCmd_WithReceipt(t *testing.T) {
	res, err := execTxCmd(successTxHash + " --rpc-url https://nodes.sequence.app/mainnet --json --full")
	assert.Nil(t, err)
	assert.NotNil(t, res)
	var txobj *tx.Transaction
	var txrec *types.Receipt
	_ = json.Unmarshal([]byte(res), &txobj)
	_ = json.Unmarshal([]byte(res), &txrec)
	assert.NotNil(t, res, txobj.Receipt)
	assert.Equal(t, txobj.Receipt.TxHash.String(), successTxHash)
}

func Test_TxCmd_TxValidJSON(t *testing.T) {
	res, err := execTxCmd(successTxHash + " --rpc-url https://nodes.sequence.app/mainnet --json")
	assert.Nil(t, err)
	txobj := tx.Transaction{}
	var p internal.Printable
	_ = p.FromStruct(txobj)
	for k := range p {
		assert.Contains(t, res, k)
	}
}

func Test_TxCmd_TxValidFieldHash(t *testing.T) {
	// validating also that -f is case-insensitive
	res, err := execTxCmd(successTxHash + " --rpc-url https://nodes.sequence.app/mainnet -f HASh")
	assert.Nil(t, err)
	assert.Equal(t, res, successTxHash+"\n")
}

func Test_TxCmd_TxInvalidField(t *testing.T) {
	res, err := execTxCmd(successTxHash + " --rpc-url https://nodes.sequence.app/mainnet -f invalid")
	assert.Nil(t, err)
	assert.Equal(t, res, "<nil>\n")
}
