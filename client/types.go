package client

type Location struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

type Weather struct {
	Current Current `json:"current" env-required:"true"`
}

type Current struct {
	Temp float64 `json:"temperature_2m"`
	Wind float64 `json:"wind_speed_10m"`
}
