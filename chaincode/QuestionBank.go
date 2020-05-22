package main

import (
	"encoding/json"
	"fmt"
	"github.com/hyperledger/fabric/core/chaincode/shim"
	pb "github.com/hyperledger/fabric/protos/peer"
)

type QuestionBank struct{}

func initLog(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	err := stub.PutState("num_of_logs", []byte("0"))
	if err != nil {
		panic(err)
	}
	return shim.Success(nil)
}

func putQuestion(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 参数校验
	if len(args) != 3 {
		return shim.Error("not enough args")
	}

	name := args[0]
	question := args[1]
	answer := args[2]

	if name == "" || question == "" {
		return shim.Error("invalid args")
	}

	question_id := getNextQuestionId(stub)

	// 存入待审核事件
	Q := &Question{
		Owner:  name,
		Id:     question_id,
		Data:   question,
		Answer: answer,
	}

	err := putQuesionCache(stub, Q)
	if err != nil {
		return shim.Error(fmt.Sprintf("put question cache err! %s", err))
	}

	// 记录到日志
	err = addLog(stub, name, "putQuestion", question_id)
	if err != nil {
		panic(err)
	}
	return shim.Success(nil)
}

func getQuestion(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数
	if len(args) != 2 {
		return shim.Error("not enough args")
	}

	name := args[0]
	question_id := args[1]
	if name == "" || question_id == "" {
		return shim.Error("invalid args")
	}

	Prefix := "question_"
	// 根据id获取试题
	/*
	if 0 == strings.Count(question_id, "put_cache_"){
		Prefix = ""
	}
	*/
	qBytes, err := stub.GetState(Prefix + question_id)
	if err != nil || qBytes == nil {
		return shim.Error("not find question")
	}

	// 获取用户角色
	role,ok := getRole(stub)
	if !ok {
		return shim.Error("get role error !")
	}

	// 根据角色确定是否返回数据
	switch role {
	case Teacher:
		return shim.Success(qBytes)
	case Student:
		Q := new(Question)
		err = json.Unmarshal(qBytes, Q)
		if err != nil {
			return shim.Error(fmt.Sprintf("unmarshal question error! %s", err))
		}
		QN := &QuestionN{
			Owner:	Q.Owner,
			Id:    	Q.Id,
			Data:  	Q.Data,
		}
		qnBytes, err := json.Marshal(QN)
		if err != nil {
			return shim.Error(fmt.Sprintf("unmarshal questionN error! %s", err))
		}
		return shim.Success(qnBytes)
	case Reviewer:
		return shim.Success(qBytes)
	case Admin:
		//TODO 暂时不对管理员开放题库
		fallthrough
	default:
		return shim.Error("user role error!  ")
	}
}

//获取所有试题
func getAllQuestions(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数
	if len(args) != 1 {
		return shim.Error("not enough args")
	}

	name := args[0]
	if name == "" {
		return shim.Error("invalid args")
	}

	// 获取用户角色
	role,ok := getRole(stub)
	if !ok {
		return shim.Error("get role error !")
	}

	// 根据角色确定是否返回数据
	switch {
	case role == Student:
		var questions []QuestionN
		it, err := stub.GetStateByRange("questiom","questioo")
		defer it.Close()
		if err != nil {
			return shim.Error(fmt.Sprintf("get questions error ! %s", err))
		}
		for it.HasNext() {
			it, _ := it.Next()
			Q := new(Question)
			err = json.Unmarshal(it.Value,Q)
			if err != nil {
				return shim.Error(fmt.Sprintf("unmarshal question error! %s", err))
			}
			QN := &QuestionN{
				Owner:	Q.Owner,
				Id:    	Q.Id,
				Data:  	Q.Data,
			}
			questions = append(questions, *QN)
		}
		questionsBytes, err := json.Marshal(questions)
		if err != nil {
			return shim.Error(fmt.Sprintf("marshal questions error !", err))
		}
		return shim.Success(questionsBytes)
	case role == Reviewer || role ==Teacher :
		var questions []Question
		it, err := stub.GetStateByRange("questiom","questioo")
		defer it.Close()
		if err != nil {
			return shim.Error(fmt.Sprintf("get questions error ! %s", err))
		}
		for it.HasNext() {
			it, _ := it.Next()
			Q := new(Question)
			err = json.Unmarshal(it.Value,Q)
			if err != nil {
				return shim.Error(fmt.Sprintf("unmarshal question error! %s", err))
			}
			questions = append(questions, *Q)
		}
		questionsBytes, err := json.Marshal(questions)
		if err != nil {
			return shim.Error(fmt.Sprintf("marshal questions error !", err))
		}
		return shim.Success(questionsBytes)
	case role == Admin:
		//TODO 暂时不对管理员开放题库
		fallthrough
	default:
		return shim.Error("user role error!  ")
	}
}

