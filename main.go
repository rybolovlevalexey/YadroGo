package main

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

func parseHHMMSS(s string) (time.Duration, error) {
	t, err := time.Parse("15:04:05", s)
	if err != nil {
		return 0, err
	}
	h := t.Hour()
	m := t.Minute()
	s2 := t.Second()
	return time.Duration(h)*time.Hour + time.Duration(m)*time.Minute + time.Duration(s2)*time.Second, nil
}


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
	curTime := strings.TrimSpace(lineDataSplited[0])  // в строке сначала идёт время
	runesCurTime := []rune(curTime)
	curTimeCleaned := string(runesCurTime[1 : len(runesCurTime)-1])  // время без кавычек
	eventId := strings.TrimSpace(lineDataSplited[1])  // id события
	compId := strings.TrimSpace(lineDataSplited[2])  // id участника соревнований
	extraParam := ""  // дополнительной информации может не быть
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
	
	saveInfoFromLine(eventIdInt, compId, curTimeCleaned, extraParam)
	printOutputLog(eventIdInt, curTime, compId, extraParam)
}


// функция для печати информации о полученном событии
func printOutputLog(eventIdInt int, curTime string, compId string, extraParam string){
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


// функция для сохранения информации о полученном событии
func saveInfoFromLine(eventIdInt int, compId string, curTimeCleaned string, extraParam string){
	switch eventIdInt{
	case 1:  // регистрация
		// инициализация словарей со временем кругов
		competitorsInfo[compId].EveryLapTimes = make(map[int][]string)
		competitorsInfo[compId].EveryPenaltyLapTimes = make(map[int][]string)
		competitorsInfo[compId].EveryLapTimes[1] = []string{"", ""}
	case 2:  // получено своё время старта
		competitorsInfo[compId].ScheduledTimeStartStr = curTimeCleaned
		competitorsInfo[compId].EveryLapTimes[1][0] = extraParam
	case 3:  // участник на старте
		// то, что участник на страте, никак не влияет на сохраняемую информацию
	case 4:  // стартует (проверить, что он уложился в свой промежуток)
		competitorsInfo[compId].ActualTimeStartStr = curTimeCleaned

		// проверка, что старт был в положенный промежуток времени
		actualTime, _ := time.Parse("15:04:05.000", competitorsInfo[compId].ActualTimeStartStr)  // время старта
		startTime, _ := time.Parse("15:04:05", configInfo.StartTimeStr)  // время старта по расписанию
		delta, _ := parseHHMMSS(configInfo.StartDeltaStr)  // временной промежуток в который можно стартовать
		endTime := startTime.Add(delta)  // время после которого нельзя стартовать
		// log.Println(startTime, "\n", endTime, "\n", delta, "\n", actualTime)
		if !(actualTime.After(startTime) && actualTime.Before(endTime)){
			competitorsInfo[compId].NotStartedFlag = true
		} else {
			competitorsInfo[compId].NotStartedFlag = false
		}
	case 5:  // на огневом рубеже
		// время стрельбы входит во время круга поэтому данное событие ни на что не влияет
	case 6:  // попал в мишень
		competitorsInfo[compId].CounterHitTargets += 1  // увеличиваю счётчик попаданий
	case 7:  // покинул огневой рубеж
		// время стрельбы входит во время круга поэтому данное событие ни на что не влияет
	case 8:  // начал штрафные круги
		competitorsInfo[compId].EveryPenaltyLapTimes[len(competitorsInfo[compId].EveryLapTimes)] = []string{curTimeCleaned, ""}
	case 9:  // закончил штрафные круги
		competitorsInfo[compId].EveryPenaltyLapTimes[len(competitorsInfo[compId].EveryLapTimes)][1] = curTimeCleaned
	case 10:  // закончил круг и одновременно начал следующий, если это был не последний круг
		competitorsInfo[compId].EveryLapTimes[len(competitorsInfo[compId].EveryLapTimes)][1] = curTimeCleaned
		if len(competitorsInfo[compId].EveryLapTimes) < configInfo.LapsCount{  // информации о кругах меньше чем должно быть всего кругов
			competitorsInfo[compId].EveryLapTimes[len(competitorsInfo[compId].EveryLapTimes) + 1] = []string{curTimeCleaned, ""}
		}
	case 11:  // не может финишировать
		competitorsInfo[compId].NoFinishedFlag = true
	}
}


// печать итогового отчёта
func printFinalReport(){
	for key, value := range competitorsInfo{
		log.Println(key, value)
	}
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
