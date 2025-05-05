package main

import (
	"YadroGo/usecases"
	"YadroGo/core"

	"log"
)

const configPath = "files/configs/config.json"
const eventsPath = "files/events/events.txt"


func main(){
	log.SetFlags(0)
	
	configInfo := core.GetConfigInfo(configPath)
	log.Println(configInfo)
	
	biathlon := usecases.Biathlon{}
	biathlon.Init(configInfo)
	biathlon.StartProcessing(eventsPath)

	finalReport := usecases.FinalReport{}
	finalReport.CreateFinalReport(&biathlon)
	finalReport.PrintSortedFinalReport()
}
