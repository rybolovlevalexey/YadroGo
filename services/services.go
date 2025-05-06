package services

import (
	"time"
	"log"
	"os"
	"encoding/json"
	"io"
	"errors"

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
func GetConfigInfo(configPath string) (models.ConfigInfo, error){
	var configInfo models.ConfigInfo
	
	file, err := os.Open(configPath)
	if err != nil{
		log.Printf("Problem in opening config file %s\n", configPath)
		return configInfo, errors.New("opening file problem")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil{
		log.Printf("Problem in reading config data of file %s\n", configPath)
		return configInfo, errors.New("problem in file data")
	}

	err = json.Unmarshal(data, &configInfo)
	if err != nil{
		log.Printf("Problem in json format of file %s\n", configPath)
		return configInfo, errors.New("json format")
	}

	log.Printf("%s correctly parsed\n", configPath)	
	return configInfo, nil
}
