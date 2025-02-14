package main

import (
	"context"
	"encoding/json"
	"fmt"
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

func (annualHolidays *AnnualHolidays) ReadFromApi() error {
	newHolidays, err := FetchAnnualHolidayInfo(time.Now().Year())
	if err != nil {
		return err
	}

	newHolidays.storePath = annualHolidays.storePath
	*annualHolidays = newHolidays
	return nil
}

func (annualHolidays *AnnualHolidays) Update() {
	if annualHolidays.year != time.Now().Year() {
		annualHolidays.ReadFromApi()
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

func FetchFromApiWithRetry(url string, maxRetries uint, timeout time.Duration) ([]byte, error) {
	var body []byte
	for i := uint(0); i <= maxRetries; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			return nil, fmt.Errorf("创建请求失败: %w", err)
		}

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			if ctx.Err() == context.DeadlineExceeded {
				log.Printf("尝试 %d: 请求超时", i+1)
				continue // 继续下一次循环
			}
			return nil, fmt.Errorf("请求失败: %w", err)
		}
		defer resp.Body.Close()

		body, err = io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("读取响应失败: %w", err)
		}

		if resp.StatusCode == http.StatusOK {
			return body, nil // 请求成功
		}

		log.Printf("尝试 %d: 状态码 %d", i+1, resp.StatusCode)
		time.Sleep(time.Second) // 等待一段时间后重试
	}

	return nil, fmt.Errorf("尝试 %d 次后仍然失败", maxRetries+1)
}

func FetchAnnualHolidayInfo(year int) (AnnualHolidays, error) {
	holidayApiUrl := "https://holiday.cyi.me/api/holidays?year=" + strconv.Itoa(year)

	holidays, err := FetchFromApiWithRetry(holidayApiUrl, 10, 10*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	return ParseAnnualHolidayInfo(holidays)
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
