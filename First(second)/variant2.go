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

func getNearestHoliday(callend Dates) (time.Time, string) {
	for _, v := range callend {
		str := "2006-01-02"
		dateHoliday, err := time.Parse(str, v.Date)
		checkErr(err)
		if dateHoliday.After(t) {
			return dateHoliday, v.LocalName
		}
	}
	return time.Time{}, "No holiday"
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func isWeekend(day time.Time) bool {
	if day.Weekday() == time.Saturday || day.Weekday() == time.Sunday {
		return true
	}
	return false
}

func main() {
	fmt.Printf("Сейчас: %s\n", t)
	url := fmt.Sprintf("https://date.nager.at/api/v2/publicholidays/" + strconv.Itoa(t.Year()) + "/UA")

	body := getAllHolidays(url)
	callend := Dates{}
	jsonErr := json.Unmarshal(body, &callend)
	checkErr(jsonErr)

	var dateHoliday, strName = getNearestHoliday(callend)
	diff := dateHoliday.Sub(t)
	duration := int(math.Ceil(diff.Hours() / 24))
	layoutISO := "2006-01-02"
	fmt.Printf("Nearest holiday = %v\n",  (dateHoliday.Format(layoutISO)))
	fmt.Printf("Name holiday = %v\n",strName)
	fmt.Printf("the holiday will be through: %v days \n", duration)

	oneDayLater := dateHoliday.AddDate(0, 0, 1)
	if isWeekend(oneDayLater) {
		fmt.Println(" The weekend adjoins the holiday!!! ")
	}
}
