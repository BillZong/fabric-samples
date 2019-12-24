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
	apiCreate         apiName = "create"
	apiCommit         apiName = "commit"
	apiQuery          apiName = "query"
	apiHistory        apiName = "history"
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
	edocKeyPrefix   = "EDoc_"
)

// Define the Smart Contract structure
type SmartContract struct {
}

// EDocument is the e-document structure for CRP(China Resource Power).
type EDocument struct {
	DocumentID  string `json:"document-id"`
	UID         string `json:"uid"`
	OpCode      opCode `json:"op-code"`
	OpTime      string `json:"op-time"`
	Hash        string `json:"hash"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type KVResult struct {
	Key   string
	Value []byte
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
	} else if api == apiCreate {
		return s.createRecord(APIstub, args)
	} else if api == apiCommit {
		return s.commitRecord(APIstub, args)
	} else if api == apiQuery {
		return s.queryRecord(APIstub, args)
	} else if api == apiHistory {
		return s.queryRecordHistory(APIstub, args)
	}

	return shim.Error(fmt.Sprintf("Not support function name: %s, args: %s", function, args))
}

func printFunctionName() {
	funcName, _, _, _ := runtime.Caller(1)
	fmt.Printf("------ %v ------\n", funcName)
}

func (s *SmartContract) initLedger(APIstub shim.ChaincodeStubInterface) sc.Response {
	printFunctionName()

	apiNames := []apiName{
		apiSupportAPIs,
		apiSupportOpCodes,
		apiInitLedger,
		apiCreate,
		apiCommit,
		apiQuery,
		apiHistory,
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

func (s *SmartContract) createRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	printFunctionName()

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 7")
	}

	edoc := EDocument{
		DocumentID:  args[0],
		UID:         args[1],
		OpCode:      opCode(args[2]),
		OpTime:      args[3],
		Hash:        args[4],
		Name:        args[5],
		Description: args[6],
	}

	docBytes, err := json.Marshal(edoc)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutState(edocKeyPrefix+args[0], docBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(APIstub.GetTxID()))
}

func (s *SmartContract) commitRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	printFunctionName()

	if len(args) != 7 {
		return shim.Error("Incorrect number of arguments. Expecting 5")
	}

	key := edocKeyPrefix + args[0]

	docBytes, err := APIstub.GetState(key)
	if err != nil { // not exists
		return shim.Error(err.Error())
	}
	var edoc EDocument
	if err := json.Unmarshal(docBytes, &edoc); err != nil {
		return shim.Error(err.Error())
	}
	edoc.UID = args[1]
	edoc.OpCode = opCode(args[2])
	edoc.OpTime = args[3]
	edoc.Hash = args[4]
	edoc.Name = args[5]
	edoc.Description = args[6]

	docBytes, err = json.Marshal(edoc)
	if err != nil {
		return shim.Error(err.Error())
	}
	if err := APIstub.PutState(edocKeyPrefix+args[0], docBytes); err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(APIstub.GetTxID()))
}

func (s *SmartContract) queryRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	printFunctionName()

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	docBytes, err := APIstub.GetState(edocKeyPrefix + args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(docBytes)
}

func (s *SmartContract) queryRecordHistory(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	printFunctionName()

	if len(args) != 1 {
		return shim.Error("Incorrect number of arguments. Expecting 1")
	}

	historyIterator, err := APIstub.GetHistoryForKey(edocKeyPrefix + args[0])
	if err != nil {
		return shim.Error(err.Error())
	}
	defer historyIterator.Close()

	var buffer bytes.Buffer
	buffer.WriteString("[")
	bArrayMemberAlreadyWritten := false
	for historyIterator.HasNext() {
		queryResponse, err := historyIterator.Next()
		if err != nil {
			return shim.Error(err.Error())
		}
		// Add a comma before array members, suppress it for the first array member
		if bArrayMemberAlreadyWritten == true {
			buffer.WriteString(",")
		}
		buffer.WriteString("{\"txid\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.GetTxId())
		buffer.WriteString("\"")

		buffer.WriteString(", \"time\":")
		buffer.WriteString(queryResponse.GetTimestamp().String())

		buffer.WriteString(", \"record\":")
		buffer.WriteString(string(queryResponse.Value))

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

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
		buffer.WriteString("{\"key\":")
		buffer.WriteString("\"")
		buffer.WriteString(queryResponse.Key)
		buffer.WriteString("\"")

		buffer.WriteString(", \"record\":")
		// record is a JSON object, so we write as-is
		buffer.WriteString(string(queryResponse.Value))
		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}
	buffer.WriteString("]")

	return buffer.Bytes(), nil
}

func (s *SmartContract) queryAllApis(APIstub shim.ChaincodeStubInterface) sc.Response {
	printFunctionName()

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
	printFunctionName()

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
// printFunctionName()

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
