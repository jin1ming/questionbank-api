package v1

// 用于sdk初始化、创建通道、安装链码
import (
	"fmt"
	"questionbank-api/sdkInit"
)

func InitSdk(initInfo *sdkInit.InitInfo, configFile string) {

	sdk, err := sdkInit.SetupSDK(configFile)
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