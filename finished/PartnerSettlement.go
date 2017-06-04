/*
This is the main package for interacting with the blockchain for the sample use case
*/

package main

import (
	"errors"
	"fmt"
	"encoding/json"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)

// SimpleChaincode example simple Chaincode implementation
type SimpleChaincode struct {
}

type CardMember struct {
	CardMemberNumber 	string
	LoyaltyAccountNumber 	string
}

type InputTransaction struct {
	ROCReferenceNumber 	string
	SENumber 	   	string
	TransactionAmount	float32
	BGCId			string
	ConversionRateInfo	ConversionRateInfo
}

type ConversionRateInfo struct {
	ConversionRate		float32
	ISOCurrency 		string
	ConversionSlabNumber	int
	ConversionSlabLowerThreshold	float32
	ConversionSlabUpperThreshold 	float32
}

type SettlementInfo struct {
	SummarySettlementId 	string
	LoyaltySettlementAmount	float32
	ISOCurrencyLoyalty	string
	SettlementAmount	float32
	ISOCurrency		string
}

func (t *SimpleChaincode) Init(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {

	fmt.Println("Init function is getting fired")
	var err error

	var blank []string

	var blankBytes, _ = json.Marshal(&blank)
	err = stub.PutState("abc", blankBytes)
	if err != nil {
		return nil, err
	}

	fmt.Println("Initialization complete")
	return nil, nil

}

// Transaction to invoke the blockchain
func (t SimpleChaincode) Invoke(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	fmt.Println("The function currently running is invoke")

	if function == "createInputTransaction" {
		return t.createInputTransaction(stub, args)
	}

	if function == "updateSettlementSummary" {
		return t.updateSettlementSummary(stub, args)
	}

	if function == "deleteTransactionEntry" {
		return t.deleteTransaction(stub, args)
	}

	return nil, nil

}

// Deletes a transaction from the chain
func (t *SimpleChaincode) deleteTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Inside delete function")
	var transactionID string
	if len(args) != 1 {
		return nil, errors.New("Incorrect number of arguments. Expecting 1")
	}

	transactionID = args[0]

	// Delete the key from the state in ledger
	err := stub.DelState(transactionID)
	if err != nil {
		return nil, errors.New("Failed to delete transaction")
	}

	return nil, nil
}

func (t SimpleChaincode) createInputTransaction(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Creating a transaction")

	var err error
	if len(args) != 1 {
		fmt.Println("Error getting transaction id")
		return nil, errors.New("createTransaction needs one arguement")
	}
	var inputTransaction InputTransaction

	err = json.Unmarshal([]byte(args[0]),&inputTransaction)
	if err != nil {
		fmt.Println("Error unmarshalling transaction details")
		return nil, errors.New("Error unmarshalling transaction")
	}

	transactionBytes, err := stub.GetState(inputTransaction.ROCReferenceNumber)

	fmt.Println(transactionBytes)

	if err != nil {
		fmt.Println("There is no entry present for " + inputTransaction.ROCReferenceNumber)
	}

	if transactionBytes == nil {
		fmt.Println("No data found in blockchain")
		fmt.Println(inputTransaction)
		inputTransactionAsString, err := json.Marshal(&inputTransaction)
		fmt.Println(inputTransactionAsString)
		err = stub.PutState(inputTransaction.ROCReferenceNumber, []byte(inputTransactionAsString))
		if err != nil {
			fmt.Println("Error writing into the blockchain")
		}

		fmt.Println("Data written to the chain")
		fmt.Println(inputTransaction.ConversionRateInfo.ConversionSlabUpperThreshold)

		/*
		transactionBytesNew, err := stub.GetState(inputTransaction.ROCReferenceNumber)
		if err != nil {
			fmt.Println("Value check")
		}
		json.Unmarshal(transactionBytesNew, &finalInputTransaction)
		fmt.Println(finalInputTransaction.SENumber)
		*/
		return nil, nil
	}

	err = stub.PutState(inputTransaction.ROCReferenceNumber, transactionBytes)
	if err != nil {
		fmt.Println("Error writing into the blockchain")
	}

	return nil,nil
}

