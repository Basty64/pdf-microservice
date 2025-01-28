package models

type RequestData struct {
	TripsStart        string   `json:"trips_start"`
	TripsEnd          string   `json:"trips_end"`
	Location          string   `json:"location"`
	PreparedFor       string   `json:"prepared_for"`
	ReservationCode   string   `json:"reservation_code"`
	PartialPrepayment string   `json:"partial_prepayment"`
	FinalPrice        string   `json:"final_price"`
	Flights           []Flight `json:"flights"`
}

type Flight struct {
	Date          string `json:"date"`
	FlightNumber  string `json:"flight_number"`
	Airline       string `json:"airline"`
	Class         string `json:"class"`
	Status        string `json:"status"`
	Departure     string `json:"departure"`
	DepartureTime string `json:"departure_time"`
	Arrival       string `json:"arrival"`
	ArrivalTime   string `json:"arrival_time"`
	Aircraft      string `json:"aircraft"`
	Distance      string `json:"distance"`
	Meals         string `json:"meals"`
	PassengerName string `json:"passenger_name"`
	Seats         string `json:"seats"`
	Booking       string `json:"booking"`
}

type RequestDataNew struct {
	Ticket Ticket `json:"ticket"`
	User   User   `json:"user"`
	QRURL  string `json:"qr_url"`
}

type Ticket struct {
	ID          int           `json:"id"`
	Price       string        `json:"price"`
	Currency    string        `json:"currency"`
	Itineraries []Itineraries `json:"itineraries"`
	Airline     string        `json:"airline"`
}

type Itineraries struct {
	Duration string     `json:"duration"`
	Segments []Segments `json:"segments"`
	Stops    int        `json:"stops"`
}

type Segments struct {
	DepartureTime    string `json:"departure_time"`
	ArrivalTime      string `json:"arrival_time"`
	DepartureAirport string `json:"departure_airport"`
	ArrivalAirport   string `json:"arrival_airport"`
	Carrier          string `json:"carrier"`
	CarrierName      string `json:"carrier_name"`
	CarrierLogo      string `json:"carrier_logo"`
	Duration         string `json:"duration"`
}

type User struct {
	Email          string `json:"email"`
	FirstName      string `json:"first_name"`
	Gender         string `json:"gender"`
	LastName       string `json:"last_name"`
	Nationality    string `json:"nationality"`
	NumberPassport int    `json:"number_passport"`
	SeriaPassport  int    `json:"seria_passport"`
}
