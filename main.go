package main

import (
	"DiscordBot/bot"
	"encoding/json"
	"io/ioutil"
	"log"
)

func main() {
	conf, _ := ioutil.ReadFile("conf.json")
	tokens := map[string]string{}
	err := json.Unmarshal(conf, &tokens)

	if err != nil {
		log.Fatal("We have an error while serialize conf file")
	}

	discordToken, openWeatherToken := tokens["DiscordBotToken"], tokens["OpenWeatherToken"]

	if openWeatherToken == "" {
		log.Fatal("Must set Open Weather token as env variable: OpenWeatherToken")
	}

	if discordToken == "" {
		log.Fatal("Must set Discord token as env variable: DiscordBotToken")
	}

	bot.BotToken = discordToken
	bot.OpenWeatherToken = openWeatherToken
	bot.Run()
}
