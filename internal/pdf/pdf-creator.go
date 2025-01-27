package pdf

import (
	"bytes"
	"fmt"
	"github.com/go-pdf/fpdf"
	"image/color"
	"math"
	"pdf-microservice/internal/models"
	"pdf-microservice/internal/qrcodes"
	"strings"
)

const (
	verifyFlights          = "Please verify flight times prior to departure"
	departure              = "DEPARTURE: "
	departingAt            = "Departing At:"
	arrivingAt             = "Arriving At:"
	aircraft               = "Aircraft:"
	distance               = "Distance (in Miles):"
	stops                  = "Stop(s):"
	meals                  = "Meals:"
	passengerName          = "Passenger Name:"
	seats                  = "Seats:"
	booking                = "Booking:"
	termsAndConditionsLOGO = "TERMS AND CONDITIONS\n\n"
)

func GeneratePDF(data models.RequestData) ([]byte, error) {
	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Define colors
	greyColor := color.RGBA{R: 240, G: 240, B: 240, A: 255}
	darkGreyColor := color.RGBA{R: 150, G: 150, B: 150, A: 255}

	//Load fonts
	pdf.AddUTF8Font("Roboto-Regular", "", "./assets/Roboto-Regular.ttf")
	pdf.AddUTF8Font("Roboto-Bold", "", "./assets/Roboto-Bold.ttf")

	// Set initial X
	currentX := 10.0

	// Set initial Y
	currentY := 7.0

	// Header text
	pdf.SetFont("Roboto-Bold", "", 13)
	pdf.SetTextColor(0, 0, 0)
	headerText := fmt.Sprintf("%s    %s         %s", data.TripsStart, data.TripsEnd, strings.ToUpper(data.Location))
	pdf.SetXY(10, currentY)
	pdf.Cell(0, 6, headerText)
	pdf.SetXY(65, currentY+0.3)
	pdf.SetFont("Roboto-Regular", "", 10)
	pdf.CellFormat(0, 6, "TRIP", "", 0, "L", false, 0, "")
	currentY = 21

	// Triangle
	pdf.SetFillColor(0, 0, 0)
	pdf.Polygon([]fpdf.PointType{{X: 37, Y: 8}, {X: 39, Y: 9.5}, {X: 37, Y: 11}}, "F")

	// QR Code
	qrCodeBytes, err := qrcodes.GenerateQRCode("https://github.com/go-pdf/fpdf")
	if err != nil {
		return nil, fmt.Errorf("failed to generate qr code: %w", err)
	}
	rdr := bytes.NewReader(qrCodeBytes)
	opt := fpdf.ImageOptions{ImageType: "png", ReadDpi: true}

	pdf.RegisterImageOptionsReader("qr-code", opt, rdr)

	pdf.Image("qr-code", 168.5, 12, 35, 35, false, "", 0, "")

	// Header line
	pdf.Line(10, 13, 200, 13)

	// Prepared for
	pdf.SetFont("Roboto-Regular", "", 11)
	pdf.SetXY(10, currentY)
	pdf.Cell(0, 4, "PREPARED FOR")
	currentY = 25.5

	pdf.SetXY(10, currentY)
	pdf.Cell(0, 4, data.PreparedFor)
	currentY = 30

	// Reservation code
	pdf.SetXY(10, currentY)
	pdf.Cell(0, 4, fmt.Sprintf("RESERVATION CODE     %s", data.ReservationCode))
	currentY = 34.5

	// Partial payment
	pdf.SetXY(10, currentY)
	pdf.Cell(0, 4, data.PartialPrepayment)
	currentY = 39

	// Final price
	pdf.SetXY(10, currentY)
	pdf.Cell(0, 4, fmt.Sprintf("FINAL PRICE: %s (taxes included)", data.FinalPrice))
	currentY = 51

	// Flights data
	for _, flight := range data.Flights {

		// 2nd Line
		pdf.Line(10, currentY, 200, currentY)

		// DEPARTURE TEXT-LINE
		depatureDate := fmt.Sprintf(strings.ToUpper(flight.Date))
		pdf.SetTextColor(0, 0, 0)
		pdf.SetXY(10, currentY)
		pdf.SetFillColor(0, 0, 0)

		currentY += 0.5 //51.5
		pdf.Polygon([]fpdf.PointType{{X: 12, Y: currentY + 0.5}, {X: 14, Y: currentY + 2}, {X: 12, Y: currentY + 3.5}}, "F")

		pdf.SetXY(14, currentY)
		currentX = pdf.GetX()
		pdf.Cell(0, 5, departure)
		currentX += pdf.GetStringWidth(departure)

		pdf.SetFont("Roboto-Bold", "", 11)
		pdf.SetXY(currentX+1, currentY)
		pdf.Cell(0, 5, depatureDate)
		currentX += pdf.GetStringWidth(depatureDate)

		pdf.SetFont("Roboto-Regular", "", 8)
		pdf.SetTextColor(int(darkGreyColor.R), int(darkGreyColor.G), int(darkGreyColor.B))
		pdf.SetXY(currentX+4, currentY)
		pdf.Cell(0, 6, verifyFlights)
		currentY += 4 //55.5

		pdf.SetFont("Roboto-Regular", "", 11)
		pdf.SetTextColor(0, 0, 0)

		// Flight grey background
		pdf.SetFillColor(int(greyColor.R), int(greyColor.G), int(greyColor.B))
		pdf.Rect(10, currentY+1, 50, 40, "F")

		// Table top-line
		currentY += 1 //56.5
		pdf.SetDrawColor(int(darkGreyColor.R), int(darkGreyColor.G), int(darkGreyColor.B))
		pdf.Line(60.5, currentY, 200, currentY)

		// Table bottom-line
		pdf.Line(60.5, currentY+40, 200, currentY+40)

		// Table left-line
		pdf.Line(60.5, currentY, 60.5, currentY+40)

		// Table right-line
		pdf.Line(200, currentY, 200, currentY+40)

		// Dotted lines
		pdf.SetLineCapStyle("round")
		rectSize := 0.1
		spaceLen := 0.5
		drawDashedRectLine(pdf, 158, currentY, 158, currentY+40, rectSize, spaceLen)
		drawDashedRectLine(pdf, 60.5, currentY+21.5, 158, currentY+21.5, rectSize, spaceLen)
		drawDashedRectLine(pdf, 105, currentY+21.5, 105, currentY+40, rectSize, spaceLen)

		// FLIGHT
		currentX = 11
		currentY += 1 //57.5
		pdf.SetXY(currentX, currentY)
		pdf.Cell(30, 4, "FLIGHT")

		// Flight number
		pdf.SetFont("Roboto-Bold", "", 11)
		currentY += 8.5 //66
		pdf.SetXY(currentX, currentY)
		pdf.Cell(30, 4, flight.FlightNumber)

		// Airline
		pdf.SetFont("Roboto-Regular", "", 8)
		currentY += 8 //74
		pdf.SetXY(currentX, currentY)
		pdf.Cell(30, 4, fmt.Sprintf("Airline: %s", flight.Airline))
		pdf.SetX(60)

		// Class
		currentY += 4 //78
		pdf.SetXY(currentX, currentY)
		flightClass := fmt.Sprintf("Class: %s", strings.ToUpper(flight.Class))
		pdf.Cell(30, 4, flightClass)

		// Status
		currentY += 8 //86
		pdf.SetXY(currentX, currentY)
		pdf.Cell(30, 4, fmt.Sprintf("Status: %s", flight.Status))
		// ----------------------------------------------------------

		// FLIGHT AIRPORTS CODES
		// Start airport-code
		pdf.SetFont("Roboto-Regular", "", 11)
		currentX = 64
		currentY -= 27.5 //58.5
		startGeo := strings.Split(flight.Departure, " ")
		airportCode := startGeo[0]
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 4, airportCode)

		// Finish airport-code
		currentX = 110
		finishGeo := strings.Split(flight.Arrival, " ")
		airportCode = finishGeo[0]
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 4, airportCode)

		// Triangle
		currentX = 108
		pdf.SetFillColor(0, 0, 0)
		pdf.SetXY(currentX, currentY)
		pdf.Polygon([]fpdf.PointType{{X: currentX, Y: currentY}, {X: currentX + 2, Y: currentY + 1.5}, {X: currentX, Y: currentY + 3}}, "F")

		// Start airport city and country
		currentX = 64
		currentY += 4.5 //62
		pdf.SetFont("Roboto-Regular", "", 8)
		cityAndCountry := strings.ToUpper(startGeo[1] + " " + startGeo[2])
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 4, cityAndCountry)

		// Finish airport city and country
		currentX = 110
		cityAndCountry = strings.ToUpper(finishGeo[1] + " " + finishGeo[2])
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 4, cityAndCountry)

		//Departing At
		currentX = 64
		currentY += 18 //80
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 4, departingAt)

		// Departure date
		currentX = 64
		currentY += 3 //83
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 4, flight.Date)

		// Departure time
		pdf.SetFont("Roboto-Regular", "", 12)
		currentX = 64
		currentY += 4 //87
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 4, flight.DepartureTime)

		//Arriving at
		pdf.SetFont("Roboto-Regular", "", 8)
		currentX = 112
		currentY -= 7 //80
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 4, arrivingAt)

		// Arrival date
		currentX = 112
		currentY += 3 //83
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 4, flight.Date)

		// Arrival time
		pdf.SetFont("Roboto-Regular", "", 12)
		currentX = 112
		currentY += 4 //87
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 4, flight.ArrivalTime)
		//----------------------------------------------------------

		// FLIGHT RIGHT DATA
		// Aircraft
		pdf.SetFont("Roboto-Regular", "", 8)
		currentX = 160
		currentY -= 29 //58
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 4, aircraft)

		// Aircraft number
		currentX = 160
		currentY += 4 // 62
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 4, flight.Aircraft)

		// Distance
		currentX = 160
		currentY += 4 // 66
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 4, distance)

		// Distance measure
		currentX = 160
		currentY += 4 // 70
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 4, flight.Distance)

		// Stops
		currentX = 160
		currentY += 4 // 74
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 4, stops)

		//Space
		currentY += 5 // 79

		// Meals
		currentX = 160
		currentY += 4 // 83
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 4, meals)

		// Meals data
		currentX = 160
		currentY += 4 //87
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 4, flight.Meals)

		//---------------------------------------------

		// BOTTOM TABLE
		// Grey background
		pdf.SetFillColor(int(greyColor.R), int(greyColor.G), int(greyColor.B))
		pdf.Rect(10, currentY+12.5, 190, 4.5, "F")

		// Dotted lines
		drawDashedRectLine(pdf, 99, currentY+13, 99, currentY+21, rectSize, spaceLen)
		drawDashedRectLine(pdf, 154, currentY+13, 154, currentY+21, rectSize, spaceLen)

		// Passenger name
		currentX = 10
		currentY += 12 //100
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 5, passengerName)

		// Passenger data
		currentX = 10
		currentY += 4 // 104
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 5, flight.PassengerName)

		// Seats
		currentX = 100
		currentY -= 4 // 100
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 5, seats)

		// Seats data
		currentX = 100
		currentY += 4 // 104
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 5, flight.Seats)

		//	Booking
		currentX = 155
		currentY -= 4 // 100
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 5, booking)

		//	Booking Data
		currentX = 155
		currentY += 4 // 104
		pdf.SetXY(currentX, currentY)
		pdf.Cell(0, 5, flight.Booking)

		currentY += 8
	}
	//----------------------------------

	pdf.Line(10, 260, 190, 260)

	// Terms and conditions
	currentX = 10.0
	currentY += 5.0

	// Выводим заголовок жирным шрифтом
	pdf.SetFont("Roboto-Bold", "", 12)
	pdf.SetXY(currentX, currentY)
	pdf.Cell(0, 5, termsAndConditionsLOGO)

	// Выводим основной текст обычным шрифтом
	pdf.SetFont("Roboto-Regular", "", 9)
	currentY += 5.0 // Сдвигаем Y на высоту заголовка
	pdf.SetXY(currentX, currentY)

	// Вывод многострочного текста, обернув в MultiCell
	termsAndConditions := "If air carriage is provided for hereon, this document must be exchanged for a ticket and at such time prior to departure as may be required\n" +
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

	currentX = 10.0
	currentY += 3.0
	pdf.SetXY(currentX, currentY)
	// Выводим многострочный текст с использованием MultiCell, устанавливаем ширину 0, т.е. на всю строку
	pdf.MultiCell(0, 3, termsAndConditions, "", "", false)

	var buf bytes.Buffer
	err = pdf.Output(&buf)
	if err != nil {
		return nil, fmt.Errorf("failed to write PDF to buffer: %w", err)
	}

	return buf.Bytes(), nil
}

// drawDashedRectLine рисует пунктирную линию из квадратов
func drawDashedRectLine(pdf *fpdf.Fpdf, x1, y1, x2, y2, rectSize, spaceLen float64) {
	dx := x2 - x1
	dy := y2 - y1
	lineLen := math.Hypot(dx, dy)

	segmentLen := rectSize + spaceLen
	segments := int(lineLen / segmentLen)

	currentX := x1
	currentY := y1

	for i := 0; i < segments; i++ {
		pdf.Rect(currentX, currentY, rectSize, rectSize, "FD")
		currentX = x1 + dx*(float64(i)*segmentLen+segmentLen)/lineLen
		currentY = y1 + dy*(float64(i)*segmentLen+segmentLen)/lineLen
	}
	pdf.Rect(currentX, currentY, rectSize, rectSize, "FD")
}
