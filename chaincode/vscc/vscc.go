package main

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type myNewCoolVSCC struct {
}

func (m *myNewCoolVSCC) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("Started")
	return shim.Success(nil)
}

func (m *myNewCoolVSCC) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("Checking VSCC :-)")
	return shim.Success(nil)
}

func main() {
	err := shim.Start(&myNewCoolVSCC{})
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
