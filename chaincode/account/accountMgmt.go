package main

import (
	"encoding/json"
	"fmt"

	"github.com/hyperledger/fabric/common/util"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type Account struct {
	PersonID string  `json:"personID"`
	Number   string  `json:"number"`
	Balance  float64 `json:"balance"`
}

type AccountAction func(params []string, stub shim.ChaincodeStubInterface) peer.Response

type accountManagement struct {
	actions          map[string]AccountAction
	personManagement string
	channelName      string
}

// Init
func (am *accountManagement) Init(stub shim.ChaincodeStubInterface) peer.Response {
	am.actions = map[string]AccountAction{
		"openAccount": am.OpenAccount,
	}
	_, params := stub.GetFunctionAndParameters()

	// get person management chaincode name during chaincode instantiation
	am.personManagement = params[0]
	am.channelName = params[1]

	fmt.Println("Chaincode initialized, person management chaincode name", am.personManagement,
		"channel name", am.channelName)
	return shim.Success(nil)
}

// Invoke
func (am *accountManagement) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	actionName, params := stub.GetFunctionAndParameters()

	if action, ok := am.actions[actionName]; ok {
		return action(params, stub)
	}

	return shim.Error(fmt.Sprintf("[ERROR] No <%s> action defined", actionName))
}

// OpenAccount
func (am *accountManagement) OpenAccount(params []string, stub shim.ChaincodeStubInterface) peer.Response {
	jsonObj := params[0]
	var account Account

	// Read person info into struct
	err := json.Unmarshal([]byte(jsonObj), &account)
	if err != nil {
		errMsg := fmt.Sprintf("[ERROR] cannot unmarshal account information, because of %s", err)
		fmt.Println(errMsg)
		return shim.Error(errMsg)
	}

	chainCodeArgs := util.ToChaincodeArgs("getPerson", account.PersonID)
	response := stub.InvokeChaincode(am.personManagement, chainCodeArgs, am.channelName)

	if response.Status != shim.OK {
		errMsg := fmt.Sprintf("Cannot verify existance of person with id %s, reason %s", account.PersonID, response.Message)
		fmt.Println(errMsg)
		return shim.Error(errMsg)
	}

	if err := stub.PutState(account.PersonID, []byte(jsonObj)); err != nil {
		errMsg := fmt.Sprintf("Cannot store account information, account id %s, failed to open account due to %s",
			account.PersonID, err)
		fmt.Println(errMsg)
		return shim.Error(errMsg)
	}

	fmt.Println("Opened new account", account)
	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(accountManagement))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
