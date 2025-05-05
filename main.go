package main

import (
	"YadroGo/usecases"
	"YadroGo/services"

	"log"
)

const configPath = "files/configs/config.json"
const eventsPath = "files/events/events.txt"


func main(){
	// логгер не пишет текущую дату для более аккуратного вывода
	log.SetFlags(0)
	
	// получение конфигов конктреной гонки
	configInfo := services.GetConfigInfo(configPath)
	log.Println(configInfo)
	
	// инициализация конкретной гонки со своими конфигами и запуск обработки эвентов
	biathlon := usecases.Biathlon{}
	biathlon.Init(configInfo)
	biathlon.StartProcessing(eventsPath)

	// инициализация, создание и печать финального отчёта
	finalReport := usecases.FinalReport{}
	finalReport.CreateFinalReport(&biathlon)
	finalReport.PrintSortedFinalReport()
}
