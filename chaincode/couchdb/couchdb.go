package main

import (
	"github.com/hyperledger/fabric/protos/peer"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"fmt"
	"encoding/json"
	"os"
)

type couchdbChaincode struct {
}

type Record struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

func (couchdbChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("Chaincode initialized")
	return shim.Success(nil)
}

func (couchdbChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	funcName, args := stub.GetFunctionAndParameters() // we not expect to get any parameters
	value := args[0]

	// we are testing couchdb capabilities to query and store values in couchdb
	if funcName == "store" {

		// create a test record
		record := Record{
			ID:    stub.GetTxID(),
			Value: value,
		}

		dbRecord, _ := json.Marshal(record)
		if err := stub.PutState(stub.GetTxID(), dbRecord); err != nil {
			return shim.Error(err.Error())
		}

	} else if funcName == "query" {
		// create selector to query database
		query, _ := json.Marshal(struct {
			Selector struct {
				Data struct {
					Value string `json:"value"`
				} `json:"data"`
			} `json:"selector"`
		}{
			Selector: struct {
				Data struct {
					Value string `json:"value"`
				} `json:"data"`
			}{
				Data: struct {
					Value string `json:"value"`
				}{
					Value: value,
				}},
		})

		fmt.Println("XXXXX >>>> (QUERY) ---->  ", string(query))
		it, err := stub.GetQueryResult(string(query))
		if err != nil {
			return shim.Error(err.Error())
		}


		for it.HasNext() {
			val, _ := it.Next()
			return shim.Success(val.Value)
		}

		fmt.Println("XXX: First trial wasn't able to find a record, attempt I.")

		// create selector to query database
		query, _ = json.Marshal(struct {
			Selector struct {
					Value string `json:"value"`
			} `json:"selector"`
		}{
			Selector: struct {
					Value string `json:"value"`
			}{
					Value: value,
				},
		})

		fmt.Println("XXXXX >>>> (QUERY) ---->  ", string(query))
		it, err = stub.GetQueryResult(string(query))
		if err != nil {
			return shim.Error(err.Error())
		}


		for it.HasNext() {
			val, _ := it.Next()
			return shim.Success(val.Value)
		}

		fmt.Println("XXX: First trial wasn't able to find a record, attempt II.")

		return shim.Error("no record found")
	}

	return shim.Success(nil)

}

func main() {
	if err := shim.Start(&couchdbChaincode{}); err != nil {
		fmt.Println("Failed to start chaincode,", err)
		os.Exit(-1)
	}
}
