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
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// Person
type Person struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Address string `json:"address"`
}

// personManagement the chaincode interface implementation to manage
// the ledger of person records
type personManagement struct {
}

func (p *personManagement) Init(stub shim.ChaincodeStubInterface) peer.Response {
	return shim.Success(nil)
}

func (p *personManagement) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	actionName, params := stub.GetFunctionAndParameters()

	if actionName == "addPerson" {

	}

}

func main() {
	/*
		err := shim.Start(new(personManagement))
		if err != nil {
			fmt.Printf("Error starting Simple chaincode: %s", err)
		}
	*/
}
