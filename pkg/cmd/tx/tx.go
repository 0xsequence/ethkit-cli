package tx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/url"
	"strconv"

	"github.com/0xsequence/ethkit-cli/internal"
	"github.com/0xsequence/ethkit/go-ethereum"
	"github.com/0xsequence/ethkit/go-ethereum/common"
	"github.com/0xsequence/ethkit/go-ethereum/common/hexutil"
	"github.com/0xsequence/ethkit/go-ethereum/core/types"
	"github.com/0xsequence/ethkit/go-ethereum/ethclient"
	"github.com/0xsequence/ethkit/go-ethereum/rpc"

	"github.com/spf13/cobra"
)

const (
	flagTxRaw    = "raw"
	flagTxFull   = "full"
	flagTxField  = "field"
	flagTxRpcUrl = "rpc-url"
	flagTxJson   = "json"
)

type tx struct {
}

// NewTxCommand returns a new build command to retrieve a transaction.
func NewTxCmd() *cobra.Command {
	c := &tx{}
	cmd := &cobra.Command{
		Use:     "tx [hash] | [blockNumber] [txIndex] | [blockHash] [txIndex]",
		Short:   "Get the information about the transaction",
		Aliases: []string{"t"},
		Args:    cobra.RangeArgs(1, 2),
		RunE:    c.Run,
	}

	cmd.Flags().BoolP(flagTxRaw, "", false, "Print the raw RLP encoded transaction")
	cmd.Flags().BoolP(flagTxFull, "", false, "Print extensive transaction info with receipt")
	cmd.Flags().StringP(flagTxField, "f", "", "Get the specific field of a block")
	cmd.Flags().StringP(flagTxRpcUrl, "r", "", "The RPC endpoint to the blockchain node to interact with")
	cmd.Flags().BoolP(flagTxJson, "j", false, "Print the transaction as JSON")

	return cmd
}

// Run executes the command
func (c *tx) Run(cmd *cobra.Command, args []string) error {
	fRaw, err := cmd.Flags().GetBool(flagTxRaw)
	if err != nil {
		return err
	}
	fFull, err := cmd.Flags().GetBool(flagTxFull)
	if err != nil {
		return err
	}
	fField, err := cmd.Flags().GetString(flagTxField)
	if err != nil {
		return err
	}
	fRpc, err := cmd.Flags().GetString(flagTxRpcUrl)
	if err != nil {
		return err
	}
	fJson, err := cmd.Flags().GetBool(flagTxJson)
	if err != nil {
		return err
	}

	if _, err = url.ParseRequestURI(fRpc); err != nil {
		return internal.ErrInvalidRpcUrl
	}

	rpcClient, err := rpc.Dial(fRpc)
	if err != nil {
		return err
	}
	ethClient, err := ethclient.Dial(fRpc)
	if err != nil {
		return err
	}

	var rawtx *rpcTransaction
	if len(cmd.Flags().Args()) == 1 {
		rawtx, err = RawTransactionByHash(rpcClient, context.Background(), cmd.Flags().Args()[0])
		if err != nil {
			return err
		}
	} else {
		cmdArgs, err := setCmdArgs(cmd.Flags().Args())
		if err != nil {
			return err
		}

		if cmdArgs.blockHash == nil {
			rawtx, err = TransactionByBlockNumberAndIndex(rpcClient, context.Background(), *cmdArgs.blockNumber, *cmdArgs.index)
			if err != nil {
				return err
			}
		} else {
			rawtx, err = TransactionByBlockHashAndIndex(rpcClient, context.Background(), *cmdArgs.blockHash, *cmdArgs.index)
			if err != nil {
				return err
			}
		}
	}

	tx := NewTransaction(rawtx)

	receipt, err := tx.TransactionReceipt(ethClient, tx.Hash)
	if err != nil {
		return err
	}
	tx.SetTxStatus(receipt.Status)
	tx.SetTxIndex((uint64(receipt.TransactionIndex)))

	if fFull {
		tx.WithReceipt(receipt)
	}

	if fField != "" {
		// TODO: Not working for all fields
		fmt.Fprintln(cmd.OutOrStdout(), internal.GetValueByJSONTag(tx, fField))
	} else if fRaw {
		buf := new(bytes.Buffer)
		if err := rawtx.tx.EncodeRLP(buf); err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), common.Bytes2Hex(buf.Bytes()))
	} else if fJson {
		json, err := internal.PrettyJSON(tx)
		if err != nil {
			return err
		}
		fmt.Fprintln(cmd.OutOrStdout(), *json)
	} else {
		fmt.Fprintln(cmd.OutOrStdout(), tx)
	}

	return nil
}

