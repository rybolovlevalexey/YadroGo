package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

type ConfigInfo struct{
	LapsCount int `json:"laps"`  // количество кругов
	LapLen int `json:"lapLen"`  // длина круга
	PenaltyLapLen int `json:"penaltyLen"`  // длина штрафного круга
	FiringLinesCount int `json:"firingLines"`  // количество огневых рубежей на каждом круге (на одном рубеже 5 мишеней)
	StartTimeStr string `json:"start"`  // время старта
	StartDeltaStr string `json:"startDelta"`  // разница с которой надо стартовать
}

type CompetitorInfo struct{  // информация о конкретном участнике
	NotStartedFlag bool  // флаг о том, что участник не стартовал
	NoFinishedFlag bool  // флаг о том, что участник не финишировал
	ScheduledTimeStartStr string  // время запланированного страта
	ActualTimeStartStr string  // время старта на самом деле
	EveryLapTimes map[int]int  // номер круга: время прохождения
	EveryPenaltyLapTimes map[int]int  //  номер штрафного круга: время прохождения
	CounterHitTargets  int // счётчик попаданий по мишеням
}

var configInfo ConfigInfo
var competitorsInfo map[string]CompetitorInfo = make(map[string]CompetitorInfo)
var timeLayout = "10:00:00.000"
var timeLayoutConfig = "15:04:05"

func main(){
	log.SetFlags(0)
	getConfigInfo("config.json")
	log.Println(configInfo)
	
	inputFile, err := os.Open("events/test_events.txt")
	if err != nil{
		log.Fatal("Problem in opening input file")
	}
	inputData, err := io.ReadAll(inputFile)
	if err != nil{
		log.Fatal("Problem in reading input data")
	}
	inputDataStr := string(inputData)

	for _, line := range strings.Split(inputDataStr, "\n"){
		getInfoFromCurrentLine(line)
	}
}


func getInfoFromCurrentLine(lineData string){
	// обработка полученной строки
	lineDataSplited := strings.Split(lineData, " ")
	curTime := strings.TrimSpace(lineDataSplited[0])
	eventId := strings.TrimSpace(lineDataSplited[1])
	compId := strings.TrimSpace(lineDataSplited[2])
	extraParam := ""
	if len(lineDataSplited) > 3{
		extraParam = strings.TrimSpace(strings.Join(lineDataSplited[3:], " "))
	}
	
	// создание нового пользователя, если его не существует в общем словаре
	if _, exists := competitorsInfo[compId]; !exists{
		competitorsInfo[compId] = CompetitorInfo{
			NotStartedFlag: true,
			NoFinishedFlag: true,
		}
	}

	// обработка событий в зависимости от их id
	eventIdInt, _ := strconv.Atoi(eventId)
	switch eventIdInt{
	case 1:
		log.Printf(curTime + " The competitor(%s) registered\n", compId)
	case 2:
		log.Printf(curTime + " The start time for the competitor(%s) was set by a draw to (%s)\n", compId, extraParam)
	case 3:
		log.Printf(curTime + " The competitor(%s) is on the start line\n", compId)
	case 4:
		log.Printf(curTime + " The competitor(%s) has started\n", compId)
	case 5:
		log.Printf(curTime + " The competitor(%s) is on the firing range(%s)\n", compId, extraParam)
	case 6:
		log.Printf(curTime + " The target(%s) has been hit by competitor(%s)\n", extraParam, compId)
	case 7:
		log.Printf(curTime + " The competitor(%s) left the firing range\n", compId)
	case 8:
		log.Printf(curTime + " The competitor(%s) entered the penalty laps\n", compId)
	case 9:
		log.Printf(curTime + " The competitor(%s) left the penalty laps\n", compId)
	case 10:
		log.Printf(curTime + " The competitor(%s) ended the main lap\n", compId)
	case 11:
		log.Printf(curTime + " The competitor(%s) can`t continue: %s\n", compId, extraParam)
	}
}


func getConfigInfo(configPath string) ConfigInfo{
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
