package services

import (
	"time"
	"log"
	"os"
	"encoding/json"
	"io"

	"YadroGo/models"
)

func ParseHHMMSS(s string) (time.Duration, error) {
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		return 0, err
	}
	h := t.Hour()
	m := t.Minute()
	s2 := t.Second()
	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s2)*time.Second, nil
}


// получение информации из конфигурационного файла
func GetConfigInfo(configPath string) models.ConfigInfo{
	var configInfo models.ConfigInfo
	
	file, err := os.Open(configPath)
	if err != nil{
		log.Fatal("Problem in opening config file")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil{
		log.Fatal("Problem in reading config data")
	}

	err = json.Unmarshal(data, &configInfo)
	if err != nil{
		log.Fatal("Problem in json format of config file")
	}

	log.Println("config.json correctly parsed")	
	return configInfo
}
