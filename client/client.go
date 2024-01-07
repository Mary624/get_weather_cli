package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

type WeatherClient struct {
	apiKey string
	client http.Client
}

const (
	hostLonLan     = "api.openweathermap.org"
	hostWeather    = "api.open-meteo.com"
	getCoordsPath  = "geo/1.0/direct"
	getWeatherPath = "/v1/forecast"
	metric         = "metric" // Celsius
	exclude        = "current"
	limit          = 1
	tempStr        = "temperature_2m"
	windStr        = "wind_speed_10m"
	msUnit         = "ms"
	kmhUnit        = "kmh"
	CUnit          = "celsius"
	FUnit          = "fahrenheit"
)

func New(apiKey string) WeatherClient {
	return WeatherClient{
		apiKey: apiKey,
		client: http.Client{},
	}
}

func (cl WeatherClient) GetWeather(city string, isF, isKMH bool) (string, error) {
	loc, err := cl.getLocation(city)
	if err != nil {
		return "", fmt.Errorf("can't get location:\n%s", err)
	}

	temp, err := cl.GetTemperature(loc.Lat, loc.Lon, isF)
	if err != nil {
		return "", fmt.Errorf("can't get temperature:\n%s", err)
	}

	wind, err := cl.GetWind(loc.Lat, loc.Lon, isKMH)
	if err != nil {
		return "", fmt.Errorf("can't get wind speed:\n%s", err)
	}

	unitsTemp := CUnit
	if isF == true {
		unitsTemp = FUnit
	}
	unitsWind := msUnit
	if isKMH == true {
		unitsWind = kmhUnit
	}

	return fmt.Sprintf(`The weather in %s is:
	Temperature: %.1f %s;
	Wind: %.1f %s`, city, temp, unitsTemp, wind, unitsWind),
		nil
}

func (cl WeatherClient) GetTemperature(lat, lon float64, isF bool) (float64, error) {
	units := CUnit
	if isF == true {
		units = FUnit
	}
	qTemp := url.Values{}
	qTemp.Add("latitude", fmt.Sprintf("%f", lat))
	qTemp.Add("longitude", fmt.Sprintf("%f", lon))
	qTemp.Add("current", tempStr)
	qTemp.Add("temperature_unit", units)

	dataTemp, err := cl.doRequest(hostWeather, getWeatherPath, qTemp)

	if err != nil {
		return 0, fmt.Errorf("can't do request:\n%s", err)
	}

	var resTemp Weather
	if err := json.Unmarshal(dataTemp, &resTemp); err != nil {
		return 0, fmt.Errorf("can't unmarshal json:\n%s", err)
	}

	temp := resTemp.Current.Temp

	return temp, nil
}

func (cl WeatherClient) GetWind(lat, lon float64, isKMH bool) (float64, error) {
	units := msUnit
	if isKMH == true {
		units = kmhUnit
	}
	qWind := url.Values{}
	qWind.Add("latitude", fmt.Sprintf("%f", lat))
	qWind.Add("longitude", fmt.Sprintf("%f", lon))
	qWind.Add("current", windStr)
	qWind.Add("wind_speed_unit", units)

	dataWind, err := cl.doRequest(hostWeather, getWeatherPath, qWind)

	if err != nil {
		return 0, fmt.Errorf("can't do request:\n%s", err)
	}

	var resWind Weather
	if err := json.Unmarshal(dataWind, &resWind); err != nil {
		return 0, fmt.Errorf("can't unmarshal json:\n%s", err)
	}

	wind := resWind.Current.Wind

	return wind, nil
}

func (cl WeatherClient) getLocation(city string) (*Location, error) {
	qLocation := url.Values{}
	qLocation.Add("q", city)
	qLocation.Add("limit", strconv.Itoa(limit))
	qLocation.Add("appid", cl.apiKey)

	data, err := cl.doRequest(hostLonLan, getCoordsPath, qLocation)

	if err != nil {
		return nil, fmt.Errorf("can't do request:\n%s", err)
	}

	var res []Location
	if err := json.Unmarshal(data, &res); err != nil {
		return nil, fmt.Errorf("can't unmarshal json:\n%s", err)
	}

	return &res[0], nil
}

func (c WeatherClient) doRequest(host, method string, query url.Values) (data []byte, err error) {
	u := url.URL{
		Scheme: "http",
		Host:   host,
		Path:   method,
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}
