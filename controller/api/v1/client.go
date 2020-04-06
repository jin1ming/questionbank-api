package v1

import (
	"encoding/hex"
	"github.com/hyperledger/fabric-sdk-go/pkg/client/channel"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/errors/retry"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"log"
	"strings"
)



func GetKeyFile(id msp.SigningIdentity) (string, string) {
	priFile := hex.EncodeToString(id.PrivateKey().SKI()) + "_sk"
	pubFile := id.Identifier().ID + "@" + id.Identifier().MSPID + "-cert.pem"
	return priFile, pubFile
}

func RegisterUser(userName string, role string) (priFile string, pubFile string, err error) {

	secret := userName + orgName
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(orgName))
	if err != nil {
		log.Panicf("Failed to create msp client: %s\n", err)
		return
	}
	//判断是否存在
	id, err := mspClient.GetSigningIdentity(userName)
	if err == nil {
		log.Println("user exists: ", userName)
		priFile, pubFile = GetKeyFile(id)
		return
	}
	//注册用户
	var attris []mspclient.Attribute
	attri := mspclient.Attribute{
		Name:  "Role",
		Value: role,
		ECert: true,
	}
	attris = append(attris, attri)

	request := &mspclient.RegistrationRequest{
		Name: userName,
		Type: "client",
		Secret: secret,
		Attributes: attris,
		Affiliation: orgName,
		MaxEnrollments: 10,
		CAName: "ca.org1.questionbank.com",
	}
	_, err = mspClient.Register(request)
	if err != nil && !strings.Contains(err.Error(), "is already registered") {
		log.Fatalf("register %s [%s]\n", userName, err)
		return
	}
	//登记保存证书到stores
	err = mspClient.Enroll(userName, mspclient.WithSecret(secret))
	if err != nil {
		log.Panicf("Failed to enroll user: %s\n", err)
		return
	}

	id, _ = mspClient.GetSigningIdentity(userName)
	priFile, pubFile = GetKeyFile(id)
	log.Printf("register %s successfully\n", userName)
	return
}

func GetId(userName string, userOrg string) (*mspclient.IdentityResponse, error) {
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(userOrg))
	if err != nil {
		log.Panicf("Failed to create msp client: %s\n", err)
		return nil, err
	}
	//判断是否存在
	id, err := mspClient.GetIdentity(userName)
	if err != nil {
		return nil, err
	}
	return id, nil
}

func RemoveUser(userName string, userOrg string) error {
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(userOrg))
	if err != nil {
		log.Panicf("Failed to create msp client: %s\n", err)
		return err
	}
	//判断是否存在
	_, err = mspClient.GetSigningIdentity(userName)
	if err != nil {
		log.Println("user not exists: ", userName)
		return err
	}
	request := mspclient.RemoveIdentityRequest{
		ID:     userName,
		Force:  false,
		CAName: "ca.org1.questionbank.com",
	}
	_, err  = mspClient.RemoveIdentity(&request)
	if err != nil {
		log.Println("remove identify failed !")
		return err
	}
	return nil
}

func queryCC(client *channel.Client, fcn string, queryArgs [][]byte) ([]byte, error) {
	response, err := client.Query(channel.Request{ChaincodeID: ccID, Fcn: fcn, Args: queryArgs},
		channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargetEndpoints("grpcs://peer0.org1.questionbank.com:7051"))
	if err != nil {
		log.Println("Failed to query: %s", err)
		return nil, err
	}
	log.Println(response)

	return response.Payload, nil
}

func executeCC(client *channel.Client, fcn string, txArgs [][]byte) error {
	_, err := client.Execute(channel.Request{ChaincodeID: ccID, Fcn: fcn, Args: txArgs},
		channel.WithRetry(retry.DefaultChannelOpts), channel.WithTargetEndpoints("grpcs://peer0.org1.questionbank.com:7051"))
	if err != nil {
		log.Println("Failed to execute: %s", err)
		return err
	}
	return nil
}