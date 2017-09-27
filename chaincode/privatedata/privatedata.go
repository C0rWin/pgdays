package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type smartContractAction func(args []string, stub shim.ChaincodeStubInterface) peer.Response

type smartContract struct {
}

var (
	actions = map[string]smartContractAction{
		"write": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
			ns := args[0]
			key := args[1]
			val := args[2]
			if err := stub.PutPrivateData(ns, key, []byte(val)); err != nil {
				fmt.Printf(fmt.Sprintf("Error while storing private data <%s, %s, %s> the error is %s", ns, key, val, err))
				return shim.Error(fmt.Sprintf("%s", err))
			}
			fmt.Printf("Storing private data <collection, key, value>  = <%s, %s, %s>\n", ns, key, val)
			return shim.Success(nil)
		},
		"read": func(args []string, stub shim.ChaincodeStubInterface) peer.Response {
			ns := args[0]
			key := args[1]
			fmt.Printf("Reading chaincode private data <%s, %s> \n", ns, key)
			if data, err := stub.GetPrivateData(ns, key); err != nil {
				fmt.Printf(fmt.Sprintf("Error reading key %s from collection %s, due to %s", ns, key, err))
				return shim.Error(fmt.Sprintf("%s", err))
			} else {
				return shim.Success(data)
			}
		},
	}
)

func (mock *smartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Printf("Initializing chaincode \n")
	return shim.Success(nil)
}

func (mock *smartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Printf("Invoking chaincode \n")
	function, args := stub.GetFunctionAndParameters()
	if action, ok := actions[function]; !ok {
		return shim.Error(fmt.Sprintf("Wrong function name, %s", function))
	} else {
		return action(args, stub)
	}
}

func main() {
	err := shim.Start(new(smartContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
