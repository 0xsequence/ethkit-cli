package main

import "errors"

var (
	ErrInvalidBlockInfo = errors.New("invalid block height, tag or hash")
	ErrInvalidRpcUrl = errors.New("invalid rpc url")
	ErrBlockNotFound = errors.New("block not found")
)