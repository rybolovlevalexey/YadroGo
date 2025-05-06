package usecases

import (
	"YadroGo/models"
	"YadroGo/services"
	"YadroGo/settings"
	"reflect"

	"testing"
)

func TestBiathlonInit(t *testing.T){
	b := Biathlon{}
	settings.ConfigPath = "../files/configs/config.json"
	cnf, _ := services.GetConfigInfo(settings.ConfigPath)
	t.Log(cnf)
	b.Init(cnf)

	if b.configInfo.LapLen != 3500 || b.configInfo.StartDeltaStr != "00:01:30"{
		t.Errorf("not correct config info in biathlon class")
	}
}

func TestGetInfoFromCurrentLine(t *testing.T){
	b := Biathlon{}
	settings.ConfigPath = "../files/configs/config1.json"
	cnf, _ := services.GetConfigInfo(settings.ConfigPath)
	b.Init(cnf)

	b.getInfoFromCurrentLine("[09:05:59.867] 1 1")
	if len(b.competitorsInfo) == 0{
		t.Error("must be one competitor, got zero")
	}

	b.getInfoFromCurrentLine("[09:15:00.841] 2 1 09:30:00.000")
	if b.competitorsInfo["1"].ScheduledTimeStartStr != "09:30:00.000"{
		t.Error("competitor(1) must have scheduled time")
	}

	b.getInfoFromCurrentLine("[09:29:45.734] 3 1")
	if b.competitorsInfo["1"].CounterHitTargets != 0{
		t.Error("competitor 1 has not shooting")
	}

	b.getInfoFromCurrentLine("[09:30:01.005] 4 1")
	if b.competitorsInfo["1"].ActualTimeStartStr != "09:30:01.005"{
		t.Error("not correct actual time start")
	}

	for _, elem := range []string{
		"[09:49:31.659] 5 1 1", 
		"[09:49:33.123] 6 1 1", 
		"[09:49:34.650] 6 1 2", 
		"[09:49:35.937] 6 1 4", 
		"[09:49:37.364] 6 1 5",
		"[09:49:38.339] 7 1", 
	}{
		b.getInfoFromCurrentLine(elem)	
	}
	if b.competitorsInfo["1"].CounterHitTargets != 4{
		t.Error("not correct count of hited targets")
	}

	b.getInfoFromCurrentLine("[09:49:55.915] 8 1")
	b.getInfoFromCurrentLine("[09:51:48.391] 9 1")
	if b.competitorsInfo["1"].EveryPenaltyLapTimes[1][0] != "09:49:55.915" || 
			b.competitorsInfo["1"].EveryPenaltyLapTimes[1][1] != "09:51:48.391"{
		t.Error("no info or not correct info about penalty laps")
	}

	b.getInfoFromCurrentLine("[09:59:03.872] 10 1")
	t.Log(b.competitorsInfo["1"].EveryLapTimes)
	if b.competitorsInfo["1"].EveryLapTimes[1][0] != "09:30:00.000" || 
			b.competitorsInfo["1"].EveryLapTimes[1][1] != "09:59:03.872"{
			t.Error("not correct info about first lap timings")
	}
}


func TestCreateLapsInfo(t *testing.T){
	finalRep := FinalReport{}
	b := Biathlon{}
	settings.ConfigPath = "../files/configs/config1.json"
	settings.EventsPath = "../files/events/events1.txt"
	cnf, _ := services.GetConfigInfo(settings.ConfigPath)
	b.Init(cnf)
	b.StartProcessing(settings.EventsPath)
	
	curLapsDone := len(b.competitorsInfo["1"].EveryLapTimes)
	if b.competitorsInfo["1"].EveryLapTimes[len(b.competitorsInfo["1"].EveryLapTimes)][1] == ""{
		curLapsDone -= 1
	}
	resFinalRepLapInfo := finalRep.createEachLapInfo(b.competitorsInfo["1"].EveryLapTimes, 
		curLapsDone,
		b.configInfo.LapLen,
		b.configInfo.LapsCount,
	)
	if resFinalRepLapInfo != "[{00:29:03.872, 2.094}; {,}]"{
		t.Error("not correct string with info about every lap times and avg speed")
	}

	curPenaltyLapsDone := len(b.competitorsInfo["1"].EveryPenaltyLapTimes)
	if b.competitorsInfo["1"].EveryPenaltyLapTimes[len(b.competitorsInfo["1"].EveryPenaltyLapTimes)][1] == ""{
		curPenaltyLapsDone -= 1
	}
	resFinalRepPenaltyLapInfo := finalRep.createPenaltyLapsInfo(
		b.competitorsInfo["1"].EveryPenaltyLapTimes,
		curPenaltyLapsDone, 
		b.configInfo.PenaltyLapLen,
	)
	if resFinalRepPenaltyLapInfo != "{00:01:52.476, 0.445}"{
		t.Log("not correct string with info about penalty laps")
	}
}

func TestSortFinalReport(t *testing.T){
	b := Biathlon{}
	settings.ConfigPath = "../files/configs/config.json"
	settings.EventsPath = "../files/events/events.txt"
	cnf, _ := services.GetConfigInfo(settings.ConfigPath)
	b.Init(cnf)
	b.StartProcessing(settings.EventsPath)
	
	finalRep := FinalReport{}
	finalRep.sortFinalReport()

	compIdsAfterSort := []string{}
	for _, array := range [][]models.CompetitorResultInfo{finalRep.ResultMapFinished, finalRep.ResultMapDNSF}{
		for _, elem := range array{
			compIdsAfterSort = append(compIdsAfterSort, elem.CompetitorId)
		}
	}

	if !reflect.DeepEqual(compIdsAfterSort, []string{"2", "1", "3", "4", "5"}){
		t.Error("not correct order of competitors ids after sorting final report results")
	}
}