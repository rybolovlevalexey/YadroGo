package usecases

import (
	"YadroGo/models"
	
	"strconv"
	"log"
	"time"
	"fmt"
	"sort"
)

// класс по созданию финального отчёта, хранит два списка: гонщиков финишировавших и всех остальных
type FinalReport struct{
	ResultMapFinished []models.CompetitorResultInfo
	ResultMapDNSF []models.CompetitorResultInfo
}

// генерация финального отчёта (на вход подаётся конкретная гонка):
// делит гонщиков на финишировавших и нет
// преобразует входные данные (вычисляет время, потраченное на каждый круг, подготавливает строки, которые будут выводиться)
func (f *FinalReport) CreateFinalReport(b *Biathlon){
	// массив с итоговой инфорацией участников, которые закончили гонку, чтобы их потом можно было отсортировать
	resultMapFinished := []models.CompetitorResultInfo{}
	// массив с итоговой инфорацией участников, которые не начали / не закончили гонку, чтобы их не сортировать, 
	// а вывести в конце списка в том порядке, в котором они встречаются
	resultMapDNSF := []models.CompetitorResultInfo{}

	for key, value := range b.competitorsInfo{
		// log.Println(key, value)
		// вычисление количества завершёённых кругов
		var curLapsDone int
		if value.EveryLapTimes[len(value.EveryLapTimes)][1] == ""{
			curLapsDone = len(value.EveryLapTimes) - 1
		} else {
			curLapsDone = len(value.EveryLapTimes)
		}

		curCompetitorResultInfo := models.CompetitorResultInfo{
			CompetitorId: key,  // id участника
			ShotsInfo: strconv.Itoa(value.CounterHitTargets) + "/" + strconv.Itoa(5 * b.configInfo.FiringLinesCount), // вычисление результатов стрельбы
			EachLapInfo: f.createEachLapInfo(value.EveryLapTimes, curLapsDone, 
				b.configInfo.LapLen, b.configInfo.LapsCount),  // информация по кругам
			PenaltyLapsInfo: f.createPenaltyLapsInfo(
				value.EveryPenaltyLapTimes, 
				5 * b.configInfo.FiringLinesCount - value.CounterHitTargets, 
				b.configInfo.PenaltyLapLen),
			TotalTime: -1,
			TotalTimeStr: "-1",
		}

		// наполнение остальной информации в зависимости от результатов
		if value.NoFinishedFlag || value.NotStartedFlag || curLapsDone != b.configInfo.LapsCount{
			if value.NotStartedFlag{
				curCompetitorResultInfo.DNSFInfo = "[NotStarted]"
			} else{
				curCompetitorResultInfo.DNSFInfo = "[NotFinished]"
			}

			resultMapDNSF = append(resultMapDNSF, curCompetitorResultInfo)
		} else {
			curCompetitorResultInfo.DNSFInfo = "[Finished]"
			
			// определение итогового времени, затраченного на всю гонку
			startTime, _ := time.Parse("15:04:05.000", value.ScheduledTimeStartStr)
			finishTime, _ := time.Parse("15:04:05.000", value.EveryLapTimes[curLapsDone][1])

			curCompetitorResultInfo.TotalTimeStr = time.Time{}.Add(finishTime.Sub(startTime)).Format("15:04:05.000")
			curCompetitorResultInfo.TotalTime = finishTime.Sub(startTime).Seconds()

			resultMapFinished = append(resultMapFinished, curCompetitorResultInfo)
		}
	}

	f.ResultMapFinished = resultMapFinished
	f.ResultMapDNSF = resultMapDNSF
}


