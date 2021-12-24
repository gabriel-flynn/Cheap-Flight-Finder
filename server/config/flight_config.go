package config

import (
	"encoding/json"
	"github.com/senseyeio/duration"
	"log"
	"os"
)

type Config struct {
	FlightParameters []FlightSearch
	SaveToDynamo     bool
}

type FlightSearch struct {
	SrcAirport    string
	DestAirport   string
	NumPassengers int
	StartWindow   duration.Duration
	EndWindow     duration.Duration
}

type configDTO struct {
	FlightParameters []flightSearchDTO `json:"flights"`
	SaveToDynamo     bool              `json:"saveToDynamo"`
}

type flightSearchDTO struct {
	SrcAirport    string `json:"srcAirport"`
	DestAirport   string `json:"destAirport"`
	NumPassengers int    `json:"numPassengers"`
	StartWindow   string `json:"startWindow"`
	EndWindow     string `json:"endWindow"`
}

var config *Config

func GetConfig() *Config {
	//try the default location and config name
	if config == nil {
		LoadConfiguration("config.json")
	}
	return config
}

func LoadConfiguration(file string) *Config {
	var configDTO configDTO
	configFile, err := os.Open(file)
	if err != nil {
		log.Fatal(err.Error())
	}
	defer configFile.Close()
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&configDTO)

	flightParameters := make([]FlightSearch, len(configDTO.FlightParameters))
	for indx, flightParam := range configDTO.FlightParameters {

		startWindow, err := duration.ParseISO8601(flightParam.StartWindow)
		if err != nil {
			log.Fatal(err)
		}

		endWindow, err := duration.ParseISO8601(flightParam.EndWindow)
		if err != nil {
			log.Fatal(err)
		}

		flightParameters[indx] = FlightSearch{
			SrcAirport:    flightParam.SrcAirport,
			DestAirport:   flightParam.DestAirport,
			NumPassengers: flightParam.NumPassengers,
			StartWindow:   startWindow,
			EndWindow:     endWindow,
		}
	}

	config = &Config{FlightParameters: flightParameters, SaveToDynamo: configDTO.SaveToDynamo}
	return config
}
