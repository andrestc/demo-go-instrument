package weather

import (
	"encoding/json"
	"net/http"
)

type cityResult struct {
	Main mainResult
}

type mainResult struct {
	Temp float64
}

type CityTemp struct {
	Temp float64
	Unit string
}

func GetCityTemp(name string) (*CityTemp, error) {
	client, err := getClient()
	if err != nil {
		return nil, err
	}
	resp, err := client.do(http.MethodGet, "data/2.5/weather", map[string]string{"q": name, "units": "metric"})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result cityResult
	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return nil, err
	}
	return &CityTemp{Temp: result.Main.Temp, Unit: "C"}, nil
}
