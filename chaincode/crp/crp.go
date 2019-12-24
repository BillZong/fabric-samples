package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"strconv"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type apiName string

const (
	apiSupportAPIs    apiName = "apis"
	apiSupportOpCodes apiName = "op-codes"
	apiInitLedger     apiName = "init-ledger"
	apiCommit         apiName = "commit"
)

type opCode string

const (
	opCreate opCode = "create"
	opDelete opCode = "delete"
	opReview opCode = "review"
	opModify opCode = "modify"
)

const (
	apiKeyPrefix    = "API"
	opCodeKeyPrefix = "OPCode"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// EDocument is the e-document structure for CRP(China Resource Power).
type EDocument struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	UID         string `json:"uid"`
	OpCode      string `json:"op-code"`
	OpTime      string `json:"op-time"`
	Description string `json:"description,omitempty"`
	Hash        string `json:"hash"`
}

func (s *SmartContract) Init(APIstub shim.ChaincodeStubInterface) sc.Response {
	return shim.Success(nil)
}

func (s *SmartContract) Invoke(APIstub shim.ChaincodeStubInterface) sc.Response {
	function, args := APIstub.GetFunctionAndParameters()

	api := apiName(function)

	if api == apiInitLedger {
		return s.initLedger(APIstub)
	} else if api == apiSupportAPIs {
		return s.queryAllApis(APIstub)
	} else if api == apiSupportOpCodes {
		return s.queryAllOpCodes(APIstub)
	}

	return shim.Error(fmt.Sprintf("Not support function name: %s, args: %s", function, args))
}

// func (s *SmartContract) queryCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

//  if len(args) != 1 {
// 	 return shim.Error("Incorrect number of arguments. Expecting 1")
//  }

//  carAsBytes, _ := APIstub.GetState(args[0])
//  return shim.Success(carAsBytes)
// }

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	funcName, _, _, _ := runtime.Caller(0)
	fmt.Printf("------ %s ------\n", funcName)

	apiNames := []apiName{
		apiSupportAPIs,
		apiSupportOpCodes,
		apiInitLedger,
		apiCommit,
	}

	for i, v := range apiNames {
		fmt.Println("i is ", i)
		apiNameAsBytes, _ := json.Marshal(v)
		APIstub.PutState(apiKeyPrefix+strconv.Itoa(i), apiNameAsBytes)
		fmt.Println("Added", v)
	}

	opCodes := []opCode{
		opCreate,
		opDelete,
		opReview,
		opModify,
	}

	for i, v := range opCodes {
		fmt.Println("i is ", i)
		asBytes, _ := json.Marshal(v)
		APIstub.PutState("OPCode"+strconv.Itoa(i), asBytes)
		fmt.Println("Added", v)
	}

	return shim.Success(nil)
}

// func (s *SmartContract) createCar(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

//  if len(args) != 5 {
// 	 return shim.Error("Incorrect number of arguments. Expecting 5")
//  }

//  var car = Car{Make: args[1], Model: args[2], Colour: args[3], Owner: args[4]}

//  carAsBytes, _ := json.Marshal(car)
//  APIstub.PutState(args[0], carAsBytes)

//  return shim.Success(nil)
// }

func (s *SmartContract) queryArrayByIteratorKeys(keys shim.StateQueryIteratorInterface) ([]byte, error) {
	// buffer is a JSON array containing QueryResults
	var buffer bytes.Buffer
	buffer.WriteString("[")

	bArrayMemberAlreadyWritten := false
	for keys.HasNext() {
		queryResponse, err := keys.Next()
		if err != nil {
			return nil, err
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"Key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"Record\":")
		// Record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return buffer.Bytes(), nil
}

func (s *SmartContract) queryAllApis(APIstub shim.ChaincodeStubInterface) sc.Response {
	funcName, _, _, _ := runtime.Caller(0)
	fmt.Printf("------ %s ------\n", funcName)

	startKey := apiKeyPrefix + "0"
	endKey := apiKeyPrefix + "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	ret, err := s.queryArrayByIteratorKeys(resultsIterator)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(ret)
}

func (s *SmartContract) queryAllOpCodes(APIstub shim.ChaincodeStubInterface) sc.Response {
	funcName, _, _, _ := runtime.Caller(0)
	fmt.Printf("------ %s ------\n", funcName)

	startKey := opCodeKeyPrefix + "0"
	endKey := opCodeKeyPrefix + "999"

	resultsIterator, err := APIstub.GetStateByRange(startKey, endKey)
	if err != nil {
		return shim.Error(err.Error())
	}
	defer resultsIterator.Close()

	ret, err := s.queryArrayByIteratorKeys(resultsIterator)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(ret)
}

// func (s *SmartContract) changeCarOwner(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {

//  if len(args) != 2 {
// 	 return shim.Error("Incorrect number of arguments. Expecting 2")
//  }

//  carAsBytes, _ := APIstub.GetState(args[0])
//  car := Car{}

//  json.Unmarshal(carAsBytes, &car)
//  car.Owner = args[1]

//  carAsBytes, _ = json.Marshal(car)
//  APIstub.PutState(args[0], carAsBytes)

//  return shim.Success(nil)
// }

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
