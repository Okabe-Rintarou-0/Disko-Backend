// Code generated by goctl. DO NOT EDIT.
package types

type LoginRequest struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

type RegisterRequest struct {
	Name     string `form:"name"`
	Email    string `form:"email"`
	Password string `form:"password"`
}

type Response struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type LoginResponse struct {
	Message  string `json:"message"`
	Ok       bool   `json:"ok"`
	Token    string `json:"token"`
	ExpireAt int64  `json:"expireAt"`
}

type RegisterResponse struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

type UserDTO struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
