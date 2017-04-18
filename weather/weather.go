package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	apiURL = "http://api.openweathermap.org"
)

var (
	client *weatherClient
)

type weatherClient struct {
	httpClient *http.Client
	apiKey     string
}

func (w *weatherClient) do(method string, path string, params map[string]string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s?appid=%s", apiURL, path, w.apiKey)
	for k, v := range params {
		url += fmt.Sprintf("&%s=%s", k, v)
	}
	fmt.Printf("DEBUG: %s\n", url)
	request, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	return w.httpClient.Do(request)
}

func getClient() (*weatherClient, error) {
	apiKey := os.Getenv("WEATHER_API_KEY")
	if apiKey == "" {
		return nil, errors.New("must set WEATHER_API_KEY env")
	}
	if client == nil {
		client = &weatherClient{
			httpClient: &http.Client{Timeout: time.Second * 15},
			apiKey:     apiKey,
		}
	}
	return client, nil
}

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
	resp, err := client.do(http.MethodGet, fmt.Sprintf("data/2.5/weather"), map[string]string{"q": name, "units": "metric"})
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
