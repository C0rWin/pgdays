package main

import (
	"fmt"

	"github.com/hyperledger/fabric-sdk-go/api/apitxn"
	"github.com/hyperledger/fabric-sdk-go/def/fabapi"
	"github.com/hyperledger/fabric-sdk-go/pkg/config"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/channel"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/orderer"
	"github.com/hyperledger/fabric-sdk-go/pkg/fabric-client/peer"
	bccspFactory "github.com/hyperledger/fabric/bccsp/factory"
)

func main() {
	configImpl, err := config.InitConfig("client/client_config.yaml")
	if err != nil {
		fmt.Printf("Cannot initialize config, because of %s\n", err)
		return
	}
	client := fabricclient.NewClient(configImpl)

	peersConf, err := configImpl.PeersConfig("org1")
	if err != nil {
		fmt.Printf("Unable to read peers configuration for org1, because of %s\n", err)
		return
	}

	if len(peersConf) == 0 {
		fmt.Println("No peers defines in org1 configuration")
		return
	}

	peerConf := peersConf[0]
	peerEndpoint := fmt.Sprintf("%s:%d", peerConf.Host, peerConf.Port)
	node, err := peer.NewPeerTLSFromCert(peerEndpoint, peerConf.TLS.Certificate, peerConf.TLS.ServerHostOverride, configImpl)
	if err != nil {
		fmt.Println("Cannot create peer", err)
	}

	ordererConf, err := configImpl.OrdererConfig("orderer0")
	if err != nil {
		fmt.Printf("Cannot read orderer config %s\n", err)
		return
	}
	ordererEndpoint := fmt.Sprintf("%s:%d", ordererConf.Host, ordererConf.Port)
	orderer, err := orderer.NewOrderer(ordererEndpoint, ordererConf.TLS.Certificate, ordererConf.TLS.ServerHostOverride, configImpl)
	if err != nil {
		fmt.Printf("wasn't able to create orderer, due to %s\n", err)
		return
	}

	channel, err := channel.NewChannel("mychannel", client)
	if err != nil {
		fmt.Printf("wasn't able to create channel, due to %s\n", err)
		return
	}

	bccspFactory.InitFactories(configImpl.CSPConfig())

	user, err := fabapi.NewPreEnrolledUser(configImpl, "/Users/bartem/golang/src/github.com/C0rwin/pgdays/artifacts/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/keystore/926f42f5b4537e68ed6aab06c1b3395fa0cda9a0db4f760f9eb2e2a89c23bd09_sk",
		"/Users/bartem/golang/src/github.com/C0rwin/pgdays/artifacts/crypto-config/peerOrganizations/org1.example.com/users/Admin@org1.example.com/msp/signcerts/Admin@org1.example.com-cert.pem",
		"peerorg1Admin", "Org1MSP", bccspFactory.GetDefault())

	if err != nil {
		fmt.Printf("cannot create org1 user, due to %s\n", err)
		return
	}
	client.SetUserContext(user)

	channel.AddOrderer(orderer)
	channel.AddPeer(node)

	proposal, err := channel.CreateTransactionProposal("helloworld", "mychannel", []string{}, true, nil)
	if err != nil {
		fmt.Printf("wasn't able to create ctransaction proposal, due to %s\n", err)
		return
	}

	response, err := channel.SendTransactionProposal(proposal, 0, []apitxn.ProposalProcessor{node})
	if err != nil {
		fmt.Printf("error while receiving response for transaction proposal, due to %s\n", err)
		return
	}

	fmt.Print(response)
}
