package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime"
	"time"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	sc "github.com/hyperledger/fabric/protos/peer"
)

type apiName string

const (
	apiInitLedger  apiName = "init"
	apiSupportAPIs apiName = "apis"
	apiCreate      apiName = "create"
	apiModify      apiName = "modify"
	apiArchive     apiName = "archive"
	apiDelete      apiName = "delete"
	apiQuery       apiName = "query"
	apiHistory     apiName = "history"
)

type opCode string

const (
	opCreate  opCode = "create"
	opModify  opCode = "modify"
	opArchive opCode = "archive"
	opDelete  opCode = "delete"
)

const (
	edocKeyPrefix = "EDoc_"
)

func edocKey(systemName, documentID string) string {
	return edocKeyPrefix + systemName + documentID
}

// SmartContract defines the smart contract for CRP(China Resource Power)
type SmartContract struct {
}

// EDocument defines the e-document in CRP
type EDocument struct {
	SystemName string      `json:"system-name"`
	DocumentID string      `json:"document-id"`
	UID        string      `json:"uid"`
	OpCode     opCode      `json:"op-code"`
	OpTime     string      `json:"op-time"`
	Hash       string      `json:"hash,omitempty"`
	MetaData   interface{} `json:"meta-data,omitempty"`
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
	} else if api == apiCreate {
		return s.createRecord(APIstub, args)
	} else if api == apiModify {
		return s.commitRecord(APIstub, args, opModify)
	} else if api == apiArchive {
		return s.commitRecord(APIstub, args, opArchive)
	} else if api == apiDelete {
		return s.commitRecord(APIstub, args, opDelete)
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
	return shim.Success(nil)
}

func (s *SmartContract) createRecord(APIstub shim.ChaincodeStubInterface, args []string) sc.Response {
	printFunctionName()

	if len(args) < 5 {
		return shim.Error("Usage: (SystemName,DocumentID,UID,OpTime,Hash,MetaData(optional))")
	}

	key := edocKey(args[0], args[1])
	if docBytes, err := APIstub.GetState(key); len(docBytes) > 0 && err == nil {
		return shim.Error("document(" + args[1] + ") in system(" + args[0] + ") already exists")
	}

	edoc := EDocument{
		SystemName: args[0],
		DocumentID: args[1],
		UID:        args[2],
		OpCode:     opCreate,
		OpTime:     args[3],
		Hash:       args[4],
	}
	// if len(args) >= 5 {
	// 	edoc.Hash = args[4]
	// }
	if len(args) >= 6 {
		edoc.MetaData = args[5]
	}

	docBytes, err := json.Marshal(edoc)
	if err != nil {
		return shim.Error(err.Error())
	}
	err = APIstub.PutState(edocKey(edoc.SystemName, edoc.DocumentID), docBytes)
	if err != nil {
		return shim.Error(err.Error())
	}

	return shim.Success([]byte(APIstub.GetTxID()))
}

func (s *SmartContract) commitRecord(APIstub shim.ChaincodeStubInterface, args []string, code opCode) sc.Response {
	printFunctionName()

	if len(args) < 5 {
		return shim.Error("Usage: (SystemName,DocumentID,UID,OpTime,Hash,MetaData(optional))")
	}

	key := edocKey(args[0], args[1])

	docBytes, err := APIstub.GetState(key)
	if err != nil { // not exists
		return shim.Error(err.Error())
	}
	var edoc EDocument
	if err := json.Unmarshal(docBytes, &edoc); err != nil {
		return shim.Error(err.Error())
	}

	// check status
	statusInfo := map[opCode]int{
		opCreate:  0,
		opModify:  1,
		opArchive: 2,
		opDelete:  3,
	}
	if statusInfo[edoc.OpCode] > statusInfo[code] {
		return shim.Error("document status could not reverse back")
	}

	// update
	edoc.UID = args[2]
	edoc.OpCode = code
	edoc.OpTime = args[3]
	edoc.Hash = args[4]

	if len(args) >= 6 {
		edoc.MetaData = args[5]
	}

	docBytes, err = json.Marshal(edoc)
	if err != nil {
		return shim.Error(err.Error())
	}
	if err := APIstub.PutState(key, docBytes); err != nil {
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

	if len(args) != 2 {
		return shim.Error("Incorrect number of arguments. Expecting 2 (systemName, documentID)")
	}

	historyIterator, err := APIstub.GetHistoryForKey(edocKey(args[0], args[1]))
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
		ts := queryResponse.GetTimestamp()
		t := time.Unix(ts.Seconds, int64(ts.Nanos))
		buffer.WriteString(t.Format("2006-01-02 15:04:05"))

		buffer.WriteString(", \"record\":")
		buffer.WriteString(string(queryResponse.Value))

		buffer.WriteString("}")
		bArrayMemberAlreadyWritten = true
	}

	buffer.WriteString("]")

	return shim.Success(buffer.Bytes())
}

func (s *SmartContract) queryAllApis(APIstub shim.ChaincodeStubInterface) sc.Response {
	printFunctionName()

	apiNames := []apiName{
		apiSupportAPIs,
		apiInitLedger,
		apiCreate,
		apiModify,
		apiArchive,
		apiDelete,
		apiQuery,
		apiHistory,
	}

	ret, err := json.Marshal(apiNames)
	if err != nil {
		return shim.Error(err.Error())
	}
	return shim.Success(ret)
}

// The main function is only relevant in unit test mode. Only included here for completeness.
func main() {
	// Create a new Smart Contract
	err := shim.Start(new(SmartContract))
	if err != nil {
		fmt.Printf("Error creating new Smart Contract: %s", err)
	}
}
