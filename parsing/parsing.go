package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type Result struct {
	Links            Link                         `json:"links"`
	ElementCount     int                          `json:"element_count"`
	NearEarthObjects map[string][]NearEarthObject `json:"near_earth_objects"`
}

type NearEarthObject struct {
	Links                          Links                   `json:"links"`
	ID                             string                  `json:"id"`
	NeoReferenceId                 string                  `json:"neo_reference_id"`
	Name                           string                  `json:"name"`
	NasaJplUrl                     string                  `json:"nasa_jpl_url"`
	EstimatedDiameter              EstimatedDiameter       `json:"estimated_diameter"`
	IsPotentiallyHazardousAsteroid bool                    `json:"is_potentially_hazardous_asteroid"`
	CloseApproachData              []closeApproachDataForm `json:"close_approach_data"`
}

type closeApproachDataForm struct {
	CloseApproachDate      string           `json:"close_apporach_date"`
	CloseApproachDateFull  string           `json:"close_approach_date_full"`
	EpochDateCloseApproach int              `json:"epoch_date_close_approach"`
	RelativeVelocity       RelativeVelocity `json:"relative_velocity"`
	MissDistance           MissDistance     `json:"miss_distance"`
	OrbitingBody           string           `json:"orbiting_body"`
}

type RelativeVelocity struct {
	KilometersPerSecond string `json:"kilometers_per_second"`
	KilometersPerHour   string `json:"kilometers_per_hour"`
	MilesPerHour        string `json:"miles_per_hour"`
}

type MissDistance struct {
	Astronomical string `json:"astronomical"`
	Lunar        string `json:"lunar"`
	Kilometers   string `json:"kilometers"`
}

type EstimatedDiameter struct {
	Kilometers Kilometers
	Meters     Meters
	Miles      Miles
	Feet       Feet
}

type Kilometers struct {
	EstimatedDiameterMin float64 `json:"estimated_diameter_min"`
	EstimatedDiameterMax float64 `json:"estimated_diameter_max"`
}

type Meters struct {
	EstimatedDiameterMin float64 `json:"estimated_diameter_min"`
	EstimatedDiameterMax float64 `json:"estimated_diameter_max"`
}

type Miles struct {
	EstimatedDiameterMin float64 `json:"estimated_diameter_min"`
	EstimatedDiameterMax float64 `json:"estimated_diameter_max"`
}

type Feet struct {
	EstimatedDiameterMin float64 `json:"estimated_diameter_min"`
	EstimatedDiameterMax float64 `json:"estimated_diameter_max"`
}

type Links struct {
	Self string `json:"self"`
}

type Link struct {
	Next     string `json:"next"`
	Previous string `json:"previous"`
	Self     string `json:"self"`
}

func main() {
	jsonFile, err := os.Open("asteroids_2022_10_p1.json")

	if err != nil {
		log.Fatal(err)
	}

	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result Result
	err = json.Unmarshal([]byte(byteValue), &result)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(result.ElementCount)

	for key, val := range result.NearEarthObjects {
		fmt.Println(key)
		for i := 0; i < len(val); i++ {
			fmt.Println(val[i].ID, val[i].Name, val[i].EstimatedDiameter.Kilometers.EstimatedDiameterMax)
		}
		fmt.Println()
	}

	fmt.Println("Complete")
	defer jsonFile.Close()
}
