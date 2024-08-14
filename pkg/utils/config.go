package utils

import (
	"encoding/json"
	"log"
	"os"
)

type Config struct {
	ChainId         int64  `json:"chainId"`
	RPCAddress      string `json:"RPCAddress"`
	Nickname        string `json:"nickname"`
	KeyStore        string `json:"keystore"`
	ContractAddress string `json:"contractAddress"`
	MasterKey       string `json:"masterKey"`
}

// func LoadConfiguration(file string) Config {
// 	var config Config
// 	configFile, err := os.Open(file)
// 	defer configFile.Close()
// 	if err != nil {
// 		fmt.Println(err.Error())
// 	}
// 	jsonParser := json.NewDecoder(configFile)
// 	jsonParser.Decode(&config)
// 	return config
// }

func LoadConfiguration(file string) (Config, error) {
	var config Config
	configFile, err := os.Open(file)
	if err != nil {
		log.Println(err)
		config = Config{
			ChainId:         1337,
			RPCAddress:      "http://127.0.0.1:8545",
			MasterKey:       "3ECB00DB9C0F56D72861E88A02D5D914629525EF03072B516A523FF92BB14F5D",
			Nickname:        "Guest",
			KeyStore:        "",
			ContractAddress: "",
		}
		return config, err
	} else {
		defer configFile.Close()
		jsonParser := json.NewDecoder(configFile)
		jsonParser.Decode(&config)
		return config, nil
	}
}
