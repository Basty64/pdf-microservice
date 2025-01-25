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
