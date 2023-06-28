package interfaces

type UserLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Forever  bool   `json:"forever"`
}
