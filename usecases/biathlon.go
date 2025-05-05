package usecases

import (
	"YadroGo/models"
	"YadroGo/core"
	
	"strings"
	"strconv"
	"log"
	"time"
	"os"
	"io"
)




type Biathlon struct{
	configInfo models.ConfigInfo
	competitorsInfo map[string]*models.CompetitorInfo
}

func (b *Biathlon) Init(configInfo models.ConfigInfo){
	b.configInfo = configInfo
	b.competitorsInfo = make(map[string]*models.CompetitorInfo)
}

func (b *Biathlon) StartProcessing(eventsPath string){
	inputFile, err := os.Open(eventsPath)
	if err != nil{
		log.Fatal("Problem in opening input file")
	}
	inputData, err := io.ReadAll(inputFile)
	if err != nil{
		log.Fatal("Problem in reading input data")
	}
	inputDataStr := string(inputData)

	for _, line := range strings.Split(inputDataStr, "\n"){
		if line == ""{
			continue
		}
		b.getInfoFromCurrentLine(line)
	}
}

// получение информации из конкретной строки IncomingEvents
func (b *Biathlon) getInfoFromCurrentLine(lineData string){
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
	if _, exists := b.competitorsInfo[compId]; !exists{
		b.competitorsInfo[compId] = &models.CompetitorInfo{
			NotStartedFlag: true,
			NoFinishedFlag: true,
		}
	}

	// обработка событий в зависимости от их id
	eventIdInt, _ := strconv.Atoi(eventId)
	
	b.saveInfoFromLine(eventIdInt, compId, curTimeCleaned, extraParam)
	b.printOutputLog(eventIdInt, curTime, compId, extraParam)
}


// функция для печати информации о полученном событии
func (b *Biathlon) printOutputLog(eventIdInt int, curTime string, compId string, extraParam string){
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
func (b *Biathlon) saveInfoFromLine(eventIdInt int, compId string, curTimeCleaned string, extraParam string){
	switch eventIdInt{
	case 1:  // регистрация
		// инициализация словарей со временем кругов
		b.competitorsInfo[compId].EveryLapTimes = make(map[int][]string)
		b.competitorsInfo[compId].EveryPenaltyLapTimes = make(map[int][]string)
		b.competitorsInfo[compId].EveryLapTimes[1] = []string{"", ""}
	case 2:  // получено своё время старта
		b.competitorsInfo[compId].ScheduledTimeStartStr = extraParam
		b.competitorsInfo[compId].EveryLapTimes[1][0] = extraParam
	case 3:  // участник на старте
		// то, что участник на страте, никак не влияет на сохраняемую информацию
	case 4:  // стартует (проверить, что он уложился в свой промежуток)
		b.competitorsInfo[compId].ActualTimeStartStr = curTimeCleaned

		// проверка, что старт был в положенный промежуток времени
		actualTime, _ := time.Parse("15:04:05.000", b.competitorsInfo[compId].ActualTimeStartStr)  // время старта
		startTime, _ := time.Parse("15:04:05", b.competitorsInfo[compId].ScheduledTimeStartStr)  // время старта по расписанию
		delta, _ := core.ParseHHMMSS(b.configInfo.StartDeltaStr)  // временной промежуток в который можно стартовать
		endTime := startTime.Add(delta)  // время после которого нельзя стартовать
		// log.Println(startTime, "\n", endTime, "\n", delta, "\n", actualTime)
		if !(actualTime.After(startTime) && actualTime.Before(endTime)){
			b.competitorsInfo[compId].NotStartedFlag = true
		} else {
			b.competitorsInfo[compId].NotStartedFlag = false
		}
	case 5:  // на огневом рубеже
		// время стрельбы входит во время круга поэтому данное событие ни на что не влияет
	case 6:  // попал в мишень
		b.competitorsInfo[compId].CounterHitTargets += 1  // увеличиваю счётчик попаданий
	case 7:  // покинул огневой рубеж
		// время стрельбы входит во время круга поэтому данное событие ни на что не влияет
	case 8:  // начал штрафные круги
		b.competitorsInfo[compId].EveryPenaltyLapTimes[len(b.competitorsInfo[compId].EveryLapTimes)] = []string{curTimeCleaned, ""}
	case 9:  // закончил штрафные круги
		b.competitorsInfo[compId].EveryPenaltyLapTimes[len(b.competitorsInfo[compId].EveryLapTimes)][1] = curTimeCleaned
	case 10:  // закончил круг и одновременно начал следующий, если это был не последний круг
		b.competitorsInfo[compId].EveryLapTimes[len(b.competitorsInfo[compId].EveryLapTimes)][1] = curTimeCleaned
		if len(b.competitorsInfo[compId].EveryLapTimes) < b.configInfo.LapsCount{  // информации о кругах меньше чем должно быть всего кругов
			b.competitorsInfo[compId].EveryLapTimes[len(b.competitorsInfo[compId].EveryLapTimes) + 1] = []string{curTimeCleaned, ""}
		} else {
			b.competitorsInfo[compId].NoFinishedFlag = false
		}
	case 11:  // не может финишировать
		b.competitorsInfo[compId].NoFinishedFlag = true
	}
}

