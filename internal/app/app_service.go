package app

import (
	"github.com/ghostrepo00/go-note/config"
	"github.com/supabase-community/supabase-go"
)

type appService struct {
	AppConfig *config.AppConfig
	DbClient  *supabase.Client
}

func NewAppService(appConfig *config.AppConfig, dbClient *supabase.Client) *appService {
	return &appService{appConfig, dbClient}
}

func (r *appService) do() {

}
