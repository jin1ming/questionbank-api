package v1

import (
	md5 "crypto/md5"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	//连接至数据库
	var err error
	db,err = sql.Open("mysql", DBUserName+":"+DBPassWord+"@tcp("+DBHostsIp+")/"+DBName)
	CheckErr(err)
	// 不关闭数据库
	// defer db.Close()
}

// 用户注册
func registerDb(name string, pwd string, role string)  {
	pwdMd5 := getPwdMd5(pwd)
	stmt,err := db.Prepare("INSERT user (username, pwd_md5, role) values (?,?,?)")
	CheckErr(err)
	//执行插入操作
	_, err = stmt.Exec(name, pwdMd5, role)
	CheckErr(err)
}


// 检查用户名是否存在
func checkName(name string) bool{
	var res string
	err := db.QueryRow("SELECT username FROM user WHERE username=?", name).Scan(&res)
	if err != nil || res != name {
		return false
	}
	return true
}

// 检查用户密码是否正确
func checkPwd(name string, pwd string, role string) bool{
	pwdmd5 := getPwdMd5(pwd)
	var q,r string
	err := db.QueryRow("SELECT pwd_md5, role FROM user WHERE username=?", name).Scan(&q, &r)
	if err != nil || q != pwdmd5 || r != role {
		return false
	}
	return true
}


// 删除用户
func delUserFromDb(name string){
	_, err := db.Exec("DELETE FROM user WHERE username=?", name)
	CheckErr(err)
}

// 添加试卷
func addPaper2Db(title string, owner string, questionIds []string )  {
	stmt,err := db.Prepare("INSERT paper (title, owner) values (?,?)")
	CheckErr(err)
	_, err = stmt.Exec(title, owner)
	CheckErr(err)
	var paperId int64
	err = db.QueryRow("SELECT id FROM paper ORDER by id DESC LIMIT 1").Scan(&paperId)
	CheckErr(err)
	for i := range questionIds{
		stmt,err = db.Prepare("INSERT paper_question (question_id, paper_id) values (?,?)")
		CheckErr(err)
		_, err = stmt.Exec(questionIds[i], paperId)
		CheckErr(err)
	}
}

// 获取所有试卷
func getAllPapersFromDb() (papers []Paper) {
	rows, err := db.Query("SELECT id, title, owner FROM paper")
	CheckErr(err)
	for rows.Next() {
		var p Paper
		if err := rows.Scan(&p.Id, &p.Title, &p.Owner); err == nil {
			papers = append(papers, p)
		}
	}
	return papers
}

// 获取指定试卷的试题id
func getPaperQuestionsFromDb(paperId int64) (questionIds []string) {
	rows, err := db.Query("SELECT question_id FROM paper_question WHERE paper_id=? ", paperId)
	CheckErr(err)
	for rows.Next() {
		var q string
		if err := rows.Scan(&q); err == nil {
			//log.Println("q:",q)
			questionIds = append(questionIds, q)
		}
	}
	//log.Println("questionIds:", questionIds)
	return questionIds
}

// 删除试卷中的某些试题
func delPaperItemFromDb(paperId int64, questionId []string) {
	for i := range questionId {
		stmt,err := db.Prepare("DELETE FROM paper_question WHERE paper_id = ? AND question_id = ?")
		CheckErr(err)
		_, err = stmt.Exec(paperId, questionId[i])
		CheckErr(err)
	}
}

// 为密码生成md5,
func getPwdMd5(pwd string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(pwd)))
}

func CheckErr(err error)  {
	if err != nil {
		panic(err)
	}
}