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
	Vcode    string `form:"vcode"`
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

type LogoutRequest struct {
	Token string `header:"Authorization"`
}

type LogoutResponse struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

type RegisterVcodeRequest struct {
	Email string `form:"email"`
}

type RegisterVcodeResponse struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

type FileUploadResponse struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

type FileUploadRequest struct {
	Parent *uint `form:"parent,optional"`
}

type CreateDirectoryResponse struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

type DeleteFileResponse struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

type DeleteFileRequest struct {
	UUID string `path:"uuid"`
}

type DeleteFilesRequest struct {
	UUIDs []string `form:"uuids"`
}

type DeleteFilesResponse struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

type CreateDirectoryRequest struct {
	Name   string `form:"name,optional"`
	Parent *uint  `form:"parent,optional"`
}

type UpdateFileRequest struct {
	ID      uint   `path:"id"`
	Name    string `form:"name,optional"`
	Private *bool  `form:"private,optional"`
	Parent  *uint  `form:"parent,optional"`
}

type UpdateFileResponse struct {
	Message string `json:"message"`
	Ok      bool   `json:"ok"`
}

type FileDownloadRequest struct {
	UUID string `path:"uuid"`
}

type UserDTO struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type FileDTO struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Ext     string `json:"ext"`
	Size    int64  `json:"size"`
	UUID    string `json:"uuid"`
	Owner   uint   `json:"owner"`
	IsDir   bool   `json:"isDir"`
	Private bool   `json:"private"`
	Parent  *uint  `json:"parent"`
}

type GetMyFileRequest struct {
	Parent     *uint    `form:"parent,optional"`
	Keyword    string   `form:"keyword,optional"`
	Extensions []string `form:"extensions,optional"`
}

type GetOneFileRequest struct {
	ID uint `path:"id"`
}
