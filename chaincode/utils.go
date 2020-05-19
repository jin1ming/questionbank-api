package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/lib/cid"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	"strconv"
	"time"
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

// 写入日志
func addLog(stub shim.ChaincodeStubInterface, name string, op string, question_id string) error {
	t := time.Now().Format("2006-01-02 15:04")
	idBytes, err := stub.GetState("num_logs")
	if err != nil {
		panic(err)
	}
	idStr := ""
	if idBytes == nil {
		idStr = "1"
	} else {
		id, err := strconv.Atoi(string(idBytes))
		if err != nil {
			panic(err)
		}
		idStr = strconv.Itoa(id + 1)
	}
	err = stub.PutState("num_logs", []byte(idStr))
	if err != nil {
		panic(err)
	}
	l := Log{
		Id:         idStr,
		Time:       t,
		Name:       name,
		Op:         op,
		QuestionId: question_id,
	}
	lBytes, err := json.Marshal(l)
	if err != nil {
		return err
	}
	err = stub.PutState("log_" + idStr, lBytes)
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

func getNextQuestionId(stub shim.ChaincodeStubInterface)(id string) {
	idBytes, err := stub.GetState("num_question")
	if err != nil {
		panic(err)
	}
	if idBytes == nil {
		id = "1"
	}
	idInt, err := strconv.Atoi(string(idBytes))
	if err != nil {
		panic(err)
	}
	idNext := strconv.Itoa(idInt + 1)
	err = stub.PutState("num_question", []byte(idNext))
	if err != nil {
		panic(err)
	}
	return
}