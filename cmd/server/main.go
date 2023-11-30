package main

import (
	"io"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/ghostrepo00/go-note/config"
	"github.com/ghostrepo00/go-note/internal/app"
	appconstant "github.com/ghostrepo00/go-note/internal/pkg/app_constant"
)

func getLogFileName(appConfig *config.AppConfig) (result string) {
	currentDate := time.Now()
	return strings.Join([]string{appConfig.Log.Path, "/", currentDate.Format(appconstant.TimestampFormat), "_", appConfig.Log.FileName}, "")
}

func useSlog(appConfig *config.AppConfig) (logFile *os.File, err error) {
	logFile, err = os.OpenFile(getLogFileName(appConfig), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)

	slogHandler := slog.NewJSONHandler(
		io.MultiWriter(os.Stdout, logFile),
		&slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		})
	logger := slog.New(slogHandler)
	slog.SetDefault(logger)
	return
}

func main() {
	if config, err := config.NewAppConfig("./config/config.json"); err != nil {
		panic(err)
	} else {
		if fileLog, err := useSlog(config); err == nil {

			defer fileLog.Close()

			slog.Info("App started")
			config.SupabaseUrl = os.Getenv("SUPABASE_URL")
			config.SupabaseKey = os.Getenv("SUPABASE_KEY")

			webserver := app.NewWebServer(config)
			webserver.Run()

		} else {
			panic(err)
		}
	}
}
