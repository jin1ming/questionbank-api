package v1


type Paper struct {
	Id int64		`json:"id"`
	Title string	`json:"title"`
	Owner string	`json:"owner"`
}
