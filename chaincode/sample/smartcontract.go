package sample

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
	"github.com/hyperledger/fabric/protos/peer"
)

type smartContract struct {
}

func (smartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {

}

func (smartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {

	funcName, params := stub.GetFunctionAndParameters()

	indexName := "txID~key"

	if funcName == "addNewKey" {

		key := params[0]
		value := params[1]

		keyTxIdKey, err := stub.CreateCompositeKey(indexName, []string{stub.GetTxID(), key})
		if err != nil {
			return shim.Error(err.Error())
		}

		creator, _ := stub.GetCreator()

		// Add key and value to the state
		stub.PutState(key, []byte(value))
		stub.PutState(keyTxIdKey, creator)

	} else if funcName == "checkTxID" {
		txID := params[0]

		it, _ := stub.GetStateByPartialCompositeKey(indexName, []string{txID})

		for it.HasNext() {
			keyTxIdRange, err := it.Next()
			if err != nil {
				return shim.Error(err.Error())
			}

			_, keyParts, _ := stub.SplitCompositeKey(keyTxIdRange.Key)
			key := keyParts[1]
			fmt.Printf("key affected by txID %s is %s\n", txID, key)
			txIDCreator := keyTxIdRange.Value

			sId := &msp.SerializedIdentity{}
			err := proto.Unmarshal(txIDCreator, sId)
			if err != nil {
				return shim.Error(fmt.Sprintf("Could not deserialize a SerializedIdentity, err %s", err))
			}

			bl, _ := pem.Decode(sId.IdBytes)
			if bl == nil {
				return shim.Error(fmt.Sprintf("Could not decode the PEM structure"))
			}
			cert, err := x509.ParseCertificate(bl.Bytes)
			if err != nil {
				return shim.Error(fmt.Sprintf("ParseCertificate failed %s", err))
			}

			fmt.Printf("Certificate of txID %s creator is %s", txID, cert)
		}
	}

	return shim.Success(nil)
}
