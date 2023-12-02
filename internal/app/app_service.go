package app

import (
	"crypto/rand"
	"errors"
	"fmt"
	"log/slog"
	"math/big"

	"github.com/ghostrepo00/go-note/config"
	"github.com/ghostrepo00/go-note/internal/pkg/model"
	"github.com/supabase-community/supabase-go"
)

type appService struct {
	AppConfig *config.AppConfig
	DbClient  *supabase.Client
}

type AppService interface {
	GetbyId(id string) (result []*model.FormData, err error)
	DeletebyId(id string) (err error)
	Save(data *model.FormData) error
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
		_, err = r.DbClient.From("notes").Select("id, content, password, is_encrypted", "", false).Eq("id", id).ExecuteTo(&result)
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

func (r *appService) Save(data *model.FormData) (err error) {
	if data.Id == "" {
		if data.Id, err = GenerateRandomId(5); err == nil {
			a, cc, d := r.DbClient.From("notes").Select("id", "", false).Eq("id", data.Id).ExecuteString()
			slog.Info("", "a", a, "e", d, "c", cc)
		}
	}

	if a, b, err := r.DbClient.From("notes").Upsert(&data, "", "", "").Execute(); err != nil {
		return errors.New("ssssssssssssssssssssssssss")
	} else {
		slog.Info("supabase", "a", a, "b", b, "c", err)
	}
	return nil
}
