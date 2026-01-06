package model

// RegisterReq is a request model for user registration.
type RegisterReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// RegisterResp is a response model for user registration.
type RegisterResp struct {
	ID    int64  `json:"id"`
	Token string `json:"token"`
}

// LoginReq is a request model for user login.
type LoginReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// LoginResp is a response model for user login.
type LoginResp struct {
	Token string `json:"token"`
}
