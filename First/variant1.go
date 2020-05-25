package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"time"
)

//Dates ...
type Dates []struct {
	Date        string      `json:"date"`
	LocalName   string      `json:"localName"`
	Name        string      `json:"name"`
	CountryCode string      `json:"countryCode"`
	Fixed       bool        `json:"fixed"`
	Global      bool        `json:"global"`
	Counties    interface{} `json:"counties"`
	LaunchYear  interface{} `json:"launchYear"`
	Type        string      `json:"type"`
}

var t = time.Now()

func getAllHolidays(url string) []byte {
	res, err := http.Get(url)
	checkErr(err)
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	checkErr(err)
	return body
}

func getNearestHoliday(callend Dates, mapDates map[int]string) int {
	first := 0
	min := 0
	for _, v := range callend {
		str := "2006-01-02"
		dates, err := time.Parse(str, v.Date)
		checkErr(err)
		difference := dates.Sub(t)
		duration := int(math.Ceil(difference.Hours() / 24))

		if first == 0 && duration >= 0 {
			min = duration
			mapDates[min] = v.LocalName
			first = 1
		}
		if duration >= 0 {
			if duration <= min {
				min := duration
				mapDates[min] = v.LocalName
			}
		}
	}
	return min
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	fmt.Printf("Сейчас: %s\n", t)

	url := fmt.Sprintf("https://date.nager.at/api/v2/publicholidays/" + strconv.Itoa(t.Year()) + "/UA")

	body := getAllHolidays(url)
	callend := Dates{}
	jsonErr := json.Unmarshal(body, &callend)
	checkErr(jsonErr)

	var mapDates = make(map[int]string)
	min := getNearestHoliday(callend, mapDates)

	if min == 0 {
		fmt.Printf("Today holiday = %v\n", mapDates[min])
	} else {
		fmt.Printf("Nearest holiday = %v\n", mapDates[min])

		str := "2006-01-02"
		for j := range callend {
			if callend[j].LocalName == mapDates[min] {
				fmt.Printf("day = %v\n", callend[j].Date)
				dates, err := time.Parse(str, callend[j].Date)
				checkErr(err)
				fmt.Println(dates.Weekday())

				oneDayLater := dates.AddDate(0, 0, 1)
				if (oneDayLater.Weekday() == time.Saturday || oneDayLater.Weekday() == time.Sunday) {
					fmt.Println(" The weekend adjoins the holiday!!! ")
				}
				break
			}
		}
	}
}
