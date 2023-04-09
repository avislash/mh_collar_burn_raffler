package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/avislash/mh_collar_burn_raffler/client"
	"github.com/avislash/mh_collar_burn_raffler/config"
	"github.com/avislash/mh_collar_burn_raffler/decoder"
	"github.com/avislash/mh_collar_burn_raffler/metadata"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/yaml.v3"
)

var ineligibleMouths map[string]struct{}
var ineligibleFaces map[string]struct{}
var ineligibleTorsos map[string]struct{}
var ineligibleForms map[string]struct{}

var configParams config.Config

func init() {
	ineligibleMouths = make(map[string]struct{})
	ineligibleFaces = make(map[string]struct{})
	ineligibleTorsos = make(map[string]struct{})
	ineligibleForms = make(map[string]struct{})
	configFile, err := ioutil.ReadFile("config.yaml")
	if err != nil {
		panic(fmt.Sprintf("Failed to read in config.yaml: %s", err))
	}

	err = yaml.Unmarshal(configFile, &configParams)
	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal config.yaml: %s", err))
	}

	for _, mouth := range configParams.IneligibleTraits.Mouths {
		ineligibleMouths[mouth] = struct{}{}
	}

	for _, face := range configParams.IneligibleTraits.Faces {
		ineligibleFaces[face] = struct{}{}
	}

	for _, torso := range configParams.IneligibleTraits.Torsos {
		ineligibleTorsos[torso] = struct{}{}
	}

	for _, form := range configParams.IneligibleTraits.Forms {
		ineligibleForms[form] = struct{}{}
	}

}

func main() {
	collarsAddress := "0x354634c4621cDfb7a25E6486cCA1E019777D841B" //TODO: Query directly from proxy contract
	burnedCollars := make(map[string][]uint64)
	eligibleWallets := make(map[string]struct{})
	ineligibleWallets := make(map[string]struct{})
	collarDecoder, err := decoder.NewMutantHoundCollarsDecoder()
	if err != nil {
		panic(err)
	}

	log.Printf("Starting Raffler with configParams: %s", spew.Sdump(configParams))

	etherscanClient := client.NewEtherscanClient("ethereum", "mainnet", os.Getenv("ETHERSCAN_API_KEY"), nil, client.WithRateLimitWithContext(context.Background(), configParams.EtherscanRateLimit))
	txns, err := etherscanClient.QueryEtherscanTransactions(configParams.Snapshot.Start, configParams.Snapshot.Stop, collarsAddress)
	if err != nil {
		log.Fatalf("Error querying Etherscan Txns: %s", err)
	}

	log.Printf("Found %d Txns", len(txns))

	for i, txn := range txns {
		tokenIDs, err := collarDecoder.ParseBurn2RedeemInput(txn.Input)
		if err != nil {
			log.Printf("Error decoding txn #%d (from %s): %s", i, txn.From, err)
			continue
		}
		if len(tokenIDs) == 0 {
			continue
		}
		burnedCollars[txn.From] = tokenIDs
	}

	metadataFetcher := client.NewHoundsMetadataFetcher(configParams.MetadataEndpoint)
	for address, tokenIDs := range burnedCollars {
		houndsMetadata, err := metadataFetcher.FetchMetdata(tokenIDs)
		if err != nil {
			log.Printf("Error fetching metadata for %s: %s", address, err)
			continue
		}

		for _, metadata := range houndsMetadata {
			if houndIsEligible(metadata) {
				eligibleWallets[address] = struct{}{}
			} else {
				ineligibleWallets[address] = struct{}{}
			}
		}
	}

	//Compare Eligible and Ineligble Wallet entries. Remove Ineligible Wallets from the Eligble Wallets
	for ineligibleWallet, _ := range ineligibleWallets {
		delete(eligibleWallets, ineligibleWallet)
	}

	raffleWallets := make([]string, len(eligibleWallets))
	i := 0
	for address, _ := range eligibleWallets {
		raffleWallets[i] = address
		i++
	}

	log.Printf("Found Total Wallets %d", len(eligibleWallets)+len(ineligibleWallets))
	log.Printf("Found %d Ineligible Wallets: %+v", len(ineligibleWallets), ineligibleWallets)
	log.Printf("Found %d Eligible Wallets: %+v", len(eligibleWallets), eligibleWallets)
	log.Printf("Found drawing over %d wallets: %+v", len(raffleWallets), raffleWallets)

	winners := drawWinners(raffleWallets, configParams.MaxWinners)

	log.Printf("Selected %d Winning Entries:", len(winners))
	for i, winner := range winners {
		log.Printf("%d. %s", i+1, winner)
	}

}

func houndIsEligible(metadata metadata.HoundMetadata) bool {
	_, hasIneligibleMouth := ineligibleMouths[metadata.Mouth]
	_, hasIneligibleFace := ineligibleMouths[metadata.Face]
	_, hasIneligibleTorso := ineligibleTorsos[metadata.Torso]
	_, hasIneligibleForm := ineligibleForms[metadata.Form]

	return !(hasIneligibleMouth || hasIneligibleFace || hasIneligibleTorso || hasIneligibleForm)
}

func drawWinners(entries []string, maxWinners uint64) []string {
	rand.Seed(time.Now().UnixNano())
	if maxWinners > uint64(len(entries)) { //Everybody wins
		return entries
	}

	_winners := make(map[string]struct{}) //use map to ensure each person only wins once
	winners := make([]string, maxWinners)
	for {
		if uint64(len(_winners)) == maxWinners {
			break
		}
		index := rand.Uint64() % uint64(len(entries))
		fmt.Println(index)
		_winners[entries[index]] = struct{}{}
	}

	i := 0
	for address, _ := range _winners {
		winners[i] = address
		i++
	}

	return winners
}
