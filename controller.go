package main

import (
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"gopkg.in/redis.v5"
)

const (
	baseUrl string = "http://www.ace.utoronto.ca/bookings/f?p=200:3:0::NO::"
)

func fetch(url string) (*http.Response, error) {
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	// TODO: Figure out how to this without hard-coding a cookie value
	req.AddCookie(&http.Cookie{
		Name:  "WWV_CUSTOM-F_1410000632844518_200",
		Value: "845B97E883105AC19173D1B9E65DE4B4",
	})

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func scrapeBuildingRooms(redis *redis.Client, building *Building) error {
	resp, err := fetch(fmt.Sprintf("%sP3_BLDG:%s", baseUrl, building.Code))
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}

	// Get building name
	building.Name = strings.TrimSpace(strings.TrimLeft(
		doc.Find("select#P3_BLDG option[selected=\"selected\"]").Text(),
		building.Code))

	// Get list of rooms for this building
	var rooms []string
	doc.Find("select#P3_ROOM option").Each(func(i int, s *goquery.Selection) {
		value, exists := s.Attr("value")

		if !exists || value == "%null%" {
			return
		}

		rooms = append(rooms, value)
	})

	var wg sync.WaitGroup
	wg.Add(len(rooms) - 1)

	for _, roomNumber := range rooms {
		go func(roomNumber string) {
			defer wg.Done()
			room := Room{Number: roomNumber}
			scrapeSingleRoom(redis, building.Code, &room)
			building.Rooms = append(building.Rooms, room)
		}(roomNumber)
	}

	wg.Wait()
	return nil
}

func scrapeSingleRoom(redis *redis.Client, buildingCode string, room *Room) error {
	resp, err := fetch(fmt.Sprintf("%sP3_BLDG,P3_ROOM:%s,%s", baseUrl, buildingCode, room.Number))
	if err != nil {
		return err
	}

	doc, err := goquery.NewDocumentFromResponse(resp)
	if err != nil {
		return err
	}

	dateMap := make(map[string][]Booking)
	var dates []string

	doc.Find("table.t3WeekCalendarAlternative1").Find("td").Each(func(i int, s *goquery.Selection) {
		if s.HasClass("t3Hour") {
			return
		}

		rawDate, exists := s.Find("input[type=\"hidden\"]").Attr("value")
		if !exists {
			return
		}

		date := rawDate[:8]
		// Remove seconds (190000 -> 1900)
		time := rawDate[8 : len(rawDate)-2]

		if time == "0000" {
			return
		}

		text := strings.TrimSpace(s.Find("div#apex_cal_data_grid_src").Text())
		// Replace multiple spaces with single space
		text = regexp.MustCompile(`[\n\r\s]+`).ReplaceAllString(text, " ")

		dateMap[date] = append(dateMap[date], Booking{
			Time:        time,
			Description: text,
		})

		dates = append(dates, date)
	})

	sort.Strings(dates)

	for date, bookings := range dateMap {
		room.Schedule = append(room.Schedule, Date{date, bookings})
	}

	return nil
}
