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
func putQuesionCache(stub shim.ChaincodeStubInterface, Q *Question) error {
	qBytes, err := json.Marshal(Q)
	if err != nil {
		return err
	}
	err = stub.PutState("cache_put_" + Q.Id, qBytes)
	if err != nil {
		return nil
	}
	return nil
}

// 放入待删除缓存
func putDelCache(stub shim.ChaincodeStubInterface, question_id string ) error {
	qBytes, err := stub.GetState(question_id)

	if err != nil {
		return err
	}
	err = stub.PutState("cache_del_" + question_id[9:], qBytes)
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
