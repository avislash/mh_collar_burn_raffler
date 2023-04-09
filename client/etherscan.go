package client

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/nanmu42/etherscan-api"
	"golang.org/x/time/rate"
)

type Network string

const (
	Mainnet         Network = "mainnet"
	Ethereum        Network = "ethereum"
	Ropsten         Network = "ropsten"
	Kovan           Network = "kovan"
	Rinkeby         Network = "rinkeby"
	Goerli          Network = "goerli"
	Tobalaba        Network = "tobalaba"
	Polygon         Network = "polygon"
	Mumbai          Network = "mumbai"
	BSC             Network = "bsc"
	BSCTestnet      Network = "bsc-testnet"
	Arbitrum        Network = "arbitrum"
	ArbitrumTestnet Network = "arbitrum-testnet"
)

func (n Network) IsValid() bool {
	switch n {
	case Mainnet, Ropsten, Kovan, Rinkeby, Goerli, Tobalaba, Mumbai:
		return true
	default:
		return false
	}
}

func (n Network) ToEtherscanAPINetwork() etherscan.Network {
	switch n {
	case Ropsten:
		return etherscan.Ropsten
	case Kovan:
		return etherscan.Kovan
	case Rinkeby:
		return etherscan.Rinkby
	case ArbitrumTestnet, Goerli: //Arbiturm Testnet and Ethereum Testnet use the same EtherscanAPINetwork Prefix
		return etherscan.Goerli
	case Tobalaba:
		return etherscan.Tobalaba
	case Mumbai:
		return "api-mumbai"
	case BSCTestnet:
		return "api-testnet"
	default:
		return etherscan.Mainnet
	}
}

func (n Network) ToEtherscanBaseEndPoint(blockchain string) string {
	switch Network(blockchain) {
	case Mainnet, Ethereum, Ropsten, Kovan, Rinkeby, Goerli, Tobalaba:
		return fmt.Sprintf(`https://%s.etherscan.io/api?`, n.ToEtherscanAPINetwork().SubDomain())
	case Arbitrum, ArbitrumTestnet:
		return fmt.Sprintf(`https://%s.arbiscan.io/api?`, n.ToEtherscanAPINetwork().SubDomain())
	case Polygon, Mumbai:
		return fmt.Sprintf(`https://%s.polygonscan.com/api?`, n.ToEtherscanAPINetwork().SubDomain())
	case BSC, BSCTestnet:
		return fmt.Sprintf(`https://%s.bscscan.com/api?`, n.ToEtherscanAPINetwork().SubDomain())
	default:
		return ""
	}
}

func (n Network) String() string {
	return string(n)
}

type CustomizationOptions func(config *etherscan.Customization)

func WithRateLimitWithContext(ctx context.Context, txnPerSecond int) func(config *etherscan.Customization) {
	return func(config *etherscan.Customization) {
		rateLimiter := rate.NewLimiter(rate.Every(time.Second), txnPerSecond)
		rateLimit := func(_, _ string, _ map[string]interface{}) error {
			if err := rateLimiter.Wait(ctx); err != nil {
				return fmt.Errorf("Rate Limiter Error: %w", err)
			}
			return nil
		}
		config.BeforeRequest = rateLimit
	}

}

type EtherscanClient struct {
	*etherscan.Client
}

func WithRateLimit(txnPerSecond int) func(config *etherscan.Customization) {
	return WithRateLimitWithContext(context.Background(), txnPerSecond)
}

func NewEtherscanClient(blockchain string, network Network, key string, httpClient *http.Client, options ...CustomizationOptions) *EtherscanClient {
	clientConfig := etherscan.Customization{
		Key:     key,
		BaseURL: network.ToEtherscanBaseEndPoint(blockchain),
		Client:  httpClient,
	}

	if nil != httpClient {
		clientConfig.Timeout = clientConfig.Client.Timeout
	}

	for _, applyOption := range options {
		applyOption(&clientConfig)
	}

	return &EtherscanClient{etherscan.NewCustomized(clientConfig)}
}

func (ec *EtherscanClient) QueryEtherscanTransactions(start, end time.Time, address string) ([]etherscan.NormalTx, error) {
	startBlock, err := ec.BlockNumber(start.Unix(), "before")
	if err != nil {
		return []etherscan.NormalTx{}, fmt.Errorf("Unable to query Starting Block Number from Etherscan: %w", err)
	}

	stopBlock, err := ec.BlockNumber(end.Unix(), "before")
	if err != nil {
		return []etherscan.NormalTx{}, fmt.Errorf("Unable to query Stopping Block Number from Etherscan: %w", err)
	}

	log.Printf("Querying for Txns between blocks %d and %d", startBlock, stopBlock)

	txns, err := ec.NormalTxByAddress(address, &startBlock, &stopBlock, 0, 0, false)
	if err != nil {
		if strings.Contains(err.Error(), "No transactions found") {
			//No transactions over the queried range is a normal outcome. No need to flag an error
			return []etherscan.NormalTx{}, nil
		}
		return []etherscan.NormalTx{}, fmt.Errorf("Unable to get Normal Txns for %s from Etherscan: %w", address, err)
	}
	return txns, nil
}
