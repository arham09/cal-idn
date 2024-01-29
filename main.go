package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/arham09/cal-idn/cache"
)

const (
	reset = "\033[0m"
	red   = "\033[31m"
	green = "\033[32m"

	key       = "holiday-schedule"
	cacheFile = "./cache.json"

	holidayExternalData = "https://raw.githubusercontent.com/guangrei/APIHariLibur_V2/main/holidays.json"
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
		date := firstOfMonth.AddDate(0, 0, i-1)
		holiday, ok := holidaySchedules[date.Format("2006-01-02")]["summary"]
		if ok {
			holidayInMonth = append(holidayInMonth, fmt.Sprintf("%d %s - %s", date.Day(), date.Month().String(), holiday))
		}

		printDate(date.Day(), ok, date.Day() == now.Day())
	}

	fmt.Printf("\n")
	for i := 8 - int(firstDay); i <= daysInMonth; i += 7 {
		for j := i; j < i+7; j++ {
			if j > daysInMonth {
				break
			}
			date := firstOfMonth.AddDate(0, 0, j-1)
			holiday, ok := holidaySchedules[date.Format("2006-01-02")]["summary"]
			if ok {
				holidayInMonth = append(holidayInMonth, fmt.Sprintf("%d %s - %s", date.Day(), date.Month().String(), holiday))
			}

			printDate(date.Day(), ok, date.Day() == now.Day())
		}
		fmt.Printf("\n")
	}

	fmt.Printf("\n")
	for i := 0; i < len(holidayInMonth); i++ {
		fmt.Println(holidayInMonth[i])
	}
}

func getHoliday() (map[string]map[string]string, error) {
	c := cache.NewCache()

	if err := c.LoadFromFile(cacheFile); err != nil {
		if os.IsNotExist(err) {
			goto FETCH
		}
		return nil, err
	}

	if data, exists := c.Get(key); exists {
		return data, nil
	}

FETCH:
	resp, err := http.Get(holidayExternalData)
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

	c.Set(key, data, 30*24*time.Hour)
	if err := c.SaveToFile(cacheFile); err != nil {
		return nil, err
	}

	return data, nil
}

func printDate(date int, isHoliday, isToday bool) {
	color := reset
	if isToday {
		color = green
	}
	if isHoliday {
		color = red
	}

	fmt.Printf("%s%2d%s ", color, date, reset)
}
