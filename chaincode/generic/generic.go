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
package generic

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

// ChaincodeAction defines simple action which should be used to define
// chaincode behaviour for certain actions
type ChaincodeAction func(params []string, stub shim.ChaincodeStubInterface) peer.Response

// genericChaincode a template structure which has to be initialized with
// map of chaincode actions, and will be use this map to dispatch actions
// based on chaincode function name parameter
type genericChaincode struct {
	actions map[string]ChaincodeAction
}

// Init initialize chaincode with map of actions and initialize the chaincode
func (cc *genericChaincode) Init(stub shim.ChaincodeStubInterface) peer.Response {
	_, params := stub.GetFunctionAndParameters()
	// Initializing chaincode
	fmt.Println("Following actions are available")
	for action := range cc.actions {
		fmt.Printf("\t\t%s\n", action)
	}
	if action, ok := cc.actions["Init"]; ok {
		return action(params, stub)
	}

	return shim.Success(nil)
}

func (cc *genericChaincode) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	actionName, params := stub.GetFunctionAndParameters()

	if action, ok := cc.actions[actionName]; ok {
		return action(params, stub)
	}

	return shim.Error(fmt.Sprintf("No <%s> action defined", actionName))
}

func NewGenericChaincode(actions map[string]ChaincodeAction) shim.Chaincode {
	return &genericChaincode{actions}
}
