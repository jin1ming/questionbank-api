package v1

import (
	"encoding/hex"
	mspclient "github.com/hyperledger/fabric-sdk-go/pkg/client/msp"
	"github.com/hyperledger/fabric-sdk-go/pkg/common/providers/msp"
	"log"
	"strings"
)


const (
	Admin = "Admin"
)


func GetKeyFile(id msp.SigningIdentity) (string, string) {
	priFile := hex.EncodeToString(id.PrivateKey().SKI()) + "_sk"
	pubFile := id.Identifier().ID + "@" + id.Identifier().MSPID + "-cert.pem"
	return priFile, pubFile
}

func RegisterUser(userName string, userOrg string, role string) (priFile string, pubFile string, err error) {

	secret := userName + userOrg
	if sdk == nil {
		log.Fatal("sdk is nil!!!!!!!!!!!!!!!!!!")
	}
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(userOrg))
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
		Type: "user",
		Secret: secret,
		Attributes: attris,
		Affiliation: "org1",
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

func DeleteUser(userName string, userOrg string) error {
	if sdk == nil {
		log.Fatal("sdk is nil!!!!!!!!!!!!!!!!!!")
	}
	mspClient, err := mspclient.New(sdk.Context(), mspclient.WithOrg(userOrg))
	if err != nil {
		log.Panicf("Failed to create msp client: %s\n", err)
		return err
	}
	//判断是否存在
	id, err := mspClient.GetSigningIdentity(userName)
	if err != nil {
		log.Println("user not exists: ", userName)
		return err
	}
}
