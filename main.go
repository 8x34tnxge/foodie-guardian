package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CatchZeng/feishu/pkg/feishu"
	"github.com/go-co-op/gocron/v2"
)

type FeishuBotAPI struct {
	token  string
	secret string
}

type LLM_API struct {
	api_key string
	model   string
}

func WorkNotification(bot *feishu.Client, annualHolidays *AnnualHolidays, f func(*LLM_API, string) string, api *LLM_API, trigger_time string) {
	annualHolidays.Update()
	if annualHolidays.DetermineIfTodayShouldWork() {
		fmt.Println("Today is work day, need to send notification")
		for {
			_, resp, err := bot.Send(feishu.NewPostMessage().AppendZHContent([]feishu.PostItem{feishu.NewText(f(api, trigger_time))}))
			if err == nil && 0 == resp.Code {
				break
			}
			fmt.Println("Send message failed, wait 10 ms and try again")
			time.Sleep(10 * time.Millisecond)
		}
		fmt.Println("Send notification success")
	}
}

func ParseTimes(times string) ([]gocron.AtTime, error) {
	timeVector := strings.Split(times, ";")
	retTimes := make([]gocron.AtTime, len(timeVector))
	for idx, timeString := range timeVector {
		timeSplits := strings.Split(timeString, ":")
		if len(timeSplits) != 3 {
			return retTimes, fmt.Errorf("time format is invalid, it should be HH:MM:SS")
		}

		hours, err := strconv.ParseUint(timeSplits[0], 10, 0)
		if err != nil {
			return retTimes, err
		}
		minutes, err := strconv.ParseUint(timeSplits[1], 10, 0)
		if err != nil {
			return retTimes, err
		}
		seconds, err := strconv.ParseUint(timeSplits[2], 10, 0)
		if err != nil {
			return retTimes, err
		}

		retTimes[idx] = gocron.NewAtTime(uint(hours), uint(minutes), uint(seconds))
	}

	return retTimes, nil
}

func main() {
	// 1. 定义命令行参数
	orderFlag := flag.Bool("order", false, "触发点餐提醒")
	mealFlag := flag.Bool("meal", false, "触发吃饭提醒")
	time_param := flag.String("time", "", "触发时间")
	flag.Parse() // 解析传入的参数

	feishu_bot_api := FeishuBotAPI{os.Getenv("FEISHU_BOT_TOKEN"), os.Getenv("FEISHU_BOT_SECRET")}
	feishuBot := feishu.NewClient(feishu_bot_api.token, feishu_bot_api.secret)
	fmt.Printf("Create client, token %s, secret %s\n", feishu_bot_api.token, feishu_bot_api.secret)

	gemini_setting := LLM_API{os.Getenv("GEMINI_API_KEY"), os.Getenv("GEMINI_MODEL")}

	// TODO: migrate holiday parse into scheduler
	annualHolidays, err := FetchAnnualHolidayInfo(time.Now().Year())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Fetch holiday info success")

	fmt.Println("Parse succeed, begin to process command")
	if *orderFlag {
		WorkNotification(feishuBot, &annualHolidays, GenerateOrderingNotificationMsg, &gemini_setting, *time_param)
		fmt.Println("Send order notification")
	}
	if *mealFlag {
		WorkNotification(feishuBot, &annualHolidays, GenerateMealNotificationMsg, &gemini_setting, *time_param)
		fmt.Println("Send meal notification")
	}
}
