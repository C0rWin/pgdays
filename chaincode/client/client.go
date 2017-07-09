package main

import (
	"fmt"
	"path/filepath"

	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	"github.com/hyperledger/fabric-sdk-go/def/fabapi"
	"github.com/hyperledger/fabric-sdk-go/pkg/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/orderer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/peer"
	bccspFactory "github.com/hyperledger/fabric/bccsp/factory"
)

func main() {
	var err error
	conf, err := config.InitConfig("./client/config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	cl := fabricclient.NewClient(conf)
	bccspFactory.InitFactories(nil)
	cl.SetCryptoSuite(bccspFactory.GetDefault())

	privKey := filepath.Join(conf.CryptoConfigPath(), "peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/045a627439618dbedc9a68888fee2fb89a0dc2b0b31b965cf57ddbc9e43dda0b_sk")
	pubKey := filepath.Join(conf.CryptoConfigPath(), "peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem")

	mspID, err := conf.MspID("org1")
	if err != nil {
		fmt.Println(err)
		return
	}

	user, err := fabapi.NewPreEnrolledUser(conf, privKey, pubKey, "user1", mspID, cl.GetCryptoSuite())
	if err != nil {
		fmt.Println(err)
		return
	}
	cl.SetUserContext(user)

	ordererConf, err := conf.OrdererConfig("orderer0")
	if err != nil {
		fmt.Println(err)
		return
	}

	o, err := orderer.NewOrderer(fmt.Sprintf("%s:%d", ordererConf.Host, ordererConf.Port),
		filepath.Join(conf.CryptoConfigPath(), "ordererOrganizations/example.com/orderers/orderer.example.com/msp/cacerts/example.com-cert.pem"),
		"orderer.example.com", conf)
	if err != nil {
		fmt.Println(err)
		return
	}

	peers, err := conf.PeersConfig("org1")
	if err != nil {
		fmt.Println(err)
		return
	}

	p, err := peer.NewPeer(fmt.Sprintf("%s:%d", peers[0].Host, peers[0].Port), conf)
	if err != nil {
		fmt.Println(err)
	}

	ch, err := cl.NewChannel("mychannel")
	if err != nil {
		fmt.Println(err)
		return
	}

	ch.AddOrderer(o)
	ch.AddPeer(p)
	ch.SetPrimaryPeer(p)
	cl.SaveUserToStateStore(user, true)

	txRequest := apitxn.ChaincodeInvokeRequest{
		Targets:      []apitxn.ProposalProcessor{p},
		Fcn:          "AddKey",
		Args:         []string{"AddKey"},
		TransientMap: map[string][]byte{},
		ChaincodeID:  "helloworld",
	}

	proposalResponse, _, err := ch.SendTransactionProposal(txRequest)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%v\n", proposalResponse[0].ProposalResponse)

	tx, err := ch.CreateTransaction(proposalResponse)
	if err != nil {
		fmt.Println(err)
		return
	}

	txResponse, err := ch.SendTransaction(tx)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(txResponse[0])
}
