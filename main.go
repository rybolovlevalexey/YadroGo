package main

import (
	"YadroGo/usecases"
	"YadroGo/services"
	"YadroGo/settings"

	"log"
)


func main(){
	// логгер не пишет текущую дату для более аккуратного вывода
	log.SetFlags(0)
	
	// получение конфигов конктреной гонки
	configInfo, err := services.GetConfigInfo(settings.ConfigPath)
	if err != nil{
		log.Println("Не получены корректные конфиги")
		return 
	}
	log.Println(configInfo)
	
	// инициализация конкретной гонки со своими конфигами и запуск обработки эвентов
	biathlon := usecases.Biathlon{}
	biathlon.Init(configInfo)
	biathlon.StartProcessing(settings.EventsPath)

	// инициализация, создание и печать финального отчёта
	finalReport := usecases.FinalReport{}
	finalReport.CreateFinalReport(&biathlon)
	finalReport.PrintSortedFinalReport()
}