func delQuestion(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数
	if len(args) != 2 {
		return shim.Error("not enough args")
	}

	name := args[0]
	question_id := args[1]
	if name == "" || question_id == "" {
		return shim.Error("invalid args")
	}

	//TODO 根据角色或者出题人做出约束，可切换
	role,ok := getRole(stub)
	if !ok {
		return shim.Error("get role error !")
	}
	if role != Teacher {
		return shim.Error(fmt.Sprintf("this role : %s  not is teacher! ", role))
	}
	if err := putDelCache(stub, question_id); err != nil {
		return shim.Error(fmt.Sprintf("put del cache error ! %s", err))
	}

	// 记录到日志
	err := addLog(stub, name, "delQuestion", question_id)
	if err != nil {
		panic(err)
	}

	return shim.Success(nil)
}

func putScore(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数
	if len(args) != 3 {
		return shim.Error("not enough args")
	}

	name := args[0]
	exam := args[1]
	score := args[2]
	if name == "" || exam == "" || score == ""{
		return shim.Error("invalid args")
	}

	// 尝试获取ExamScore
	esBytes, err := stub.GetState(conExamScoreId(exam))
	if err != nil {
		return shim.Error(fmt.Sprintf("get exam score error ! %s", err))
	}

	var es *ExamScore

	if esBytes == nil { // 该ExamScore未创建，先创建
		es = &ExamScore{
			Exam:  exam,
			Score: make(map[string]string),
		}
	} else {			// 已存在ExamScore，转为结构体
		err = json.Unmarshal(esBytes, es)
		if err != nil {
			return shim.Error(fmt.Sprintf("unmarshal esBytes error ! %s", err))
		}
	}
	// 存入分数
	es.Score[name] = score
	nes, err := json.Marshal(es)
	err = stub.PutState(conExamScoreId(exam), nes)
	if err != nil {
		return shim.Error(fmt.Sprintf("put exam score error ! %s", err))
	}
	return shim.Success(nil)
}

func getScore(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数
	if len(args) != 3 {
		return shim.Error("not enough args")
	}

	name := args[0]		//发起查询者
	stu_name := args[1]	//被查询对象
	exam := args[2]
	if name =="" || stu_name == "" || exam == "" {
		return shim.Error("invalid args")
	}

	// 检查权限
	role,ok := getRole(stub)
	if !ok {
		return shim.Error("get role error !")
	}
	if name != stu_name && role != Teacher {
		return shim.Error("not right to score !")
	}

	// 获取成绩
	esBytes, err := stub.GetState(conExamScoreId(exam))
	if err != nil {
		return shim.Error(fmt.Sprintf("get exam score error ! %s", err))
	}
	examScore := new(ExamScore)
	err = json.Unmarshal(esBytes, examScore)
	if err != nil {
		return shim.Error(fmt.Sprintf("unmarshal exam score error ! %s", err))
	}
	score := examScore.Score[stu_name]
	if score == "" {
		return shim.Error("not find score !")
	}
	return shim.Success([]byte(score))
}

func getLogs(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查权限
	role,ok := getRole(stub)
	if !ok {
		return shim.Error("get role error !")
	}
	if role != Admin {
		return shim.Error("not right to get logs !")
	}
	var logs []json.RawMessage
	it, err := stub.GetStateByRange("lof","loh")
	defer it.Close()
	if err != nil {
		return shim.Error(fmt.Sprintf("get logs error ! %s", err))
	}
	for it.HasNext() {
		it, _ := it.Next()
		logs = append(logs, it.Value)
		/*
		L := new(Log)
		err = json.Unmarshal(it.Value,L)
		if err != nil {
			return shim.Error(fmt.Sprintf("unmarshal question error! %s", err))
		}
		logs.Data[it.Key] = *L
		*/
	}
	logsBytes, err := json.Marshal(logs)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal logs error !", err))
	}
	return shim.Success(logsBytes)
}

