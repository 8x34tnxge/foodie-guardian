package main

import (
	"context"
	"fmt"
	"google.golang.org/genai"
	"math/rand"
)

func GenerateKaomoji() string {

	kaomojis := []string{
		"~^_^~",
		"d(d＇∀＇)",
		"d(`･∀･)b",
		"(*´∀`)~♥",
		"σ`∀´)σ",
		"(〃∀〃)",
		"ε٩(๑> ₃ <)۶з",
		"(◔౪◔)",
		"ლ(╹◡╹ლ)",
		"(⁰▿⁰)",
		"(๑´ㅂ`๑)",
		"_(:3 ⌒ﾞ)_",
		"(灬ºωº灬)",
		"(❛◡❛✿)",
		"( ^ω^)",
		"( ﾟ∀ﾟ)o彡ﾟ",
		"( ～'ω')～",
		"(≧∀≦)ゞ",
		"(・ε・)",
		"(=´ω`=)",
	}
	randomIndex := rand.Intn(len(kaomojis))

	return kaomojis[randomIndex]
}

func GetMessageFromLLM(gemini_api LLM_API, current_time string, command string) (string, error) {
	ctx := context.Background()
	// The client gets the API key from the environment variable `GEMINI_API_KEY`.
	client, err := genai.NewClient(ctx, nil)
	if err != nil {
		return "", err
	}
	prompt := "当前时间为" + current_time + "，提醒内容为" + command + "提醒，请返回提醒语句"
	fmt.Println("prompt: " + prompt)
	config := genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: "食堂每日会提供午餐和晚餐，但是需要当天提前报餐才能去食用。"},
				{Text: "而你是一个在公司群里的报餐提醒机器人，负责在可以报餐以及到点吃饭的时间，推送对应的提醒语句。"},
				{Text: "每次需要提醒的时候，你收到的内容包括当前时间(e.g. 17:30)以及当前提醒的内容（e.g. 报餐、吃饭）。"},
				{Text: "你的任务是，根据收到的内容，生成对应的提醒语句，返回内容必须仅包含提醒语句，提醒语句格式要求为“提醒” + “好处”，你可以在任意位置添加 emoji 和颜文字，可以适当添加语气词和符号。"},
				{Text: "注意，报餐必须是当天报餐当天食用，例如上午是报当天的午餐，下午是报当天的晚餐，不会跨天报餐。"},
			},
		},
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-pro",
		genai.Text(prompt),
		&config,
	)
	if err != nil {
		return "", err
	}
	if result.Text() == "" {
		fmt.Println("Receive empty message")
		return "", fmt.Errorf("Receive empty message from LLM")
	}
	fmt.Println("result is" + result.Text())
	return result.Text(), nil
}

func GeneratePrevOrderingMsg() string {
	msgs := []string{
		"报餐时间到！",
		"爱护身体是一种投资，要投资你的胃，让身体更加健康。",
		"吃饱了才有力气干活！",
		"时间紧迫！",
		"吃货们，",
		"报餐不仅是一种行为，更是一种态度，",
		"还没报餐的同学，你们的午餐/晚餐有着落了吗？",
		"报餐啦！报餐啦！",
		"Hi 大家，",
		"嗨嗨~",
	}
	randomIndex := rand.Intn(len(msgs))

	return msgs[randomIndex]
}

func GenerateOrderingNotificationMsg(gemini_api *LLM_API, trigger_time string) string {
	command := "报餐"
	var msg string

	legacy_mode := gemini_api == nil

	if !legacy_mode {
		msg_from_llm, err := GetMessageFromLLM(*gemini_api, trigger_time, command) // TODO: get msg from Gemini
		if err != nil {
			fmt.Printf("Get message failed, error %s", err.Error())
			legacy_mode = true
		} else {
			msg = msg_from_llm
		}
	}
	if legacy_mode {
		// legacy process
		msgs := []string{
			"工作再忙也要记得报餐哦~",
			"别忘记报餐哦~",
			"快点报餐啦~",
			"请尽快完成报餐，谢谢合作！",
			"快去快去报餐！让我们一起为美食而战！",
			"快去快去报餐！让我们一起分享美食的快乐！",
			"快来报餐吧~",
			"报餐报餐！请大家积极配合",
			"赶紧报餐，美味等你！",
			"报餐时间到，别错过哦！",
			"今天你报餐了吗？",
			"午餐/晚餐报餐进行中，请及时执行。",
			"为了你的胃，请尽快报餐！",
			"报餐啦！错过今天再等明天！",
			"亲爱的同学们，请尽快完成报餐。",
			"每天一餐，快乐一天！快来报餐吧。",
		}
		randomIndex := rand.Intn(len(msgs))
		msg = GeneratePrevOrderingMsg() + msgs[randomIndex] + GenerateKaomoji()
	}

	return msg
}

func GeneratePrevMealMsg() string {
	msgs := []string{
		"吃饭时间到！",
		"吃饱了才有力气干活！",
		"时间紧迫！",
		"吃货们，",
		"吃饭不仅是一种行为，更是一种态度，",
		"吃饭啦！吃饭啦！",
		"Hi 大家，",
		"人是铁饭是钢，一顿不吃饿得慌，",
		"嗨嗨~",
	}
	randomIndex := rand.Intn(len(msgs))

	return msgs[randomIndex]
}

func GenerateMealNotificationMsg(gemini_api *LLM_API, trigger_time string) string {
	command := "吃饭"
	var msg string

	legacy_mode := gemini_api == nil

	if !legacy_mode {
		msg_from_llm, err := GetMessageFromLLM(*gemini_api, trigger_time, command) // TODO: get msg from Gemini
		if err != nil {
			fmt.Printf("Get message failed, error %s", err.Error())
			legacy_mode = true
		} else {
			msg = msg_from_llm
		}
	}

	if legacy_mode {
		// legacy process when gemini_api
		msgs := []string{
			"工作再忙也要记得吃饭哦~",
			"别忘记吃饭哦~",
			"快点去吃饭啦~",
			"请尽快完成吃饭，谢谢合作！",
			"快去吃饭吧~",
			"吃饭吃饭！请大家积极配合",
			"吃点好吃的犒劳一下自己吧~",
			"快去补充能量吧！",
			"还没有吃的同学，你们的胃在抗议啦！",
			"是时候去吃饭啦！",
			"快来展现你的吃货精神吧！",
			"肚子饿了吗？该吃饭了！",
			"辛苦了，去吃顿好的吧！",
			"身体是革命的本钱，记得按时吃饭！",
			"别让饥饿影响你的工作效率哦！",
			"休息一下，享受美食吧！",
			"吃饭时间到，别磨蹭啦！",
			"为了健康，请按时就餐。",
			"饭点啦，快去补充体力！",
		}
		randomIndex := rand.Intn(len(msgs))
		msg = GeneratePrevMealMsg() + msgs[randomIndex] + GenerateKaomoji()
	}

	return msg
}
