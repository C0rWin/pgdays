package main

import (
	"fmt"

	"crypto/x509"
	"encoding/pem"

	"github.com/golang/protobuf/proto"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"github.com/hyperledger/fabric/protos/msp"
	"github.com/hyperledger/fabric/protos/peer"
)

type smartContract struct {
}

func (mock *smartContract) Init(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("Init")
	return shim.Success(nil)
}

func (*smartContract) Invoke(stub shim.ChaincodeStubInterface) peer.Response {
	fmt.Println("Invoke")

	serializedID, _ := stub.GetCreator()

	sId := &msp.SerializedIdentity{}
	err := proto.Unmarshal(serializedID, sId)
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

	fmt.Println(cert)

	return shim.Success(nil)
}

func main() {
	err := shim.Start(new(smartContract))
	if err != nil {
		fmt.Printf("Error starting Simple chaincode: %s", err)
	}
}
