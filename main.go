package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	reset = "\033[0m"
	red   = "\033[31m"
)

func main() {
	var (
		now   = time.Now().Local()
		month = now.Month()
		year  = now.Year()

		firstOfMonth = time.Date(year, month, 1, 0, 0, 0, 0, time.Local)
		lastOfMonth  = firstOfMonth.AddDate(0, 1, -1)

		firstDay    = firstOfMonth.Weekday()
		daysInMonth = lastOfMonth.Day()

		holidayInMonth = make([]string, 0, daysInMonth)
	)

	holidaySchedules, err := getHoliday()
	if err != nil {
		fmt.Printf("Failed to fetch holiday %v", err)
		return
	}

	fmt.Printf("Calendar for %s - %d\n\n", month, year)
	fmt.Println("Su Mo Tu We Th Fr Sa")
	for i := 0; i < int(firstDay); i++ {
		fmt.Printf("   ")
	}

	for i := 1; i <= 7-int(firstDay); i++ {
		color := reset
		date := firstOfMonth.AddDate(0, 0, i-1)
		if holiday, ok := holidaySchedules[date.Format("2006-01-02")]["summary"]; ok {
			holidayInMonth = append(holidayInMonth, fmt.Sprintf("%d %s: %s", date.Day(), date.Month().String(), holiday))
			color = red
		}

		fmt.Printf("%s%2d%s ", color, date.Day(), reset)
	}

	fmt.Printf("\n")
	for i := 8 - int(firstDay); i <= daysInMonth; i += 7 {
		for j := i; j < i+7; j++ {
			if j > daysInMonth {
				break
			}
			date := firstOfMonth.AddDate(0, 0, j-1)
			color := reset
			if holiday, ok := holidaySchedules[date.Format("2006-01-02")]["summary"]; ok {
				holidayInMonth = append(holidayInMonth, fmt.Sprintf("%d %s: %s", date.Day(), date.Month().String(), holiday))
				color = red
			}

			fmt.Printf("%s%2d%s ", color, date.Day(), reset)
		}
		fmt.Printf("\n")
	}

	fmt.Printf("\n")
	for i := 0; i < len(holidayInMonth); i++ {
		fmt.Println(holidayInMonth[i])
	}
}

func getHoliday() (map[string]map[string]string, error) {
	resp, err := http.Get("https://raw.githubusercontent.com/guangrei/APIHariLibur_V2/main/holidays.json")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var data map[string]map[string]string
	if err := json.Unmarshal(body, &data); err != nil {
		return nil, err
	}

	return data, nil
}
