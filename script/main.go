package main

import (
	"context"
	"html/template"
	"os"

	"github.com/hectormalot/omgo"
)

func tempToEmoji(temp float64) string {
	switch {
	case temp <= 10:
		return "🥶"
	case temp <= 20:
		return "🌥️"
	case temp <= 29:
		return "🌤️"
	case temp <= 34:
		return "🔥"
	default:
		return "☀️🥵"
	}
}

type WeatherInfo struct {
	LocationName string
	Latitude     float64
	Longitude    float64
	Temperature  float64
	Emoji        string
}

func main() {
	c, _ := omgo.NewClient()

	// Get the current weather for HCMC
	lat, lon := 10.823, 106.6296
	loc, _ := omgo.NewLocation(lat, lon)
	res, _ := c.CurrentWeather(context.Background(), loc, nil)

	data := WeatherInfo{
		LocationName: "Ho Chi Minh City",
		Latitude:     lat,
		Longitude:    lon,
		Temperature:  res.Temperature,
		Emoji:        tempToEmoji(res.Temperature),
	}
	// Load template from file
	tmpl, err := template.ParseFiles("../template/readme.md.tpl")
	if err != nil {
		panic(err)
	}

	// Output to stdout or file
	err = tmpl.Execute(os.Stdout, data)
	if err != nil {
		panic(err)
	}

	out, _ := os.Create("README.md")
	defer out.Close()
	tmpl.Execute(out, data)
}
