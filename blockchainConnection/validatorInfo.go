package blockchainConnection

import (
	"context"
	"log"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
)

var contractAddress = common.HexToAddress("0x7E7A86DAFF5Ba10d830d3581A490446cD40C5c0f")

// get total stack of the hole validator nodes
func GettotalStack() []byte {
	parsedABI := ContractAbi()
	client := Connect()
	callData, err := parsedABI.Pack("getTotalStakedValue")
	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}
	// Call the function.
	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatalf("Failed to call the contract: %v", err)
	}
	return result
}

// get perticular stack of the hole validator node
func GetNodeStakeWithEsg(address common.Address) []byte {
	parsedABI := ContractAbi()
	client := Connect()
	callData, err := parsedABI.Pack("validators", address)
	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}
	// Call the function.
	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatalf("Failed to call the contract: %v", err)
	}
	return result
}

// get total stack of the hole validator nodes
func GetValidatorList() []byte {
	parsedABI := ContractAbi()
	client := Connect()
	callData, err := parsedABI.Pack("getValidatorList")
	callMsg := ethereum.CallMsg{
		To:   &contractAddress,
		Data: callData,
	}
	// Call the function.
	result, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		log.Fatalf("Failed to call the contract: %v", err)
	}

	return result
}
