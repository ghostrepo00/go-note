package model

type FormData struct {
	Id          string `form:"id" json:"id"`
	Content     string `form:"content" json:"content"`
	Password    string `form:"password" json:"password"`
	IsEncrypted bool   `form:"isEncrypted" json:"is_encrypted"`
}
