package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/CatchZeng/feishu/pkg/feishu"
	"github.com/go-co-op/gocron/v2"
)

func WorkNotification(bot *feishu.Client, annualHolidays *AnnualHolidays, f func() string) {
	annualHolidays.Update()
	if annualHolidays.DetermineIfTodayShouldWork() {
		fmt.Println("Today is work day, need to send notification")
		for {
			_, resp, err := bot.Send(feishu.NewPostMessage().AppendZHContent([]feishu.PostItem{feishu.NewText(f())}))
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
	token := os.Getenv("FOODIE_GARDIAN_TOKEN")
	secret := os.Getenv("FOODIE_GARDIAN_SECRET")
	timezone := os.Getenv("TZ")
	orderNotifyTimes := os.Getenv("ORDER_NOTIFY_TIMES")
	mealNotifyTimes := os.Getenv("MEAL_NOTIFY_TIMES")

	feishuBot := feishu.NewClient(token, secret)
	fmt.Printf("Create client, token %s, secret %s\n", token, secret)

	// TODO: try database first, api later
	annualHolidays, err := FetchAnnualHolidayInfo(time.Now().Year())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Fetch holiday info success")

	location, err := time.LoadLocation(timezone)
	if err != nil {
		log.Fatal(err)
	}

	scheduler, err := gocron.NewScheduler(gocron.WithLocation(location))
	if err != nil {
		log.Fatal(err)
	}
	defer scheduler.Shutdown()

	fmt.Println("Parse times from env")
	orderNotifyTimesFromEnv, err := ParseTimes(orderNotifyTimes)
	if err != nil {
		log.Fatal(err)
	}
	mealNotifyTimesFromEnv, err := ParseTimes(mealNotifyTimes)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Parse succeed, begin to generate new jobs")
	_, err = scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(
			orderNotifyTimesFromEnv[0], orderNotifyTimesFromEnv[1:]...,
		)),
		gocron.NewTask(
			WorkNotification, feishuBot, &annualHolidays, GenerateOrderingNotificationMsg,
		))
	if err != nil {
		log.Fatal(err)
	}

	_, err = scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(
			mealNotifyTimesFromEnv[0], mealNotifyTimesFromEnv[1:]...,
		)),
		gocron.NewTask(
			WorkNotification, feishuBot, &annualHolidays, GenerateMealNotificationMsg,
		))
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("scheduler init succeed, start it")
	scheduler.Start()

	select {}
}
