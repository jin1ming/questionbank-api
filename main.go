package main

import (
	"net/http"
	"os"
	apiservice "questionbank-api/controller/api/v1"
	router "questionbank-api/router/v1"
	"questionbank-api/sdkInit"
)


func main()  {

	// 初始化redis
	apiservice.InitRedis("localhost:6379", "", 0)

	// 初始化sdk、创建通道、安装链码
	initInfo := &sdkInit.InitInfo{

		ChannelID: "miracle",
		ChannelConfig: os.Getenv("GOPATH") + "/src/questionbank-api/fixtures/artifacts/channel.tx",
		OrgAdmin:"Admin",
		OrgName:"Org1",
		OrdererOrgName: "orderer.questionbank.com",

		ChaincodeID: "questionbank",
		ChaincodeGoPath: os.Getenv("GOPATH"),
		ChaincodePath: "questionbank-api/chaincode/", //这里不需要src！ 开始也不能有/
		UserName:"User1",
	}
	apiservice.InitSdk(initInfo,"config.yaml")

	// 初始化gin
	r := router.InitRouter()
	s := &http.Server{
		Addr:           ":8080",
		Handler:        r,
	}
	if err := s.ListenAndServe(); err != nil {
		panic(err.Error())
	}
}


