package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"
)

type HolidayType struct {
	name     string
	isOffDay bool // true if vacation else work
}

type AnnualHolidays struct {
	year     int
	holidays map[string]HolidayType
}

func FetchAnnualHolidayInfo(year int) (AnnualHolidays, error) {
	holidayApiUrl := "https://holiday.cyi.me/api/holidays?year=" + strconv.Itoa(year)
	resp, err := http.Get(holidayApiUrl)
	if err != nil {
		return AnnualHolidays{}, err
	}
	defer resp.Body.Close()

	holidayInfo, err := io.ReadAll(resp.Body)
	if err != nil {
		return AnnualHolidays{}, err
	}

	return ParseAnnualHolidayInfo(holidayInfo)
}

func ParseAnnualHolidayInfo(annualHolidayInfo []byte) (AnnualHolidays, error) {
	var holiday map[string]interface{}

	err := json.Unmarshal(annualHolidayInfo, &holiday)
	if err != nil {
		log.Printf("解析 Json 失败")
		return AnnualHolidays{}, err
	}

	var annualHolidays AnnualHolidays
	annualHolidays.year, err = strconv.Atoi(holiday["year"].(string))
	if err != nil {
		log.Printf("json 中年份不存在, 或不为整数")
		return AnnualHolidays{}, err
	}

	annualHolidays.holidays = make(map[string]HolidayType)
	for _, holiday := range holiday["days"].([]interface{}) {
		if holidayMap, ok := holiday.(map[string]interface{}); ok {
			annualHolidays.holidays[holidayMap["date"].(string)] = HolidayType{holidayMap["name"].(string), holidayMap["isOffDay"].(bool)}
		} else {
			log.Printf("节假日格式错误, %s", holiday)
			return AnnualHolidays{}, err
		}
	}

	return annualHolidays, nil
}

func (annualHolidays *AnnualHolidays) Update() {
	if annualHolidays.year != time.Now().Year() {
		newHolidays, err := FetchAnnualHolidayInfo(time.Now().Year())
		if err != nil {
			log.Fatal(err)
		}

		annualHolidays = &newHolidays
	}
}

func (annualHolidays AnnualHolidays) DetermineIfTodayShouldWork() bool {
	currTime := time.Now()
	if holiday, ok := annualHolidays.holidays[currTime.Format("2006-01-02")]; ok {
		return !holiday.isOffDay
	} else {
		if currTime.Weekday() == time.Saturday || currTime.Weekday() == time.Sunday {
			return false
		} else {
			return true
		}
	}
}
