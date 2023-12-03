package app

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"

	"github.com/ghostrepo00/go-note/config"
	"github.com/ghostrepo00/go-note/internal/pkg/model"
	"github.com/supabase-community/supabase-go"
	"golang.org/x/crypto/bcrypt"
)

type appService struct {
	AppConfig *config.AppConfig
	DbClient  *supabase.Client
}

type AppService interface {
	GetbyId(id string) (result []*model.FormData, err error)
	DeletebyId(id string) (err error)
	Save(id string, data *model.FormData) (outputPassword string, err error)
}

func NewAppService(appConfig *config.AppConfig, dbClient *supabase.Client) *appService {
	return &appService{appConfig, dbClient}
}

func GenerateRandomId(length int) (string, error) {
	const letters = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	result := make([]byte, length)
	for i := 0; i < length; i++ {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(letters))))
		if err != nil {
			return "", err
		}
		result[i] = letters[num.Int64()]
	}

	return string(result), nil
}

func (r *appService) GetbyId(id string) (result []*model.FormData, err error) {
	if id != "" {
		_, err = r.DbClient.From("notes").Select("id, content, is_encrypted", "", false).Eq("id", id).ExecuteTo(&result)
		fmt.Print(result, err)
	}
	return
}

func (r *appService) DeletebyId(id string) (err error) {
	if id != "" {
		_, _, err = r.DbClient.From("notes").Delete("", "").Eq("id", id).Execute()
	}
	return
}

func (r *appService) GenerateNewId() (result string, err error) {
	for i := 0; i < 3; i++ {
		if result, err = GenerateRandomId(5); err == nil {
			a, _, _ := r.DbClient.From("notes").Select("id", "", false).Eq("id", result).Single().ExecuteString()
			if len(a) == 0 {
				return
			}
		}
	}

	return "", errors.New("Failed to generate new id")
}

func HashPassword(password string) string {
	bytes, _ := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes)
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func (r *appService) Save(id string, data *model.FormData) (outputPassword string, err error) {

	var hashedPass string
	if data.Password != "" {
		hashedPass = HashPassword(data.Password)
	} else if id == "" {
		p, _ := GenerateRandomId(5)
		hashedPass = HashPassword(p)
	}

	if id == "" {
		id, err = r.GenerateNewId()
	} else {
		var resultSet *model.FormData
		_, err = r.DbClient.From("notes").Select("password", "", false).Eq("id", id).Single().ExecuteTo(&resultSet)
		if err != nil {
			return
		} else {
			if !CheckPasswordHash(resultSet.Password, hashedPass) {
				err = errors.New("Invalid password")
				return
			}
		}
	}

	data.Id = id
	data.Password = hashedPass
	_, _, err = r.DbClient.From("notes").Upsert(&data, "", "", "").Execute()

	return
}
