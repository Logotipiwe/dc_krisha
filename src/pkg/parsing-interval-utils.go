package pkg

import (
	. "github.com/logotipiwe/dc_go_config_lib"
	"krisha/src/internal/domain"
	"strconv"
	"strings"
	"time"
)

func GetActiveParsingIntervalMs() int {
	defaultVal := 120 * 1000
	got := GetConfigIntOr("ACTIVE_PARSE_INTERVAL_MS", defaultVal)
	if got <= 0 {
		return defaultVal
	}
	return got
}

func GetNonActiveParsingIntervalMs() int {
	defaultVal := 600 * 1000
	got := GetConfigIntOr("NON_ACTIVE_PARSE_INTERVAL_MS", defaultVal)
	if got <= 0 {
		return defaultVal
	}
	return got
}

type ActiveParsingPeriod struct {
	StartHour int
	EndHour   int
}

func GetActiveParsingHours() *ActiveParsingPeriod {
	got := GetConfig("ACTIVE_PARSE_HOURS")
	split := strings.Split(got, "-")
	if len(split) != 2 {
		return nil
	}
	startHour, err := strconv.Atoi(split[0])
	if err != nil {
		return nil
	}
	endHour, err := strconv.Atoi(split[1])
	if err != nil {
		return nil
	}
	if startHour < 0 || endHour < 0 {
		return nil
	}
	return &ActiveParsingPeriod{
		StartHour: startHour,
		EndHour:   endHour,
	}
}

func GetParserSleepingInterval(settings *domain.ParserSettings) time.Duration {
	activeHours := GetActiveParsingHours()
	activeInterval := GetActiveParsingIntervalMs()
	nonActiveInterval := GetNonActiveParsingIntervalMs()

	if activeHours == nil {
		return time.Duration(activeInterval) * time.Millisecond
	} else {
		currHour := time.Now().Hour()
		//active - 8-23; 8-0; 8-1
		//if start < end - (curr >= start && curr < end)
		//if start > end - (curr > start || curr < end)
		var isInActiveHours bool
		if activeHours.StartHour < activeHours.EndHour {
			//example "8-23" - from 8am to 11pm
			isInActiveHours = currHour >= activeHours.StartHour && currHour < activeHours.EndHour
		} else {
			//example "8-1" - from 8am to 1am
			isInActiveHours = currHour >= activeHours.StartHour || currHour < activeHours.EndHour
		}
		if isInActiveHours {
			return time.Duration(activeInterval) * time.Millisecond
		} else {
			return time.Duration(nonActiveInterval) * time.Millisecond
		}
	}
}
