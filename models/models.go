package models


// конфиги гонки, полученные из json файла
type ConfigInfo struct{
	LapsCount int `json:"laps"`  // количество кругов
	LapLen int `json:"lapLen"`  // длина круга
	PenaltyLapLen int `json:"penaltyLen"`  // длина штрафного круга
	FiringLinesCount int `json:"firingLines"`  // количество огневых рубежей на каждом круге (на одном рубеже 5 мишеней)
	StartTimeStr string `json:"start"`  // время старта
	StartDeltaStr string `json:"startDelta"`  // разница с которой надо стартовать
}

// информация об участнике, получаемая в процессе обработки эвентов
type CompetitorInfo struct{  // информация о конкретном участнике
	NotStartedFlag bool  // флаг о том, что участник не стартовал
	NoFinishedFlag bool  // флаг о том, что участник не финишировал
	ScheduledTimeStartStr string  // время запланированного страта
	ActualTimeStartStr string  // время старта на самом деле
	EveryLapTimes map[int][]string  // номер круга: информация о прохождении этого круга (старт, финиш)
	EveryPenaltyLapTimes map[int][]string  //  номер основного круга: информация о прохождении штрафного круга (старт, финиш)
	CounterHitTargets  int // счётчик попаданий по мишеням
}

// итоговая информация об участнике гонки (часть полей получена после преобразования полей CompetitorInfo)
type CompetitorResultInfo struct{
	CompetitorId string  // id участника
	DNSFInfo string  // финишировал/не стартовал/не финишировал
	TotalTime float64  // суммарное время на забег в секундах
	TotalTimeStr string  // суммарное время, затраченное на забег
	EachLapInfo string  // строка с информацией о времени и средней скорости каждого круга
	PenaltyLapsInfo string  // строка с информацией о штрафных минутах
	ShotsInfo string  // строка с информацией о точности стрельбы
}
