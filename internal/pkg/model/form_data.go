package model

type FormData struct {
	Id          string `form:"id" json:"id,omitempty"`
	Content     string `form:"content" json:"content"`
	Password    string `form:"password" json:"password"`
	IsEncrypted bool   `form:"is_encrypted" json:"is_encrypted"`
}
