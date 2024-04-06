package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
)

// Define structs for Division and District
type Division struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	BnName       string `json:"bn_name"`
	DistrictsURL string `json:"url"`
}

type District struct {
	ID           string `json:"id"`
	DivisionID   string `json:"division_id"`
	Name         string `json:"name"`
	BnName       string `json:"bn_name"`
	Lat          string `json:"lat"`
	Lon          string `json:"lon"`
	DistrictsURL string `json:"url"`
}

// Weather struct to hold weather information
type Weather struct {
	Main        string  `json:"main"`
	Description string  `json:"description"`
	Temperature float64 `json:"temp"`
}

// Function to search for districts by division name
func searchDistrictsByDivision(divisionName string, divisions []Division, districts []District) []District {
	// Create an empty slice to store the districts
	var divisionDistricts []District

	// Find division ID by division name
	var divisionID string
	for _, division := range divisions {
		if strings.EqualFold(division.Name, divisionName) {
			divisionID = division.ID
			break
		}
	}

	// Find districts belonging to the division
	for _, district := range districts {
		if district.DivisionID == divisionID {
			divisionDistricts = append(divisionDistricts, district)
		}
	}

	return divisionDistricts
}

// Function to fetch weather forecast for a district
func getWeatherForecast(latitude, longitude string) (Weather, error) {
	apiKey := "644260e7b85256544fef67f1a6f963f9"
	url := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?lat=%s&lon=%s&units=metric&appid=%s", latitude, longitude, apiKey)

	resp, err := http.Get(url)
	if err != nil {
		return Weather{}, err
	}
	defer resp.Body.Close()

	var weatherResp struct {
		Weather []Weather `json:"weather"`
		Main    struct {
			Temperature float64 `json:"temp"`
		} `json:"main"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&weatherResp); err != nil {
		return Weather{}, err
	}

	if len(weatherResp.Weather) == 0 {
		return Weather{}, fmt.Errorf("no weather data found")
	}

	return Weather{
		Main:        weatherResp.Weather[0].Main,
		Description: weatherResp.Weather[0].Description,
		Temperature: weatherResp.Main.Temperature,
	}, nil
}

func main() {
	// Read JSON data containing divisions
	divisionsData, err := os.ReadFile("divisions.json")
	if err != nil {
		log.Fatalf("Error reading divisions file: %v", err)
	}

	// Read JSON data containing districts
	districtsData, err := os.ReadFile("districts.json")
	if err != nil {
		log.Fatalf("Error reading districts file: %v", err)
	}

	// Parse JSON data into structs
	var divisions []Division
	if err := json.Unmarshal(divisionsData, &divisions); err != nil {
		log.Fatalf("Error parsing divisions JSON: %v", err)
	}

	var districts []District
	if err := json.Unmarshal(districtsData, &districts); err != nil {
		log.Fatalf("Error parsing districts JSON: %v", err)
	}

	// Example: Search for districts by division name
	divisionName := "dhaka" // Replace with the division name you want to search for
	divisionDistricts := searchDistrictsByDivision(divisionName, divisions, districts)

	// Print out the districts found
	fmt.Println("Districts in", divisionName, "Division:")
	for _, district := range divisionDistricts {
		fmt.Println("District Name:", district.Name)

		// Fetch weather forecast for the district
		weather, err := getWeatherForecast(district.Lat, district.Lon)
		if err != nil {
			log.Printf("Error fetching weather forecast for %s: %v", district.Name, err)
			continue
		}

		// Print weather forecast
		fmt.Printf("Weather Forecast: %s - %s, Temperature: %.1f°C\n", weather.Main, weather.Description, weather.Temperature)
	}
}
