package main

import(
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

// 暂时初始化无内容
func (c *QuestionBank) Init(stub shim.ChaincodeStubInterface) pb.Response {

	return shim.Success(nil)
}

// 调用入口
func (c *QuestionBank) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	funcName, args := stub.GetFunctionAndParameters()

	switch funcName {
	case "putQuestion":	// 存入试题（需审核）
		return putQuestion(stub, args)
	case "getQuestion":	// 获取试题
		return getQuestion(stub, args)
	case "delQuestion":	// 删除试题（需审核）
		return delQuestion(stub, args)
	case "putScore":	// 添加学生成绩
		return putScore(stub, args)
	case "getScore":	// 获取成绩
		return getScore(stub, args)
	case "getLogs":		// 获取日志
		return getLogs(stub, args)
	case "userRegister":// 注册用户
		return userRegister(stub, args)
	case "userDestroy":	// 注销用户
		return userDestroy(stub, args)
	case "getCache":
		return getCache(stub, args)
	case "approve":		// 通过审核
		return approve(stub, args)
	case "reject":		// 审核失败
		return reject(stub, args)
	default:
		return shim.Error("unknow function !")
	}
}

func main(){
	err := shim.Start(new(QuestionBank))
	if err != nil {
		fmt.Printf("Error starting AssertsExchange chaincode: %s", err)
	}
}