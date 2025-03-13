package main

import (
	"fmt"

	"github.com/aptos-labs/aptos-go-sdk"
	"github.com/aptos-labs/aptos-go-sdk/bcs"
	"github.com/aptos-labs/aptos-go-sdk/internal/util"
)

// const scriptBytes = "a11ceb0b0700000a0601000202020405062107271c0843401083011f0103000207001201020d0e03040f0508000a020a0d0a0e0a030a040a0f0a050a08000a0a080000083c53454c463e5f30046d61696e06537472696e6706737472696e67ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000000000000000000114636f6d70696c6174696f6e5f6d65746164617461090003322e3003322e310000010102"

const scriptBytes = "a11ceb0b0700000a0601000202020405061007161c08324010721f0103000207000b01020d0e03040f0508000a020a0d00083c53454c463e5f30046d61696e06537472696e6706737472696e67ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff000000000000000000000000000000000000000000000000000000000000000114636f6d70696c6174696f6e5f6d65746164617461090003322e3003322e310000010102"
const FundAmount = uint64(100_000_000)

func example(networkConfig aptos.NetworkConfig) {
	// Create a client for Aptos
	client, err := aptos.NewClient(networkConfig)
	if err != nil {
		panic("Failed to create client:" + err.Error())
	}

	// Create a sender  
	sender, err := aptos.NewEd25519Account()
	if err != nil {
		panic("Failed to create sender:" + err.Error())
	}

	// Fund the sender with the faucet to create it on-chain
	println("SENDER: ", sender.Address.String())
	err = client.Fund(sender.Address, FundAmount)
	if err != nil {
		panic("Failed to fund sender:" + err.Error())
	}

	// Now run a script version
	fmt.Printf("\n== Now running script version ==\n")
	runScript(client, sender)

	if err != nil {
		panic("Failed to get store balance:" + err.Error())
	}
	// fmt.Printf("After Script: Receiver Before transfer: %d, after transfer: %d\n", receiverAfterBalance, receiverAfterAfterBalance)
}

func runScript(client *aptos.Client, alice *aptos.Account) {
	scriptBytes, err := aptos.ParseHex(scriptBytes)
	if err != nil {
		panic("Failed to parse script:" + err.Error())
	}

	u128_arg ,err:= util.StrToBigInt("128")
	if err != nil {
		panic("Failed to convert u128:" + err.Error())
	}

	u256_arg ,err:= util.StrToBigInt("256")
	if err != nil {
		panic("Failed to convert u256:" + err.Error())
	}

	vec_u16 := []uint16{1, 2, 3, 4, 5}

	vec_u16_len, err := bcs.SerializeUleb128(uint32(len(vec_u16)));
	if err != nil {
		panic("Failed to serialize uleb128:" + err.Error())
	}

	var vec_u16_arg []byte
	vec_u16_arg = append(vec_u16_arg, vec_u16_len...)
	for _, v := range vec_u16 {
		bytes, err := bcs.SerializeU16(v)
		if err != nil {
			panic("Failed to serialize u16:" + err.Error())
		}
		vec_u16_arg = append(vec_u16_arg, bytes...)
	}

	// 1. Build transaction
	rawTxn, err := client.BuildTransaction(alice.AccountAddress(), aptos.TransactionPayload{
		Payload: &aptos.Script{
			Code:     scriptBytes,
			ArgTypes: []aptos.TypeTag{ },
			Args: []aptos.ScriptArgument{
				{
					Variant: aptos.ScriptArgumentBool,
					Value:   bool(true),
				},
				{
					Variant: aptos.ScriptArgumentU8,
					Value:   uint8(8),
				},
				{
					Variant: aptos.ScriptArgumentU16,
					Value:   uint16(16),
				},
				{
					Variant: aptos.ScriptArgumentU32,
					Value:   uint32(32),
				},
				{
					Variant: aptos.ScriptArgumentU64,
					Value:   uint64(64),
				},
				{
					Variant: aptos.ScriptArgumentU128,
					Value:   *u128_arg,
				},
				{
					Variant: aptos.ScriptArgumentU256,
					Value:   *u256_arg,
				},
				{
					Variant: aptos.ScriptArgumentAddress,
					Value:   alice.Address,
				},
				{
					Variant: aptos.ScriptArgumentU8Vector,
					Value:  []byte("string"),
				},
				{
					Variant: aptos.ScriptArgumentU8Vector,
					Value:  []byte{1, 2, 3, 4, 5},
				},
				{
					Variant: aptos.ScriptArgumentSerialized,
					Value:  &bcs.Serialized{Value: vec_u16_arg},
				},
 
			},
		}})
	if err != nil {
		panic("Failed to build multiagent raw transaction:" + err.Error())
	}

 
	if err != nil {
		panic("Failed to build transaction:" + err.Error())
	}

	// 2. Simulate transaction (optional)
	// This is useful for understanding how much the transaction will cost
	// and to ensure that the transaction is valid before sending it to the network
	// This is optional, but recommended
	simulationResult, err := client.SimulateTransaction(rawTxn, alice)
	if err != nil {
		panic("Failed to simulate transaction:" + err.Error())
	}
	fmt.Printf("\n=== Simulation ===\n")
	fmt.Printf("Gas unit price: %d\n", simulationResult[0].GasUnitPrice)
	fmt.Printf("Gas used: %d\n", simulationResult[0].GasUsed)
	fmt.Printf("Total gas fee: %d\n", simulationResult[0].GasUsed*simulationResult[0].GasUnitPrice)
	fmt.Printf("Status: %s\n", simulationResult[0].VmStatus)

	// 3. Sign transaction
	signedTxn, err := rawTxn.SignedTransaction(alice)
	if err != nil {
		panic("Failed to sign transaction:" + err.Error())
	}

	// 4. Submit transaction
	submitResult, err := client.SubmitTransaction(signedTxn)
	if err != nil {
		panic("Failed to submit transaction:" + err.Error())
	}
	txnHash := submitResult.Hash

	// 5. Wait for the transaction to complete
	_, err = client.WaitForTransaction(txnHash)
	if err != nil {
		panic("Failed to wait for transaction:" + err.Error())
	}
}
func main() {
	example(aptos.DevnetConfig)
}