package configs

import (
	"errors"
	"fmt"
	"sync"

	"github.com/joho/godotenv"

	"github.com/spf13/viper"
)

var (
	instance *Configuration
	once     sync.Once
)

// Config ...
func Config() *Configuration {
	fmt.Println("Config")

	once.Do(func() {
		instance = load()
	})

	return instance
}

// TestConfig ...
func TestConfig() *Configuration {
	fmt.Println("Config")

	return &Configuration{
		PostgresHost:     "localhost",
		PostgresPort:     5432,
		PostgresUser:     "postgres",
		PostgresPassword: "12345",
		PostgresDatabase: "review_db",
		ServerHost:       "localhost",
		RPCPort:          ":8090",
	}
}

// Configuration ...
type Configuration struct {
	LogLevel    string `json:"log_level"`
	Environment string `json:"environment"`

	PostgresHost     string
	PostgresPort     int
	PostgresDatabase string
	PostgresUser     string
	PostgresPassword string
	ServerHost       string

	RPCPort string

	// context timeout in seconds
	ServerReadTimeout int
	ServiceName       string

	ApexAPIUser     string
	ApexAPIPassword string
	Apex1CURL       string
}

func load() *Configuration {

	// load .env file from given path
	err := godotenv.Load("./src/review_service/.env")
	if err != nil {
		panic(err)
	}

	var config Configuration

	v := viper.New()
	v.AutomaticEnv()

	config.Environment = v.GetString("ENVIRONMENT")
	config.LogLevel = v.GetString("LOG_LEVEL")
	config.PostgresDatabase = v.GetString("POSTGRES_DB")
	config.PostgresUser = v.GetString("POSTGRES_USER")
	config.PostgresPassword = v.GetString("POSTGRES_PASSWORD")
	config.PostgresHost = v.GetString("POSTGRES_HOST")
	config.PostgresPort = v.GetInt("POSTGRES_PORT")
	config.RPCPort = v.GetString("RPC_PORT")
	config.ServiceName = v.GetString("SERVICE_NAME")

	config.ApexAPIUser = v.GetString("APEX_API_USER")
	config.ApexAPIPassword = v.GetString("APEX_API_PASSWORD")
	config.Apex1CURL = v.GetString("APEX_1C_URL")

	return &config
}

// Validate validates the configuration
func (c *Configuration) Validate() error {
	if c.RPCPort == "" {
		return errors.New("rpc_port required")
	}
	return nil
}
