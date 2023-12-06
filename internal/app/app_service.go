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
	AppConfig    *config.AppConfig
	DbClient     *supabase.Client
	CryptoClient CryptoService
}

type AppService interface {
	GetbyId(id string) (result []*model.FormData, err error)
	DeleteById(id string, data *model.FormData) (errs []error)
	Save(id string, data *model.FormData) (errs []error)
	Create(data *model.FormData) (errs []error)
	EncryptMessage(id string, data *model.FormData) (errs []error)
	DecryptMessage(id string, data *model.FormData) (errs []error)
}

func NewAppService(appConfig *config.AppConfig, dbClient *supabase.Client, crypto CryptoService) *appService {
	return &appService{appConfig, dbClient, crypto}
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

func (r *appService) ValidatePassword(id string, inputPassword string) (initialPassword string, err error) {
	var resultSet *model.FormData
	r.DbClient.From("notes").Select("password", "", false).Eq("id", id).Single().ExecuteTo(&resultSet)
	if !CheckPasswordHash(inputPassword, resultSet.Password) {
		return "", errors.New("Invalid Password")
	}
	return resultSet.Password, nil
}

func (r *appService) GetbyId(id string) (result []*model.FormData, err error) {
	_, err = r.DbClient.From("notes").Select("id, content, is_encrypted", "", false).Eq("id", id).ExecuteTo(&result)
	return
}

func (r *appService) DeleteById(id string, data *model.FormData) (errs []error) {
	if _, err := r.ValidatePassword(id, data.Password); err == nil {
		r.DbClient.From("notes").Delete("", "").Eq("id", id).Execute()
	} else {
		errs = append(errs, err)
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

func (r *appService) CheckDuplicateId(id string) error {
	var found []*model.FormData
	r.DbClient.From("notes").Select("id", "", false).Eq("id", id).ExecuteTo(&found)
	if len(found) > 0 {
		return errors.New(fmt.Sprintf("Id \"%s\" has been used by existing record", id))
	}
	return nil
}

func (r *appService) Save(id string, data *model.FormData) (errs []error) {
	if initialPassword, err := r.ValidatePassword(id, data.Password); err == nil {
		if data.Id == "" {
			data.Id = id
		} else if id != data.Id {
			if err := r.CheckDuplicateId(data.Id); err != nil {
				errs = append(errs, err)
				return
			}
			go func() {
				r.DbClient.From("notes").Delete("", "").Eq("id", id).Execute()
			}()
		}
		data.Password = initialPassword
		r.DbClient.From("notes").Upsert(&data, "", "", "").Execute()
	} else {
		errs = append(errs, err)
	}
	return
}

func (r *appService) Create(data *model.FormData) (errs []error) {
	if data.Id == "" {
		data.Id, _ = GenerateRandomId(5)
	} else if err := r.CheckDuplicateId(data.Id); err != nil {
		errs = append(errs, err)
		return
	}

	data.Password = HashPassword(data.Password)

	if _, _, err := r.DbClient.From("notes").Insert(&data, false, "", "", "").Execute(); err != nil {
		errs = append(errs, err)
	}

	return
}

func (r *appService) EncryptMessage(id string, data *model.FormData) (errs []error) {
	if id != "" {
		if _, err := r.ValidatePassword(id, data.Password); err != nil {
			errs = append(errs, err)
			return
		}
	}

	data.IsEncrypted = true
	data.Content, _ = r.CryptoClient.Encrypt(data.Content, data.Password)
	return nil
}

func (r *appService) DecryptMessage(id string, data *model.FormData) (errs []error) {
	if id != "" {
		if _, err := r.ValidatePassword(id, data.Password); err != nil {
			errs = append(errs, err)
			return
		}
	}

	data.IsEncrypted = false
	data.Content, _ = r.CryptoClient.Decrypt(data.Content, data.Password)
	return nil
}
