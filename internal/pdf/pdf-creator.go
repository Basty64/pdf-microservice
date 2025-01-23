package pdf

import (
	"bytes"
	"fmt"
	"github.com/go-pdf/fpdf"
	"github.com/skip2/go-qrcode"
	"image"
	"image/color"
	"image/png"
	"strings"
)

func GeneratePDF(data RequestData) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Define colors
	greyColor := color.RGBA{R: 240, G: 240, B: 240, A: 255}

	//Load fonts
	pdf.AddUTF8Font("Roboto-Regular", "", "./assets/Roboto-Regular.ttf")
	pdf.AddUTF8Font("Roboto-Bold", "B", "./assets/Roboto-Bold.ttf")

	// Set initial Y
	currentY := 7.0

	// Header information
	pdf.SetFont("Roboto-Bold", "", 13)
	pdf.SetTextColor(0, 0, 0)

	// Header text
	pdf.SetTextColor(0, 0, 0)
	headerText := fmt.Sprintf("%s  %s TRIP %s", data.TripsStart, data.TripsEnd, strings.ToUpper(data.Location))
	pdf.SetXY(10, currentY)
	pdf.Cell(0, 6, headerText)
	currentY += 8

	// Triangle
	pdf.SetFillColor(0, 0, 0)
	pdf.Polygon([]fpdf.PointType{{X: 10, Y: currentY - 2}, {X: 13, Y: currentY + 1}, {X: 10, Y: currentY + 4}}, "F")

	// Header line
	pdf.Line(10, 14, 200, 14)

	// QR Code
	qrCodeBytes, err := generateQRCode("https://github.com/go-pdf/fpdf")
	if err != nil {
		return nil, fmt.Errorf("failed to generate qr code: %w", err)
	}
	rdr := bytes.NewReader(qrCodeBytes)
	opt := fpdf.ImageOptions{ImageType: "png", ReadDpi: true}

	pdf.RegisterImageOptionsReader("qr-code", opt, rdr)

	pdf.Image("qr-code", 170, 15, 35, 35, false, "", 0, "")

	// Prepared for
	pdf.SetFont("Roboto-Regular", "", 9)
	pdf.SetXY(10, currentY)
	pdf.Cell(0, 4, "PREPARED FOR")
	currentY += 5

	pdf.SetXY(10, currentY)
	pdf.Cell(0, 5, data.PreparedFor)
	currentY += 6

	// Reservation code
	pdf.SetFont("Roboto-Regular", "", 10)
	pdf.SetXY(10, currentY)
	pdf.Cell(0, 4, fmt.Sprintf("RESERVATION CODE %s", data.ReservationCode))
	currentY += 5

	// Partial payment
	pdf.SetXY(10, currentY)
	pdf.Cell(0, 4, data.PartialPrepayment)
	currentY += 5

	// Final price
	pdf.SetXY(10, currentY)
	pdf.Cell(0, 5, fmt.Sprintf("FINAL PRICE: %s (taxes included)", data.FinalPrice))
	currentY += 10

	// Flights data
	for _, flight := range data.Flights {

		// Flight grey background
		pdf.SetFillColor(int(greyColor.R), int(greyColor.G), int(greyColor.B))
		pdf.Rect(10, currentY-5, 190-10, 45, "F")

		// DEPARTURE
		pdf.SetTextColor(0, 0, 0)
		pdf.SetXY(10, currentY)
		pdf.SetFillColor(0, 0, 0)
		pdf.Polygon([]fpdf.PointType{{X: 10, Y: currentY + 1}, {X: 13, Y: currentY + 4}, {X: 10, Y: currentY + 7}}, "F")
		pdf.Cell(0, 5, " DEPARTURE: "+strings.ToUpper(flight.Departure)+" "+flight.DepartureTime+" Please verify flight times prior to departure")
		currentY += 7

		// Flight info
		pdf.SetXY(10, currentY)
		pdf.Cell(30, 4, "FLIGHT")
		pdf.SetX(100)
		pdf.Cell(40, 4, flight.Departure)
		pdf.SetX(150)
		pdf.Cell(30, 4, flight.Arrival)
		pdf.SetX(180)
		pdf.Cell(0, 4, "Aircraft:")
		currentY += 5

		pdf.SetFont("Roboto-Bold", "", 10)
		pdf.SetXY(10, currentY)
		pdf.Cell(30, 4, flight.FlightNumber)

		pdf.SetFont("Roboto-Regular", "", 10)
		pdf.SetX(100)
		pdf.SetX(150)
		pdf.SetX(180)
		pdf.Cell(0, 4, flight.Aircraft)
		currentY += 5

		pdf.SetXY(10, currentY)
		pdf.Cell(30, 4, fmt.Sprintf("Airline: %s", flight.Airline))
		pdf.SetX(100)
		pdf.Cell(40, 4, "Departing At:")
		pdf.SetX(150)
		pdf.Cell(30, 4, "Arriving At:")
		pdf.SetX(180)
		pdf.Cell(0, 4, "Distance (in Miles):")
		currentY += 4

		pdf.SetXY(10, currentY)
		pdf.Cell(30, 4, fmt.Sprintf("Class: %s", flight.Class))
		pdf.SetX(100)
		pdf.Cell(40, 4, flight.DepartureTime)
		pdf.SetX(150)
		pdf.Cell(30, 4, flight.ArrivalTime)
		pdf.SetX(180)
		pdf.Cell(0, 4, flight.Distance)
		currentY += 4

		pdf.SetXY(10, currentY)
		pdf.Cell(30, 4, fmt.Sprintf("Status: %s", flight.Status))
		pdf.SetX(180)
		pdf.Cell(0, 4, fmt.Sprintf("Stop(s):"))
		currentY += 4

		pdf.SetX(180)
		pdf.Cell(0, 4, fmt.Sprintf("Meals: %s", flight.Meals))
		currentY += 8

		pdf.SetFont("Roboto-Regular", "", 10)
		pdf.SetXY(10, currentY)
		pdf.Cell(30, 5, fmt.Sprintf("Passenger Name: %s", flight.PassengerName))
		pdf.SetX(100)
		pdf.Cell(20, 5, "Seats:")
		pdf.SetX(140)
		pdf.Cell(0, 5, fmt.Sprintf("Booking: %s", flight.Booking))
		currentY += 5
		pdf.SetXY(10, currentY)
		pdf.Cell(0, 5, fmt.Sprintf("%s", flight.Seats))
		currentY += 10

	}

	pdf.Line(10, 260, 190, 260)

	// Terms and conditions
	pdf.SetFont("Roboto-Regular", "", 8)
	pdf.SetY(270)
	pdf.SetX(10)
	termsAndConditions := "TERMS AND CONDITIONS\n\n" +
		"If air carriage is provided for hereon, this document must be exchanged for a ticket and at such time prior to departure as may be required\n" +
		"by the rules and regulations of the carrier to whom the document is directed\n\n" +
		"If this document is issued in respect to baggage, the passenger must also have a passenger ticket and bag- baggage check, since this\n" +
		"document is not the baggage check described by Article 4 of The Hague Protocol or The Warsaw Convention as amended by the Hague\n" +
		"Protocol, 1955 or the Baggage Identification Tag described by Article 3 of the Montreal Convention 1999.\n\n" +
		"This document and any carriage or services for which it provides are subject to the currently effective and applicable tariffs, conditions of\n" +
		"carriage, rules and regulations of the issuer and of the carrier to whom it is directed and of any carrier performing carriage or services\n" +
		"under the ticket or tickets issued in exchange for this order, and to all the terms and conditions under which non-air carriage services are\n" +
		"arranged, offered or provided, as well as the laws of the country wherein these services are arranged, offered or provided.\n\n" +
		"In issuing this document, the issuer acts only as agent for the carrier or carriers furnishing the carriage or the person arranging or\n" +
		"supplying the services described hereon and the issuer shall not be liable for any loss, injury, damage or delay which is occasioned by\n" +
		"such carrier or person, for which results from such carrier or person performing or failing to perform the carriage or other services, or from\n" +
		"such carrier or person failing to honour this document.\n\n" +
		"The honouring carrier or person providing services re- serves the right to obtain authorisation from the issuing carrier prior to honouring\n" +
		"this document.\n\n" +
		"The use of the term issuer, carrier or person includes all owners, subsidiaries and affiliates of such issuer, carrier or person and any\n" +
		"person with whom such issuer, carrier or person has contracted to perform the carriage or services provided for hereon.\n\n" +
		"The acceptance of this document by the person named on the face hereof, or by the person purchasing this document on behalf of such\n" +
		"named person, shall be deemed to be consent to and acceptance by such person or persons of these conditions."

	pdf.MultiCell(0, 3, termsAndConditions, "", "", false)

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to write PDF to buffer: %w", err)
	}

	return buf.Bytes(), nil
}

func generateQRCode(data string) ([]byte, error) {
	qrCode, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	img, _, err := image.Decode(bytes.NewReader(qrCode))
	if err != nil {
		return nil, fmt.Errorf("failed to decode qr code: %w", err)
	}

	err = png.Encode(&buf, img)
	if err != nil {
		return nil, fmt.Errorf("failed to encode qr code to png: %w", err)
	}

	return buf.Bytes(), nil
}

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
