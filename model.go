package main

import (
	"gopkg.in/redis.v5"
)

type (
	Building struct {
		Code  string `json:"building_code"`
		Name  string `json:"name"`
		Rooms []Room `json:"rooms"`
	}

	Room struct {
		Number   string `json:"room_number"`
		Schedule []Date `json:"schedule"`
	}

	Date struct {
		Date     string    `json:"date"`
		Bookings []Booking `json:"bookings"`
	}

	Booking struct {
		Time        string `json:"time"`
		Description string `json:"description"`
	}
)

func getBuilding(redis *redis.Client, buildingCode string) (Building, error) {
	b := Building{Code: buildingCode}

	err := scrapeBuildingRooms(redis, &b)
	if err != nil {
		return Building{}, err
	}

	return b, nil
}

func getRoom(redis *redis.Client, buildingCode string, roomNumber string) (Room, error) {
	r := Room{Number: roomNumber}

	err := scrapeSingleRoom(redis, buildingCode, &r)
	if err != nil {
		return Room{}, err
	}

	return r, nil
}
