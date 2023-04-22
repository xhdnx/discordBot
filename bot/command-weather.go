package bot

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

const URL string = "https://api.openweathermap.org/data/2.5/weather?"

type WeatherData struct {
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:speed`
	} `json:"wind"`
	Name string `json:"name"`
}

func getCurrentWeather(message string) *discordgo.MessageSend {
	r, _ := regexp.Compile(`\d{5}`)
	zip := r.FindString(message)

	if zip == "" {
		return &discordgo.MessageSend{
			Content: "Sorry that ZIP code doesn't look right",
		}
	}

	weatherURL := fmt.Sprintf("%szip=%s&util=imperial&appid=%s", URL, zip, OpenWeatherToken)

	client := http.Client{Timeout: 5 * time.Second}
	response, err := client.Get(weatherURL)

	if err != nil {
		return &discordgo.MessageSend{
			Content: "Sorry, there was an error trying to get the weather",
		}
	}

	body, _ := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	var data WeatherData
	err = json.Unmarshal(body, &data)

	if err != nil {
		log.Fatal("Data format has been changed")
	}

	city := data.Name
	conditions := data.Weather[0].Description
	temperature := strconv.FormatFloat(data.Main.Temp, 'f', 2, 64)
	humidity := strconv.Itoa(data.Main.Humidity)
	wind := strconv.FormatFloat(data.Wind.Speed, 'f', 2, 64)

	// Build Discord embed response
	embed := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{{
			Type:        discordgo.EmbedTypeRich,
			Title:       "Current Weather",
			Description: "Weather for " + city,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Conditions",
					Value:  conditions,
					Inline: true,
				},
				{
					Name:   "Temperature",
					Value:  temperature + "°F",
					Inline: true,
				},
				{
					Name:   "Humidity",
					Value:  humidity + "%",
					Inline: true,
				},
				{
					Name:   "Wind",
					Value:  wind + " mph",
					Inline: true,
				},
			},
		},
		},
	}

	return embed

}
