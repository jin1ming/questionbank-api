package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
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
func getRole(stub shim.ChaincodeStubInterface) (role string, ok bool) {
	role, ok, err := cid.GetAttributeValue(stub, "Role")
	if err != nil || !ok {
		return "", false
	}
	return
}