// создание строки с информацией о каждом круге
// входные данные: словарь - номером круга (ключ) и список начала и конца круга (значение)
// количество кругов, которые завершил данный гонщик; длина круга; общее количество кругов
func (f *FinalReport) createEachLapInfo(everyLapTimes map[int][]string, curLapsDone int, lapLen int, lapsCount int) string{
	var resultString string

	resultString = "["
	for lapNum := 1; lapNum <= len(everyLapTimes); lapNum += 1{
		lapInfo, exists := everyLapTimes[lapNum]
		if !exists{
			continue
		}
		
		if lapInfo[0] == "" || lapInfo[1] == ""{
			continue
		}

		resultString += "{"
		startTime, _ := time.Parse("15:04:05.000", lapInfo[0])
		finishTime, _ := time.Parse("15:04:05.000", lapInfo[1])
		resultString += time.Time{}.Add(finishTime.Sub(startTime)).Format("15:04:05.000")
		resultString += ", "
		resultString += fmt.Sprintf("%.3f", float64(lapLen) / finishTime.Sub(startTime).Seconds())
		resultString += "}"
		if lapNum != curLapsDone{
			resultString += "; "
		}
	}
	if curLapsDone < lapsCount{
		resultString += "; {,}"
	}
	resultString += "]"
	// log.Println(resultString)
	return resultString
}


// создание строки с информацией о прохождении штрафных кругов
// входные данные: словарь - номером  круга (ключ) и список начала и конца прохождения штрафных кругов (значение)
// количество штрафных кругов, которые завершил данный гонщик; длина штрафного круга
func (f *FinalReport) createPenaltyLapsInfo(everyPenaltyLapTimes map[int][]string, penaltyLapsCount int, penaltyLapLen int) string {
    var totalPenaltyDuration time.Duration

    for _, value := range everyPenaltyLapTimes {
        startTime, err := time.Parse("15:04:05.000", value[0])
        if err != nil {
            log.Printf("Error parsing start time: %v", err)
            continue
        }

        finishTime, err := time.Parse("15:04:05.000", value[1])
        if err != nil {
            log.Printf("Error parsing finish time: %v", err)
            continue
        }

        totalPenaltyDuration += finishTime.Sub(startTime)
    }

    // Рассчитываем общую дистанцию (в метрах, если PenaltyLapLen в метрах)
    totalDistance := float64(penaltyLapsCount * penaltyLapLen)
    
    // Рассчитываем среднюю скорость (м/с)
    totalSeconds := totalPenaltyDuration.Seconds()
    var avgSpeed float64
    if totalSeconds > 0 {
        avgSpeed = totalDistance / totalSeconds
    }

    // Форматируем общее время в строку HH:MM:SS.fff
    hours := int(totalPenaltyDuration.Hours())
    minutes := int(totalPenaltyDuration.Minutes()) % 60
    seconds := int(totalPenaltyDuration.Seconds()) % 60
    milliseconds := int(totalPenaltyDuration.Milliseconds()) % 1000

    resultString := fmt.Sprintf("{%02d:%02d:%02d.%03d, %.3f}", 
        hours, minutes, seconds, milliseconds, avgSpeed)

	if resultString == "{00:00:00.000, 0.000}"{
		resultString = "{,}"
	}

    return resultString
}


func (f *FinalReport) sortFinalReport(){
	// сортировка финишировавших участников
	sort.Slice(f.ResultMapFinished, func(i int, j int) bool{
		return f.ResultMapFinished[i].TotalTime < f.ResultMapFinished[j].TotalTime
	})

	// сортировка остальных участников по id
	sort.Slice(f.ResultMapDNSF, func(i int, j int) bool{
		return f.ResultMapDNSF[i].CompetitorId < f.ResultMapDNSF[j].CompetitorId
	})
}


// вывод итогового отчёта(финишировавшие сортируются по возрастанию общего времени; 
// остальные печатаются в конце по возрастанию id)
func (f *FinalReport) PrintSortedFinalReport(){
	f.sortFinalReport()
	
	for _, elem := range f.ResultMapFinished{
		log.Printf(
			"%s %s(%s) %s %s %s\n", elem.DNSFInfo, elem.CompetitorId, 
			elem.TotalTimeStr, elem.EachLapInfo, elem.PenaltyLapsInfo, elem.ShotsInfo,
		)
	}

	// печать не финишировавших и не стратовавших участников без сортировки
	for _, elem := range f.ResultMapDNSF{
		log.Printf(
			"%s %s %s %s %s\n", elem.DNSFInfo, elem.CompetitorId, 
			elem.EachLapInfo, elem.PenaltyLapsInfo, elem.ShotsInfo,
		)
	}
}