// TxStatus defines the status of a transaction
type TxStatus string

const (
	// TxStatusPending refers to transactions no yet finalized in a block and that are still in the mempool
	TxStatusPending TxStatus = "TX_PENDING"
	// TxStatusFail refers to finalized transactions which were not valid or failed for some reason (e.g., not enough gas)
	TxStatusFail TxStatus = "TX_FAIL"
	// TxStatusSuccess refers to finalized and valid transactions
	TxStatusSuccess TxStatus = "TX_SUCCESS"
)

// Transaction is a customized transaction object for cli
type Transaction struct {
	Hash        common.Hash  `json:"hash"`
	Status      TxStatus     `json:"status"`
	BlockHash   *common.Hash `json:"blockHash"`
	BlockNumber *big.Int     `json:"blockNumber"`
	// Time        uint64         `json:"timestamp"`
	From      common.Address `json:"from"`
	To        common.Address `json:"to"`
	Value     *big.Int       `json:"value"`
	Gas       uint64         `json:"gas"`
	GasPrice  *big.Int       `json:"gasPrice"`
	GasTipCap *big.Int       `json:"gasTipCap"`
	GasFeeCap *big.Int       `json:"gasFeeCap"`
	Nonce     uint64         `json:"nonce"`
	Index     uint64         `json:"positionInBlock"`
	Type      byte           `json:"type"`
	Data      []byte         `json:"data"`
	Receipt   *types.Receipt `json:"receipt,omitempty"`
	V         *big.Int       `json:"v"`
	R         *big.Int       `json:"r"`
	S         *big.Int       `json:"s"`
}

// NewTransaction returns the custom-built Transaction object.
func NewTransaction(rawtx *rpcTransaction) *Transaction {
	return &Transaction{
		Hash:        rawtx.tx.Hash(),
		BlockHash:   rawtx.BlockHash,
		BlockNumber: internal.HexToBigInt(rawtx.BlockNumber),
		// TODO types.Transaction.Time() is missing in the current ethkit go-ethereum package. Either upgrade ethkit or replace import with official go-ethereum module
		// Time:        uint64(rawtx.tx.Time().Unix()),
		From:      *rawtx.From,
		To:        *rawtx.tx.To(),
		Value:     rawtx.tx.Value(),
		Gas:       rawtx.tx.Gas(),
		GasPrice:  rawtx.tx.GasPrice(),
		GasTipCap: rawtx.tx.GasTipCap(),
		GasFeeCap: rawtx.tx.GasFeeCap(),
		Nonce:     rawtx.tx.Nonce(),
		Type:      rawtx.tx.Type(),
		Data:      rawtx.tx.Data(),
		V:         signatures(rawtx.tx)[0],
		R:         signatures(rawtx.tx)[1],
		S:         signatures(rawtx.tx)[2],
	}
}

// String overrides the standard behavior for Transaction "to-string".
func (tx *Transaction) String() string {
	var p internal.Printable
	if err := p.FromStruct(tx); err != nil {
		log.Fatal(err)
	}
	s := p.Columnize(*internal.NewPrintableFormat(20, 0, 0, byte(' ')))

	return s
}

// WithReceipt sets the transaction receipt object into the transaction
func (tx *Transaction) WithReceipt(receipt *types.Receipt) *Transaction {
	tx.Receipt = receipt
	return tx
}

// TransactionReceipt retrieves the transaction receipt by its hash via eth client
func (tx *Transaction) TransactionReceipt(client *ethclient.Client, txHash common.Hash) (*types.Receipt, error) {
	receipt, err := client.TransactionReceipt(context.Background(), txHash)
	if err == ethereum.NotFound {
		return nil, nil
	}

	return receipt, err
}

