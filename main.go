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
	EveryLapTimes map[int][]string  // номер круга: информация о прохождении этого круга (старт, финиш)
	EveryPenaltyLapTimes map[int][]string  //  номер основного круга: информация о прохождении штрафного круга (старт, финиш)
	CounterHitTargets  int // счётчик попаданий по мишеням
}

var configInfo ConfigInfo
var competitorsInfo map[string]*CompetitorInfo = make(map[string]*CompetitorInfo)
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

	printFinalReport()
}


// получение информации из конкретной строки IncomingEvents
func getInfoFromCurrentLine(lineData string){
	// обработка полученной строки
	lineDataSplited := strings.Split(lineData, " ")
	curTime := strings.TrimSpace(lineDataSplited[0])
	runesCurTime := []rune(curTime)
	curTimeCleaned := string(runesCurTime[1 : len(runesCurTime)-1])
	eventId := strings.TrimSpace(lineDataSplited[1])
	compId := strings.TrimSpace(lineDataSplited[2])
	extraParam := ""
	if len(lineDataSplited) > 3{
		extraParam = strings.TrimSpace(strings.Join(lineDataSplited[3:], " "))
	}
	
	// создание нового пользователя, если его не существует в общем словаре
	if _, exists := competitorsInfo[compId]; !exists{
		competitorsInfo[compId] = &CompetitorInfo{
			NotStartedFlag: true,
			NoFinishedFlag: true,
		}
	}

	// обработка событий в зависимости от их id
	eventIdInt, _ := strconv.Atoi(eventId)
	switch eventIdInt{
	case 1:
		// инициализация словарей со временем кругов
		competitorsInfo[compId].EveryLapTimes = make(map[int][]string)
		competitorsInfo[compId].EveryPenaltyLapTimes = make(map[int][]string)
		competitorsInfo[compId].EveryLapTimes[1] = []string{"", ""}

		log.Printf(curTime + " The competitor(%s) registered\n", compId)
	case 2:
		competitorsInfo[compId].ScheduledTimeStartStr = curTimeCleaned
		competitorsInfo[compId].EveryLapTimes[1][0] = extraParam

		log.Printf(curTime + " The start time for the competitor(%s) was set by a draw to (%s)\n", compId, extraParam)
	case 3:
		log.Printf(curTime + " The competitor(%s) is on the start line\n", compId)
	case 4:
		competitorsInfo[compId].ActualTimeStartStr = curTimeCleaned
		competitorsInfo[compId].NotStartedFlag = false

		log.Printf(curTime + " The competitor(%s) has started\n", compId)
	case 5:
		// время стрельбы входит во время круга поэтому данное событие ни на что не влияет
		log.Printf(curTime + " The competitor(%s) is on the firing range(%s)\n", compId, extraParam)
	case 6:
		competitorsInfo[compId].CounterHitTargets += 1  // увеличиваю счётчик попаданий

		log.Printf(curTime + " The target(%s) has been hit by competitor(%s)\n", extraParam, compId)
	case 7:
		// время стрельбы входит во время круга поэтому данное событие ни на что не влияет
		log.Printf(curTime + " The competitor(%s) left the firing range\n", compId)
	case 8:
		competitorsInfo[compId].EveryPenaltyLapTimes[len(competitorsInfo[compId].EveryLapTimes)] = []string{curTimeCleaned, ""}

		log.Printf(curTime + " The competitor(%s) entered the penalty laps\n", compId)
	case 9:
		competitorsInfo[compId].EveryPenaltyLapTimes[len(competitorsInfo[compId].EveryLapTimes)][1] = curTimeCleaned

		log.Printf(curTime + " The competitor(%s) left the penalty laps\n", compId)
	case 10:
		competitorsInfo[compId].EveryLapTimes[len(competitorsInfo[compId].EveryLapTimes)][1] = curTimeCleaned
		competitorsInfo[compId].EveryLapTimes[len(competitorsInfo[compId].EveryLapTimes) + 1] = []string{curTimeCleaned, ""}

		log.Printf(curTime + " The competitor(%s) ended the main lap\n", compId)
	case 11:
		competitorsInfo[compId].NoFinishedFlag = true
		log.Printf(curTime + " The competitor(%s) can`t continue: %s\n", compId, extraParam)
	}
}


// печать итогового отчёта
func printFinalReport(){
	/*
	for key, value := range competitorsInfo{
	
	}
	*/
}


// получение информации из конфигурационного файла
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
