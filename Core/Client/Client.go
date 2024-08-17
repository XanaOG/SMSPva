package Client

import (
	"encoding/json"
	"fmt"
	"os"
)

var (
	ConfigFile  = "./Assets/Json/Config.json"
	CountryFile = "./Assets/Json/List.json"
)

type Option struct {
	Name   string `json:"Name"`
	Number int    `json:"Number"`
}
type Config struct {
	APIKey string `json:"APIKey"`
	List   struct {
		Countries []string `json:"Countries"`
		Options   []Option `json:"Options"`
	} `json:"List"`
}

func GetConfig(file string) Config {
	var config Config
	ConfigFile, err := os.Open(file)
	if err != nil {
		fmt.Println("Config file not found." + err.Error())
	}
	defer ConfigFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(ConfigFile)
	jsonParser.Decode(&config)
	return config
}