// SetTxStatus sets the status into the transaction
func (tx *Transaction) SetTxStatus(status uint64) {
	if tx.BlockNumber == nil {
		tx.Status = TxStatusPending
	} else if status == 0 {
		tx.Status = TxStatusFail
	} else {
		tx.Status = TxStatusSuccess
	}
}

// SetTxIndex sets the index or position in block into the transaction
func (tx *Transaction) SetTxIndex(index uint64) { tx.Index = index }

// UnmarshalJSON parses the JSON-encoded data of msg into the rpc transaction object
func (tx *rpcTransaction) UnmarshalJSON(msg []byte) error {
	if err := json.Unmarshal(msg, &tx.tx); err != nil {
		return err
	}
	return json.Unmarshal(msg, &tx.txExtraInfo)
}

// RawTransactionByHash retrieves an rpc transaction by its hash via rpc client using the standard eth_getTransactionByHash ethereum JSON-RPC call
func RawTransactionByHash(client *rpc.Client, ctx context.Context, hash string) (*rpcTransaction, error) {
	if !internal.IsValidHash(hash) {
		return nil, internal.ErrInvalidHash
	}

	return callContext(client, ctx, "eth_getTransactionByHash", common.HexToHash(hash))
}

// TransactionByBlockHashAndIndex retrieves an rpc transaction via rpc client by the hash of the block is contained in and the position (index) in that block using the standard eth_getTransactionByBlockHashAndIndex ethereum JSON-RPC call
func TransactionByBlockHashAndIndex(client *rpc.Client, ctx context.Context, blockHash, index string) (*rpcTransaction, error) {
	if !(internal.IsValidHash(blockHash) && internal.IsInt(index)) {
		return nil, internal.ErrCmdInvalidArgs
	}
	i, _ := strconv.Atoi(index)

	return callContext(client, ctx, "eth_getTransactionByBlockHashAndIndex", common.HexToHash(blockHash), hexutil.Uint(i))
}

// TransactionByBlockNumberAndIndex retrieves an rpc transaction via rpc client by the number of the block is contained in and the position (index) in that block using the standard eth_getTransactionByBlockNumberAndIndex ethereum JSON-RPC call
func TransactionByBlockNumberAndIndex(client *rpc.Client, ctx context.Context, blockNumber, index string) (*rpcTransaction, error) {
	if !(internal.IsInt(blockNumber) && internal.IsInt(index)) {
		return nil, internal.ErrCmdInvalidArgs
	}
	bn, _ := strconv.Atoi(blockNumber)
	i, _ := strconv.Atoi(index)

	return callContext(client, ctx, "eth_getTransactionByBlockNumberAndIndex", hexutil.Uint(i), hexutil.Uint(bn))
}

func callContext(client *rpc.Client, ctx context.Context, method string, args ...any) (*rpcTransaction, error) {
	var json *rpcTransaction
	err := client.CallContext(ctx, &json, method, args...)
	if err != nil {
		return nil, err
	} else if json == nil {
		return nil, internal.ErrTxNotFound
	} else if _, r, _ := json.tx.RawSignatureValues(); r == nil {
		return nil, internal.ErrTxWithoutSignature
	}

	return json, nil
}

type cmdArgs struct {
	blockHash, blockNumber, index *string
}

func setCmdArgs(args []string) (*cmdArgs, error) {
	if internal.IsValidHash(args[0]) && internal.IsInt(args[1]) {
		return &cmdArgs{blockHash: &args[0], index: &args[1]}, nil
	}
	if internal.IsInt(args[0]) {
		if internal.IsInt(args[1]) {
			// assumption: blockNumber > txIndex
			if args[0] > args[1] {
				return &cmdArgs{blockNumber: &args[0], index: &args[1]}, nil
			}
			return &cmdArgs{blockNumber: &args[1], index: &args[0]}, nil
		}
		if internal.IsValidHash(args[1]) {
			return &cmdArgs{blockHash: &args[1], index: &args[0]}, nil
		}
	}

	return nil, internal.ErrCmdInvalidArgs
}

func signatures(tx *types.Transaction) [3]*big.Int {
	v, r, s := tx.RawSignatureValues()
	return [3]*big.Int{v, r, s}
}

type rpcTransaction struct {
	tx *types.Transaction
	txExtraInfo
}

type txExtraInfo struct {
	BlockNumber *string
	BlockHash   *common.Hash
	From        *common.Address
}
