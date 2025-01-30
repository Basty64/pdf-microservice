package models

type RequestData struct {
	Ticket Ticket `json:"ticket"`
	User   User   `json:"user"`
}

type Ticket struct {
	ID               int           `json:"id"`
	Price            string        `json:"price"`
	Currency         string        `json:"currency"`
	Itineraries      []Itineraries `json:"itineraries"`
	Airline          string        `json:"airline"`
	FlightClass      string        `json:"flight_class"`
	StartCityName    string        `json:"start_city_name"`
	StartCountryName string        `json:"start_country_name"`
	FinalCityName    string        `json:"final_city_name"`
	FinalCountryName string        `json:"final_country_name"`
	QRURL            string        `json:"qr_url"`
}

type Itineraries struct {
	Duration string     `json:"duration"`
	Segments []Segments `json:"segments"`
	Stops    int        `json:"stops"`
}

type Segments struct {
	DepartureTime        string `json:"departure_time"`
	ArrivalTime          string `json:"arrival_time"`
	DepartureAirport     string `json:"departure_airport"`
	ArrivalAirport       string `json:"arrival_airport"`
	Carrier              string `json:"carrier"`
	CarrierName          string `json:"carrier_name"`
	CarrierLogo          string `json:"carrier_logo"`
	Duration             string `json:"duration"`
	DepartureCityName    string `json:"departure_city_name"`
	DepartureCountryName string `json:"departure_country_name"`
	ArrivalCityName      string `json:"arrival_city_name"`
	ArrivalCountryName   string `json:"arrival_country_name"`
}

type User struct {
	Email       string  `json:"email"`
	PhoneNumber string  `json:"phone_number"`
	Adults      []Adult `json:"adults"`
}

type Adult struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	BirthDate      string `json:"birth_date"`
	Gender         string `json:"gender"`
	SeriaPassport  int    `json:"seria_passport"`
	NumberPassport int    `json:"number_passport"`
	Nationality    string `json:"nationality"`
	ValidityPeriod string `json:"validity_period"`
}
