package interfaces

type CreateUserResponse struct {
	ID     uint   `json:"id"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	Job    string `json:"job"`
	JobDes string `json:"jobDesc"`
}
