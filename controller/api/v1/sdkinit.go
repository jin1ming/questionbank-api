package v1

// 用于sdk初始化、创建通道、安装链码
import (
	"fmt"
	"os"
	"questionbank-api/sdkInit"
)

func InitSdk(initInfo *sdkInit.InitInfo, configFile string) {

	var err error
	sdk, err = sdkInit.SetupSDK(configFile)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	defer sdk.Close()


	f,err := os.OpenFile("fixtures/crypto-config/log",os.O_RDWR|os.O_CREATE,0777)
	if err!=nil && os.IsNotExist(err){
		f, err = os.Create("fixtures/crypto-config/log")
		if err != nil {
			panic(err)
		}
	}else{
		return
	}
	defer f.Close()

	f.Write([]byte("finish"))

	// 初始化操作
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