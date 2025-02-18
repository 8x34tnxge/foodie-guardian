package main

import (
	"fmt"
	"log"
	"os"
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

func main() {
	token := os.Getenv("FOODIE_GARDIAN_TOKEN")
	secret := os.Getenv("FOODIE_GARDIAN_SECRET")
	timezone := os.Getenv("TZ")

	feishuBot := feishu.NewClient(token, secret)

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

	_, err = scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(
			gocron.NewAtTime(9, 30, 0),
			gocron.NewAtTime(14, 00, 0),
		)),
		gocron.NewTask(
			WorkNotification, feishuBot, &annualHolidays, GenerateOrderingNotificationMsg,
		))
	if err != nil {
		log.Fatal(err)
	}

	_, err = scheduler.NewJob(
		gocron.DailyJob(1, gocron.NewAtTimes(
			gocron.NewAtTime(11, 50, 0),
			gocron.NewAtTime(17, 00, 0),
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
