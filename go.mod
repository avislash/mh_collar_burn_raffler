module github.com/avislash/mh_collar_burn_raffler

go 1.20

require (
	github.com/davecgh/go-spew v1.1.1
	github.com/ethereum/go-ethereum v1.11.5
	github.com/nanmu42/etherscan-api v1.10.0
	golang.org/x/time v0.3.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/btcsuite/btcd/btcec/v2 v2.2.0 // indirect
	github.com/decred/dcrd/dcrec/secp256k1/v4 v4.0.1 // indirect
	github.com/holiman/uint256 v1.2.0 // indirect
	golang.org/x/crypto v0.1.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
)

replace github.com/nanmu42/etherscan-api v1.10.0 => github.com/avislash/etherscan-api v0.2.0
