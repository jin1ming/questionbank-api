package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
)


// 构造考试成绩id
func conExamScoreId(exam string) string {
	return "exam_score_" + exam
}

// 构造用户key
func constructUserKey(name string) string {
	return fmt.Sprintf("user_%s", name)
}

// 放入待审核缓存
func putQuesionCache(stub shim.ChaincodeStubInterface, q *Question) error {
	Q, err := json.Marshal(q)
	if err != nil {
		return err
	}
	stub.PutState("cache_put_" + q.Id, Q)
	return nil
}

// 放入待删除缓存
func putDelCache(stub shim.ChaincodeStubInterface, id string) error {
	err := stub.PutState("cache_del_" + id, []byte(id))
	if err != nil {
		return err
	}
	return nil
}

// 获取用户身份
func getRole(stub shim.ChaincodeStubInterface, name string) (role string, err error) {
	key := "user_" + name

	userBytes, err := stub.GetState(key)
	if err != nil || len(userBytes) == 0 {
		return role,fmt.Errorf("user not found")
	}

	user := new(User)
	err = json.Unmarshal(userBytes, user)
	if err != nil {
		return
	}

	role = user.Role
	return
}