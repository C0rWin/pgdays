package temperature

import (
	"fmt"

	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/peer"
)

type temperatureSmartContract struct {
}

func (contract *temperatureSmartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("Initialize chaincode if needed")
	return shim.Success(nil)
}

func (contract *temperatureSmartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	funcName, params := stub.GetFunctionAndParameters()

	if funcName == "addTemperature" {
		// Store observation into ledger
		stub.PutState("temperature", []byte(params[0]))
	} else if funcName == "getTemperatures" {
		iter, err := stub.GetHistoryForKey("temperature")
		if err != nil {
			shim.Error(fmt.Sprintf("%s", err))
		}

		var result []string
		for iter.HasNext() {
			mod, err := iter.Next()
			if err != nil {
				shim.Error(fmt.Sprintf("%s", err))
			}
			result = append(result, string(mod.Value))
		}
		return shim.Success([]byte(result))
	}
	return shim.Success(nil)
}
