package main

import (
	"log"

	"deforestation.detection.com/server/internal/app/apiserver"
)

//var (
//	configPath string
//)

//func init() {
//	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "path to config file")
//}

func main() {
	//flag.Parse()

	config := apiserver.NewConfig()
	log.Printf("Server Configs are %v", config)

	//_, err := toml.DecodeFile(configPath, config)
	//if err != nil {
	//	log.Fatal(err)
	//}

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