func getCache(stub shim.ChaincodeStubInterface) pb.Response {

	// 检查权限
	role,ok := getRole(stub)
	if !ok {
		return shim.Error("get role error !")
	}
	if role != Reviewer {
		return shim.Error("not right to score !")
	}

	/*
	list := new(ListCache)
	list.DelCache = make(map[string]Question)
	list.PutCache = make(map[string]Question)
	 */

	var list ListCache
	itPut, err := stub.GetStateByRange("cache_p","cache_q")
	defer itPut.Close()
	if err != nil {
		return shim.Error(fmt.Sprintf("get put cache list error ! %s", err))
	}
	for itPut.HasNext() {
		it, _ := itPut.Next()
		list.PutCache = append(list.PutCache, it.Value)
		/*
		Q := new(Question)
		err = json.Unmarshal(it.Value,Q)
		if err != nil {
			return shim.Error(fmt.Sprintf("unmarshal question error! %s", err))
		}
		list.PutCache[it.Key] = *Q
		 */
	}

	itDel, err := stub.GetStateByRange("cache_d","cache_e")
	defer itDel.Close()
	if err != nil {
		return shim.Error(fmt.Sprintf("get del cache list  error ! %s", err))
	}
	for itDel.HasNext() {
		it, _ := itDel.Next()
		list.DelCache = append(list.DelCache, it.Value)
		/*
		Q := new(Question)
		err = json.Unmarshal(it.Value,Q)
		if err != nil {
			return shim.Error(fmt.Sprintf("unmarshal question error! %s", err))
		}
		list.DelCache[it.Key] = *Q
		 */
	}

	listBytes, err := json.Marshal(list)
	if err != nil {
		return shim.Error(fmt.Sprintf("marshal list error !", err))
	}
	//log.Println("listBytes:",string(listBytes))
	return shim.Success(listBytes)
}

func approve(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数
	if len(args) != 3 {
		return shim.Error("not enough args")
	}

	name := args[0]
	op := args[1]
	question_id := args[2]
	if name == "" || op == "" || question_id == "" {
		return shim.Error("invalid args")
	}

	// 检查权限
	role,ok := getRole(stub)
	if !ok {
		return shim.Error("get role error !")
	}
	if role != Reviewer {
		return shim.Error("not right to approve !")
	}
	switch op {
	case Put:
		// 从缓存取出
		qBytes, err := stub.GetState(question_id)
		if err != nil {
			return shim.Error(fmt.Sprintf("get question error ! %s", err))
		}

		// 写入状态
		err = stub.PutState("question_" + question_id[10:], qBytes)
		if err != nil {
			return shim.Error(fmt.Sprintf("put question error ! %s", err))
		}

		// 从缓存删除
		err = stub.DelState(question_id)
		if err != nil {
			return shim.Error(fmt.Sprintf("del question error ! %s", err))
		}
		// 记录到日志
		err = addLog(stub, name, "approve putQuestion", question_id)
		if err != nil {
			panic(err)
		}
		return shim.Success(nil)
	case Delete:
		// 从缓存删除
		err := stub.DelState(question_id)
		if err != nil {
			return shim.Error(fmt.Sprintf("del question error ! %s", err))
		}
		err = stub.DelState("cache_del_" + question_id[10:])
		if err != nil {
			return shim.Error(fmt.Sprintf("del question error ! %s", err))
		}

		// 记录到日志
		err = addLog(stub, name, "approve delQuestion", question_id)
		if err != nil {
			panic(err)
		}
		return shim.Success(nil)
	default:
		return shim.Error("unknow op !")
	}
}

func reject(stub shim.ChaincodeStubInterface, args []string) pb.Response {
	// 检查参数
	if len(args) != 3 {
		return shim.Error("not enough args")
	}

	name := args[0]
	op := args[1]
	question_id := args[2]
	if op == "" || question_id == "" {
		return shim.Error("invalid args")
	}

	// 检查权限
	role,ok := getRole(stub)
	if !ok {
		return shim.Error("get role error !")
	}
	if role != Reviewer {
		return shim.Error("not right to approve !")
	}
	switch op {
	case Put:
		// 从缓存删除
		err := stub.DelState(question_id)
		if err != nil {
			return shim.Error(fmt.Sprintf("del question error ! %s", err))
		}

		// 记录到日志
		err = addLog(stub, name, "reject putQuestion", question_id)
		if err != nil {
			panic(err)
		}
		return shim.Success(nil)
	case Delete:
		err := stub.DelState("cache_del_" + question_id[10:])
		if err != nil {
			return shim.Error(fmt.Sprintf("del question error ! %s", err))
		}

		// 记录到日志
		err = addLog(stub, name, "reject delQuestion", question_id)
		if err != nil {
			panic(err)
		}
		return shim.Success(nil)
	default:
		return shim.Error("unknow op !")
	}
}
