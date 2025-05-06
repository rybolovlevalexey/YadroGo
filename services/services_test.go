package services

import (
	"YadroGo/settings"
	"time"

	"fmt"
	"strings"
	"testing"
)

func TestGetConfigInfo(t *testing.T){
	settings.ConfigPath = "../" + settings.ConfigPath
	
	// тестирование некорректных файлов
	for _, elem := range []string{
		"../nofile.json",
		"../files/configs/nofile.json", 
		"../files/configs/testconfig.json", 
		"../files/configs/testbadformatconfig.json",
	}{
		_, err := GetConfigInfo(elem)
		if err == nil{
			t.Error("excepted error in opening/reading/unmarshal file")
		}
	}

	// тестирование корректных файлов
	for i := 0; i <= 2; i += 1{
		var cnfPath string
		if i == 0{
			cnfPath = settings.ConfigPath
		} else {
			splittedCnfPath := strings.Split(settings.ConfigPath, ".")
			cnfPath = strings.Join(splittedCnfPath[0:len(splittedCnfPath) - 1], ".") + fmt.Sprintf(
				"%d.", i) + splittedCnfPath[len(splittedCnfPath) - 1]
		}
		cnf, err := GetConfigInfo(cnfPath)
		if err != nil{
			t.Errorf("excepted no errors")	
		}
		t.Logf("cnfPath - %s", cnfPath)
		
		switch i{
		case 0:
			if cnf.LapLen != 3500{
				t.Errorf("Expected 3500 in ConfigInfo.LapLen")
			}
			if cnf.StartTimeStr != "10:00:00.000"{
				t.Errorf("Expected '10:00:00.000' in ConfigInfo.StartTimeStr")
			}
			t.Logf("tests with cnfPath - %s done correctly", cnfPath)
		case 1:
			if cnf.LapsCount != 2{
				t.Errorf("Expected 2 in ConfigInfo.LapsCount")
			}
			if cnf.StartDeltaStr != "00:00:30"{
				t.Errorf("Expected '00:00:30'")
			}
			t.Logf("tests with cnfPath - %s done correctly", cnfPath)
		case 2:
			if cnf.PenaltyLapLen != 50{
				t.Errorf("Expected 50 in ConfigInfo.PenaltyLapLen")
			}
			if cnf.FiringLinesCount != 1{
				t.Errorf("Expected 1 in ConfigInfo.FiringLinesCount")
			}
		}
	}
}


func TestParseHHMMSS(t *testing.T){
	timeDuration, err := ParseHHMMSS("00:01:30")
	if err != nil{
		t.Errorf("err expected not nil")
	}

	if timeDuration.Minutes() != 1.5{
		t.Errorf("minutes in timeDuration expected 1.5")
	}

	parsedTime, _ := time.Parse("15:04", "05:30")
	res := parsedTime.Add(timeDuration)
	t.Log(res)
	if res.Second() != 30{
		t.Errorf("expected 30 in seconds after add duration(00:01:30)) to 5:30:00, got %d", res.Second())
	}
	if res.Minute() != 31{
		t.Errorf("expected 31 in minute after add duration(00:01:30)) to 5:30:00, got %d", res.Minute())
	}


	timeDuration, _ = ParseHHMMSS("01:21:15")
	if timeDuration.Minutes() != 81.25{
		t.Errorf("minutes in timeDuration expected 81.25")
	}
}