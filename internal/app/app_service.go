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
	DeleteById(id string, data *model.FormData) (errs []string)
	Save(id string, data *model.FormData) (errs []string)
	Create(data *model.FormData) (errs []string)
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

func (r *appService) DeleteById(id string, data *model.FormData) (errs []string) {
	if data.Password == "" {
		errs = append(errs, "Password is required")
		return
	}

	var resultSet *model.FormData
	if _, errp := r.DbClient.From("notes").Select("password", "", false).Eq("id", id).Single().ExecuteTo(&resultSet); errp != nil {
		errs = append(errs, errp.Error())
	} else {
		if !CheckPasswordHash(data.Password, resultSet.Password) {
			errs = append(errs, "Invalid password")
		}
	}

	if len(errs) == 0 {
		r.DbClient.From("notes").Delete("", "").Eq("id", id).Execute()
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

func (r *appService) Save(id string, data *model.FormData) (errs []string) {
	if data.Password == "" {
		errs = append(errs, "Password is required")
		return
	}

	var resultSet *model.FormData
	if _, errp := r.DbClient.From("notes").Select("password", "", false).Eq("id", id).Single().ExecuteTo(&resultSet); errp != nil {
		errs = append(errs, errp.Error())
	} else {
		if !CheckPasswordHash(data.Password, resultSet.Password) {
			errs = append(errs, "Invalid password")
		}
	}

	if len(errs) == 0 {
		if data.Id == "" {
			data.Id = id
		} else if id != data.Id {
			go func() {
				r.DbClient.From("notes").Delete("", "").Eq("id", id).Execute()
			}()
		}
		data.Password = resultSet.Password
		r.DbClient.From("notes").Upsert(&data, "", "", "").Execute()

	}
	return
}

func (r *appService) Create(data *model.FormData) (errs []string) {
	if data.Password == "" {
		errs = append(errs, "Password is required")
		return
	}

	data.Password = HashPassword(data.Password)

	if data.Id == "" {
		data.Id, _ = GenerateRandomId(5)
	}

	if len(errs) == 0 {
		a, _, _ := r.DbClient.From("notes").Insert(&data, false, "", "", "").Execute()
		fmt.Println(a)
	}
	return
}
