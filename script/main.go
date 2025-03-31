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
	} `json:"daily"`
}

func getTemperatureEmoji(temp float64) string {
	switch {
	case temp >= 35:
		return "🔥" // Very hot
	case temp >= 30:
		return "☀️" // Hot
	case temp >= 25:
		return "⛅" // Warm
	case temp >= 20:
		return "😊" // Pleasant
	case temp >= 15:
		return "🌥️" // Cool
	case temp >= 10:
		return "❄️" // Cold
	default:
		return "🥶" // Very cold
	}
}

func main() {
	resp, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=10.823&longitude=106.6296&daily=temperature_2m_max,temperature_2m_min&timezone=Asia%2FBangkok&forecast_days=7")
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
	for _, temp := range weather.Daily.Temperature2mMax {
		maxTemps = append(maxTemps, fmt.Sprintf("%.1f %s", temp, getTemperatureEmoji(temp)))
	}

	// Each row must be []string with same length as headers
	rows := [][]string{
		append([]string{"Temp °C"}, maxTemps...),
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
