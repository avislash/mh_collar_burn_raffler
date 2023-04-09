package decoder

import (
	"fmt"
	"math/big"
	"strings"

	mhc_abi "github.com/avislash/mh_collar_burn_raffler/abi"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
)

type MutantHoundCollarsDecoder struct {
	abi abi.ABI
}

func NewMutantHoundCollarsDecoder() (*MutantHoundCollarsDecoder, error) {
	abi, err := abi.JSON(strings.NewReader(mhc_abi.MutantHoundCollarsABI))
	if err != nil {
		return nil, fmt.Errorf("Error loading Mutant Hounds Collar ABI: %w", err)
	}
	return &MutantHoundCollarsDecoder{abi}, nil
}

func (mhcd *MutantHoundCollarsDecoder) GetMethod(funcSelector []byte) (*abi.Method, error) {
	method, err := mhcd.abi.MethodById(funcSelector)
	if err != nil {
		return nil, fmt.Errorf("Error decoding function selector %x: %w", funcSelector, err)
	}
	return method, err
}
func (mhcd *MutantHoundCollarsDecoder) IsBurn2Redeem(method *abi.Method) bool {
	burn2RedeemMethodName := "burn2Redeem"
	return nil != method && method.Name == burn2RedeemMethodName
}

func (mhcd *MutantHoundCollarsDecoder) ParseBurn2RedeemInput(txnInput string) ([]uint64, error) {
	inputBytes := common.FromHex(txnInput)

	method, err := mhcd.GetMethod(inputBytes[0:4])
	if err != nil {
		return nil, fmt.Errorf("Error getting method: %w", err)
	}

	if !mhcd.IsBurn2Redeem(method) {
		return nil, nil
	}

	rawTokenIDs, err := method.Inputs.Unpack(inputBytes[4:])
	if err != nil {
		return nil, fmt.Errorf("Error unpacking input args: %w", err)
	}

	if len(rawTokenIDs) == 0 {
		return nil, nil
	}

	tokenIDs := make([]uint64, len(rawTokenIDs[0].([]*big.Int)))
	for i, rawTokenID := range rawTokenIDs[0].([]*big.Int) {
		tokenIDs[i] = rawTokenID.Uint64()
	}

	return tokenIDs, err
}
