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
	"encoding/json"
	"fmt"

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

// String returns string representation of the account
func (a *Person) String() string {
	return fmt.Sprintf("Person id [%d], name = %s, phone numner = %s, account humber %s",
		a.ID, a.Name, a.Phone, a.Address)
}

// PersonAction
type PersonAction func(params []string, stub shim.ChaincodeStubInterface) peer.Response

// personManagement the chaincode interface implementation to manage
// the ledger of person records
type personManagement struct {
	actions map[string]PersonAction
}

// Init initialize chaincode with mapping between actions and real methods
func (pm *personManagement) Init(stub shim.ChaincodeStubInterface) peer.Response {
	pm.actions = map[string]PersonAction{
		"addPerson":        pm.AddPerson,
		"getPerson":        pm.GetPerson,
		"delPerson":        pm.DelPerson,
		"updateAddress":    pm.UpdateAddress,
		"getPreviousValue": pm.GetPreviousValue,
	}

	fmt.Println("Chaincode has been initialized")
	fmt.Println("Following actions are available")
	for action := range pm.actions {
		fmt.Printf("\t\t%s\n", action)
	}
	return shim.Success(nil)
}

// Invoke handles chaincode invocation logic, executes actual code
// for given action name
func (pm *personManagement) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	// Actions supported by this chaincode will be
	// * AddPerson - adds new account into
	// * GetPerson
	// * DelPerson
	// * UpdateAddress
	// * GetPreviousValue
	actionName, params := stub.GetFunctionAndParameters()

	if action, ok := pm.actions[actionName]; ok {
		return action(params, stub)
	}

	return shim.Error(fmt.Sprintf("No <%s> action defined", actionName))
}

// AddPerson inserts new person into ledger
func (pm *personManagement) AddPerson(params []string, stub shim.ChaincodeStubInterface) peer.Response {
	jsonObj := params[0]
	var person Person

	// Read person info into struct
	json.Unmarshal([]byte(jsonObj), &person)

	val, err := stub.GetState(person.ID)
	if err != nil {
		fmt.Printf("[ERROR] cannot get state, because of %s\n", err)
		return shim.Error(fmt.Sprintf("%s", err))
	}

	if val != nil {
		errMsg := fmt.Sprintf("[ERROR] person already exists, cannot create two accounts with same ID <%d>", person.ID)
		fmt.Println(errMsg)
		return shim.Error(errMsg)
	}

	fmt.Println("Adding new person", person)
	if err = stub.PutState(person.ID, []byte(jsonObj)); err != nil {
		errMsg := fmt.Sprintf("[ERROR] cannot store person record with id <%d>, due to %s", person.ID, err)
		fmt.Println(errMsg)
		return shim.Error(errMsg)
	}
	return shim.Success(nil)
}

// GetPerson retrieves persons record from the ledger
func (pm *personManagement) GetPerson(params []string, stub shim.ChaincodeStubInterface) peer.Response {
	// Extract person id from the parameters
	personID := params[0]

	val, err := stub.GetState(personID)
	if err != nil {
		fmt.Printf("[ERROR] cannot get state, because of %s\n", err)
		return shim.Error(fmt.Sprintf("%s", err))
	}

	if val == nil {
		errMsg := fmt.Sprintf("[ERROR] Cannot find person with id %s", personID)
		fmt.Println(errMsg)
		return shim.Error(errMsg)
	}

	// Returning person struct back to the application
	var person Person
	// Read person info into struct
	json.Unmarshal(val, &person)
	fmt.Println("Returning information about", string(val))

	return shim.Success(val)
}

// DelPerson removes person record from the ledger
func (pm *personManagement) DelPerson(params []string, stub shim.ChaincodeStubInterface) peer.Response {
	// Extract person id
	personID := params[0]

	// Trying to delete person with given id
	err := stub.DelState(personID)
	if err != nil {
		fmt.Printf("[ERROR] cannot delete person record with id %s, because of %s\n", personID, err)
		return shim.Error(fmt.Sprintf("%s", err))
	}

	// Person record deleted
	fmt.Printf("Person with id %s, has been deleted successefully \n", personID)
	return shim.Success(nil)
}

// // UpdateAddress updates persons address
func (pm *personManagement) UpdateAddress(params []string, stub shim.ChaincodeStubInterface) peer.Response {
	// Extract person id from the parameters
	personID := params[0]
	newAddress := params[1]

	val, err := stub.GetState(personID)
	if err != nil {
		fmt.Printf("[ERROR] cannot get state, because of %s\n", err)
		return shim.Error(fmt.Sprintf("%s", err))
	}

	if val == nil {
		errMsg := fmt.Sprintf("[ERROR] Cannot find person with id %d", personID)
		fmt.Println(errMsg)
		return shim.Error(errMsg)
	}

	// Returning person struct back to the application
	var person Person
	// Read person info into struct
	json.Unmarshal(val, &person)
	person.Address = newAddress
	// Marshal back to json
	val, err = json.Marshal(person)
	if err != nil {
		errMsg := fmt.Sprintf("[ERROR] person id %s, cannot marshal person struct due to %s", personID, err)
		return shim.Error(errMsg)
	}

	fmt.Println("Updating person address", person)
	// Update key with new value
	if err = stub.PutState(person.ID, val); err != nil {
		errMsg := fmt.Sprintf("[ERROR] cannot store person record with id <%d>, due to %s", person.ID, err)
		fmt.Println(errMsg)
		return shim.Error(errMsg)
	}

	fmt.Println("Address updated successeful for", person)
	return shim.Success(nil)
}

// GetPreviousValue reads previous value of given key
func (pm *personManagement) GetPreviousValue(params []string, stub shim.ChaincodeStubInterface) peer.Response {
	// Extract personID
	personID := params[0]

	historyIer, err := stub.GetHistoryForKey(personID)

	if err != nil {
		errMsg := fmt.Sprintf("[ERROR] cannot retreive history of person record with id <%s>, due to %s", personID, err)
		fmt.Println(errMsg)
		return shim.Error(errMsg)
	}

	if historyIer.HasNext() {
		modification, err := historyIer.Next()
		if err != nil {
			errMsg := fmt.Sprintf("[ERROR] cannot read person record modification, id <%s>, due to %s", personID, err)
			fmt.Println(errMsg)
			return shim.Error(errMsg)
		}
		fmt.Println("Returning information about", string(modification.Value))
		return shim.Success(modification.Value)
	}

	fmt.Printf("No history found for person with id %s\n", personID)
	return shim.Success([]byte(fmt.Sprintf("No history for person with id %s", personID)))
}

func main() {
	err := shim.Start(new(personManagement))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
