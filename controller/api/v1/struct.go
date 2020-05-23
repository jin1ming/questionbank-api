package v1

type Paper struct {
	Id int64		`json:"id"`
	Title string	`json:"title"`
	Owner string	`json:"owner"`
}

type User struct {
	Name	string	`json:"name"`
	Role    string	`json:"role"`
}