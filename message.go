package main

import (
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

func GeneratePrevOrderingMsg() string {
	msgs := []string{
		"报餐时间到！",
		"爱护身体是一种投资，要投资你的胃，让身体更加健康。",
		"吃饱了才有力气干活！",
		"时间紧迫！",
		"吃货们，",
		"报餐不仅是一种行为，更是一种态度，",
		"还没报餐的同学，你们的午餐/晚餐有着落了吗？",
		"温馨提示：为了确保大家都能吃到心仪的午餐/晚餐，",
		"报餐啦！报餐啦！",
		"Hi 大家，",
	}
	randomIndex := rand.Intn(len(msgs))

	return msgs[randomIndex]
}

func GenerateOrderingNotificationMsg() string {
	msgs := []string{
		"工作再忙也要记得报餐哦~",
		"别忘记报餐哦~",
		"快点报餐啦~",
		"请尽快完成报餐，谢谢合作！",
		"快去快去报餐！让我们一起为美食而战！",
		"快去快去报餐！让我们一起分享美食的快乐！",
		"快来报餐吧~",
		"报餐报餐！请大家积极配合",
	}
	randomIndex := rand.Intn(len(msgs))

	return GeneratePrevOrderingMsg() + msgs[randomIndex] + GenerateKaomoji()
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
	}
	randomIndex := rand.Intn(len(msgs))

	return msgs[randomIndex]
}

func GenerateMealNotificationMsg() string {
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
	}
	randomIndex := rand.Intn(len(msgs))

	return GeneratePrevMealMsg() + msgs[randomIndex] + GenerateKaomoji()
}
