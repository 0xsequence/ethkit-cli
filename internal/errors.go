package internal

import "errors"

var (
	ErrInvalidBlockInfo   = errors.New("invalid block height, tag or hash")
	ErrInvalidRpcUrl      = errors.New("invalid rpc url")
	ErrInvalidHash        = errors.New("invalid hash")
	ErrBlockNotFound      = errors.New("block not found")
	ErrTxNotFound         = errors.New("transaction not found")
	ErrCmdInvalidArgs     = errors.New("invalid arguments")
	ErrTxWithoutSignature = errors.New("server returned transaction without signature")
)