func (t SimpleChaincode) updateSettlementSummary(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Updating settlement summary")

	var err error
	var settlementInfo SettlementInfo
	var loyaltySettlementAmount float32
	var settlementAmount float32
	var prevState SettlementInfo
	var newState SettlementInfo

	if len(args) !=1 {
		fmt.Println("Error getting settlement summary")
		return nil, errors.New("SettlementId is not passed")
	}

	err = json.Unmarshal([]byte(args[0]), &settlementInfo)

	if err != nil {
		fmt.Println("Error unmarshalling the data")
	}

	loyaltySettlementAmount = settlementInfo.LoyaltySettlementAmount
	settlementAmount = settlementInfo.SettlementAmount
	settlementSummaryBytes, err := stub.GetState(settlementInfo.SummarySettlementId)

	if err != nil {
		fmt.Println("There is no entry present for " + settlementInfo.SummarySettlementId)
	}

	if settlementSummaryBytes == nil {
		fmt.Println("No data found for summary query in chain")
		newState.SummarySettlementId 		= settlementInfo.SummarySettlementId
		newState.LoyaltySettlementAmount	= loyaltySettlementAmount
		newState.SettlementAmount		= settlementAmount
		newState.ISOCurrency			= settlementInfo.ISOCurrency
		newState.ISOCurrencyLoyalty		= settlementInfo.ISOCurrencyLoyalty
	} else {
		json.Unmarshal(settlementSummaryBytes,&prevState)
		loyaltySettlementAmount = loyaltySettlementAmount + prevState.LoyaltySettlementAmount
		settlementAmount	= settlementAmount	  + prevState.SettlementAmount
		newState.SummarySettlementId 		= prevState.SummarySettlementId
		newState.LoyaltySettlementAmount	= loyaltySettlementAmount
		newState.SettlementAmount		= settlementAmount
		newState.ISOCurrency			= prevState.ISOCurrency
		newState.ISOCurrencyLoyalty		= prevState.ISOCurrencyLoyalty
	}


	fmt.Println(newState.SummarySettlementId)
	fmt.Println(newState.ISOCurrency)
	fmt.Println(newState.ISOCurrencyLoyalty)
	fmt.Println(loyaltySettlementAmount)
	fmt.Println(settlementAmount)

	jsonAsBytes, _ := json.Marshal(newState)
	err = stub.PutState(settlementInfo.SummarySettlementId, jsonAsBytes)

	return nil, nil
}

// Query callback representing the query of a chaincode
func (t *SimpleChaincode) Query(stub shim.ChaincodeStubInterface, function string, args []string) ([]byte, error) {
	if function == "getTransactionDetail" {
		return t.getTransactionDetail(stub, args)
	}

	if function == "getSettlementSummary" {
		return t.getSettlementSummary(stub, args)
	}

	return nil,nil
}


func (t *SimpleChaincode) getTransactionDetail(stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {

	var transactionID string
	var inputTransaction InputTransaction
	var err error

	if len(args) != 1  {
		fmt.Println("Expecting transaction id for the inquiry")
		return nil, errors.New("Error in number of arguements passed")
	}

	transactionID = args[0]

	inputTransactionAsBytes, err := stub.GetState(transactionID)

	if err != nil {
		fmt.Println("No information found for this transaction id")
		return nil, nil
	}

	json.Unmarshal(inputTransactionAsBytes,&inputTransaction)
	
	jsonResp := "{\""  + string(inputTransactionAsBytes) + "\"}"
	fmt.Printf("Query Response for detail:%s\n", jsonResp)
	
	return inputTransactionAsBytes,nil
}

func (t *SimpleChaincode) getSettlementSummary (stub shim.ChaincodeStubInterface, args []string) ([]byte, error) {
	fmt.Println("Inside Get Summary")

	var summarySettlementID string
	var settlementInfo SettlementInfo
	var err error


	if len(args) !=1  {
		fmt.Println("Expecting transaction id for the inquiry")
	}

	summarySettlementID = args[0]

	settlementInfoAsBytes, err := stub.GetState(summarySettlementID)

	if err != nil {
		fmt.Println("No information found for this summary settlement id")
		return nil, nil
	}

	json.Unmarshal(settlementInfoAsBytes,&settlementInfo)
	jsonResp := "{\""  + string(settlementInfoAsBytes) + "\"}"
	fmt.Printf("Query Response for summary:%s\n", jsonResp) 

	return settlementInfoAsBytes, nil
}

func main() {
	err := shim.Start(new(SimpleChaincode))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
