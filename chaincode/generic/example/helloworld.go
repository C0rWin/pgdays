/*
Copyright IBM Corp. 2017 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"fmt"

	"github.com/C0rwin/pgdays/chaincode/generic"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// Example set of action to provide for chaincode creation
var (
	myactions = map[string]generic.ChaincodeAction{
		"Init": func(params []string, stub shim.ChaincodeStubInterface) peer.Response {
			fmt.Println("Initialize chaincode")
			return shim.Success(nil)
		},
		"AddKey": func(params []string, stub shim.ChaincodeStubInterface) peer.Response {
			fmt.Println("Adding key")
			return shim.Success(nil)
		},
	}
)

func main() {
	err := shim.Start(generic.NewGenericChaincode(myactions))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
