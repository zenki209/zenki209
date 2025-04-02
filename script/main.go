package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/fbiville/markdown-table-formatter/pkg/markdown"
)

type WeatherResponse struct {
	Daily struct {
		Time             []string  `json:"time"`
		Temperature2mMax []float64 `json:"temperature_2m_max"`
		WeatherCode      []int     `json:"weather_code"`
	} `json:"daily"`
}

func getTemperatureEmoji(code int) string {
	codeMap := map[int]string{
		0:  "☀️",   // Clear sky
		1:  "🌤️",   // Mainly clear
		2:  "⛅",    // Partly cloudy
		3:  "☁️",   // Overcast
		45: "🌫️",   // Fog
		48: "🌫️",   // Depositing rime fog
		51: "🌦️",   // Light drizzle
		53: "🌦️",   // Moderate drizzle
		55: "🌧️",   // Dense drizzle
		56: "🌧️❄️", // Light freezing drizzle
		57: "🌧️❄️", // Dense freezing drizzle
		61: "🌧️",   // Slight rain
		63: "🌧️",   // Moderate rain
		65: "🌧️",   // Heavy rain
		66: "🌧️❄️", // Light freezing rain
		67: "🌧️❄️", // Heavy freezing rain
		71: "🌨️",   // Slight snow fall
		73: "🌨️",   // Moderate snow fall
		75: "❄️",   // Heavy snow fall
		77: "❄️",   // Snow grains
		80: "🌧️",   // Slight rain showers
		81: "🌧️",   // Moderate rain showers
		82: "🌧️",   // Violent rain showers
		85: "🌨️",   // Slight snow showers
		86: "❄️",   // Heavy snow showers
		95: "⛈️",   // Thunderstorm
		96: "⛈️⚡",  // Thunderstorm with slight hail
		99: "⛈️❄️", // Thunderstorm with heavy hail
	}

	if emoji, exists := codeMap[code]; exists {
		return emoji
	}
	return "❓"
}

func main() {
	resp, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=10.823&longitude=106.6296&daily=temperature_2m_max,weather_code&timezone=Asia%2FSingapore")
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var weather WeatherResponse
	if err := json.Unmarshal(body, &weather); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// for i, date := range weather.Daily.Time {
	// 	maxTemp := weather.Daily.Temperature2mMax[i]
	// 	fmt.Printf("Date: %s\n", date)
	// 	fmt.Printf("Max Temperature: %.1f°C %s\n", maxTemp, getTemperatureEmoji(maxTemp))
	// 	fmt.Println("---")
	// }
	//Prepare Header
	headers := append([]string{"Date"}, weather.Daily.Time...)
	var maxTemps []string
	var weatherCode []string
	for _, temp := range weather.Daily.Temperature2mMax {
		maxTemps = append(maxTemps, fmt.Sprintf("%.1f", temp))
	}
	for _, weather := range weather.Daily.WeatherCode {
		weatherCode = append(weatherCode, getTemperatureEmoji(weather))
	}

	// Each row must be []string with same length as headers
	rows := [][]string{
		append([]string{"Temp °C"}, maxTemps...),
		append([]string{"Weather"}, weatherCode...),
	}

	// Format table
	table, err := markdown.NewTableFormatterBuilder().
		WithPrettyPrint().
		Build(headers...).
		Format(rows)
	if err != nil {
		log.Fatalf("Error formatting markdown table: %v", err)
	}

	templateBytes, err := os.ReadFile("../template/readme.md.tpl")
	if err != nil {
		log.Fatalf("Error reading template: %v", err)
	}

	templateStr := string(templateBytes)
	finalReadme := strings.Replace(templateStr, "{{TABLE}}", table, 1)

	err = os.WriteFile("README.md", []byte(finalReadme), 0644)
	if err != nil {
		log.Fatalf("Error writing README.md: %v", err)
	}
}
