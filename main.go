package main

import (
	"fmt"
	"os"
	"questionbank-api/sdkInit"
)



const (
	configFile = "config.yaml"
	initialized = false
)

func main() {

	initInfo := &sdkInit.InitInfo{

		ChannelID: "miracle",
		ChannelConfig: os.Getenv("GOPATH") + "/src/questionbank-api/fixtures/artifacts/channel.tx",
		OrgAdmin:"Admin",
		OrgName:"Org1",
		OrdererOrgName: "orderer.questionbank.com",

		ChaincodeID: "questionbank",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath: "questionbank-api/chaincode/", //这里不需要src！
		UserName:"User1",
	}

	sdk, err := sdkInit.SetupSDK(configFile, initialized)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	defer sdk.Close()

	err = sdkInit.CreateChannel(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	channelClient, err := sdkInit.InstallAndInstantiateCC(sdk, initInfo)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println(channelClient)

}

