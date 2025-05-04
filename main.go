package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
)

type ConfigInfo struct{
	LapsCount int `json:"laps"`
	LapLen int `json:"lapLen"`
	PenaltyLapLen int `json:"penaltyLen"`
	FiringLinesCount int `json:"firingLines"`
	StartTimeStr string `json:"start"`
	StartDeltaStr string `json:"startDelta"`
}

func main(){
	configInfo := getConfigInfo("config.json")
	log.Println(configInfo)
}


func getConfigInfo(configPath string) ConfigInfo{
	file, err := os.Open(configPath)
	if err != nil{
		log.Fatal("Problem in opening file")
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil{
		log.Fatal("Problem in reading data")
	}

	var configInfo ConfigInfo
	err = json.Unmarshal(data, &configInfo)
	if err != nil{
		log.Fatal("Problem in json format")
	}

	log.Println("config.json correctly parsed")
	return configInfo
}