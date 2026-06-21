package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	toml "github.com/pelletier/go-toml"
)

// Config specifies the available configuration values
type Config struct {
	Title      string
	Root       string
	Debug      bool
	SessionKey string
	Database   struct {
		Driver   string
		Username string
		Password string
		Name     string
		Host     string
		Port     int
		URI      string
	}
	Session struct {
		Key string
	}
}

// Conf is the loaded configuration
var Conf Config

func init() {
	var err error
	Conf, err = LoadConfig("config.toml")

	if err != nil {
		Conf, _ = LoadConfig("../../config.toml")
	}
}

// LoadConfig loads a config from a file
func LoadConfig(f string) (Config, error) {
	stream, err := ioutil.ReadFile(f)

	if err != nil {
		// Will return an empty config populated by environment variables
		return SetupConfig(make([]byte, 0)), err
	}

	return SetupConfig(stream), nil
}

// SetupConfig converts a byte stream to a Config object
func SetupConfig(stream []byte) Config {
	config := Config{}
	toml.Unmarshal(stream, &config)

	if "" == config.Title {
		config.Title = "Argument"
	}

	if "" != os.Getenv("DB_DRIVER") {
		config.Database.Driver = os.Getenv("DB_DRIVER")
	}

	if "sqlite" == config.Database.Driver {
		config.Database.URI = "argument.db"

		return config
	}

	if "" != os.Getenv("DB_HOST") {
		config.Database.Host = os.Getenv("DB_HOST")
	}

	if "" != os.Getenv("DB_PORT") {
		i, _ := strconv.Atoi(os.Getenv("DB_PORT"))
		config.Database.Port = i
	}

	config.Database.URI = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", config.Database.Username, config.Database.Password, config.Database.Host, config.Database.Port, config.Database.Name)

	return config
}
