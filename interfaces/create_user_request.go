package interfaces

var CreateUserRequest struct {
	Name    string `json:"name"`
	Email   string `json:"email"`
	Job     string `json:"job"`
	JobDesc string `json:"jobdesc"`
}
