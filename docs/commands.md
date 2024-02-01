# Commands

## wallet

`wallet` handles encrypted Ethereum wallet creation and management in user-supplied keyfiles.
It allows users to create a new Ethereum wallet, import an existing Ethereum wallet from a secret mnemonic, or print an existing wallet's secret mnemonic.

```bash
Usage:
  ethkit-cli wallet [flags]

Flags:
  -h, --help              help for wallet
      --import-mnemonic   import a secret mnemonic to a new keyfile
      --keyfile string    wallet key file path
      --new               create a new wallet and save it to the keyfile
      --print-account     print wallet account address from keyfile (default) (default true)
      --print-mnemonic    print wallet secret mnemonic from keyfile (danger!)
```

## abigen

`abigen` generates Go contract client code from a JSON [truffle](https://www.trufflesuite.com/)
artifacts file.

```bash
Usage:
  ethkit-cli abigen [flags]

Flags:
      --abiFile string         path to abi json file
      --artifactsFile string   path to truffle contract artifacts file
  -h, --help                   help for abigen
      --lang string            target language, supported: [go], default=go
      --outFile string         outFile (optional), default=stdout
      --pkg string             pkg (optional)
      --type string            type (optional)
```

## artifacts

`artifacts` prints the contract ABI or bytecode from a user-supplied truffle artifacts file.

```bash
Usage:
  ethkit-cli artifacts [flags]

Flags:
      --abi           abi
      --bytecode      bytecode
      --file string   path to truffle contract artifacts file (required)
  -h, --help          help for artifacts
```

## balance

`balance` retrieves the balance of an account via RPC by a provided address at a predefined block height.
It provides an implementation of the standard [eth_getBalance](https://ethereum.org/en/developers/docs/apis/json-rpc#eth_getbalance) JSON-RPC method.

```bash
Usage:
  ethkit-cli balance [account] [flags]

Flags:
  -B, --block string     The block height to query at (default "latest")
  -e, --ether            Format the balance in ether
  -h, --help             help for balance
  -r, --rpc-url string   The RPC endpoint to the blockchain node to interact with
```

## block

`block` retrieves a block by a provided block height or tag via RPC.
It provides an implementation of the standard [eth_getBlockByNumber](https://ethereum.org/en/developers/docs/apis/json-rpc#eth_getblockbynumber) JSON-RPC method.

```bash
Usage:
  ethkit-cli block [number|tag] [flags]

Aliases:
  block, bl

Flags:
  -f, --field string     Get the specific field of a block
      --full             Get the full block information
  -h, --help             help for block
  -j, --json             Print the block as JSON

```

## block-number

`block-number` get the latest block number for a given blockchain network.
It provides an implementation of the standard [eth_getBlockNumber](https://ethereum.org/en/developers/docs/apis/json-rpc#eth_blocknumber) JSON-RPC method.

```shell
Usage:
  ethkit block-number [flags]

Aliases:
  block-number, bn

Flags:
  -h, --help             help for block-number
  -r, --rpc-url string   The RPC endpoint to the blockchain node to interact with
```
