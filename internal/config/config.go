package config

import (
	"strings"
	"sync"
	"time"

	"github.com/marcos-wz/capstone-go-bootcamp/internal/logger"

	"github.com/spf13/viper"
)

var (
	cfgInstance = newConfig()
	once        sync.Once
)

// Config is the representation of the application's configuration.
type Config struct {
	Application Application
	HTTP        HTTP
	Database    Database
}

// GetInstance returns the default configuration instance.
func GetInstance() Config {
	return cfgInstance
}

func setDefaultConfig() {
	viper.SetDefault("application.version", "v0.0.0")
	viper.SetDefault("http.server.host", "")
	viper.SetDefault("http.server.port", 8080)
	viper.SetDefault("http.server.shutdown.timeout", time.Second*15)
	viper.SetDefault("http.data_api.url", "https://thecocktaildb.com/api/json/v1/1/search.php?f=a")
	viper.SetDefault("database.driver", "csv")
	viper.SetDefault("database.csv.file_name", "cocktails.csv")
	viper.SetDefault("database.csv.data_dir", "./data")
}

// newConfig creates a new Config instance of type singleton.
// The default configuration instance can be configured only once.
func newConfig() Config {
	var cfg Config
	once.Do(func() {
		setDefaultConfig()
		viper.SetEnvPrefix("capstone")
		viper.AutomaticEnv()
		viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

		cfg = Config{
			Application: Application{
				version: viper.GetString("application.version"),
			},
			HTTP: HTTP{
				Server: HttpServer{
					host:            viper.GetString("http.server.host"),
					port:            viper.GetInt("http.server.port"),
					shutdownTimeout: viper.GetDuration("http.server.shutdown.timeout"),
				},
				DataAPI: DataAPI{
					url: viper.GetString("http.data_api.url"),
				},
			},
			Database: Database{
				driver: viper.GetString("database.driver"),
				Csv: CsvDB{
					fileName: viper.GetString("database.csv.file_name"),
					dataDir:  viper.GetString("database.csv.data_dir"),
				},
			},
		}
		logger.Log().Debug().
			Str("version", cfg.Application.Version()).
			Str("base_path", cfg.Application.BasePath()).
			Msg("configuration loaded")
	})
	return cfg
}
