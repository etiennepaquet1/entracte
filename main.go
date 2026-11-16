package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"
	"unicode"

	"golang.org/x/text/unicode/norm"
)

const doorString = "door"

func main() {

	if err := LoadConfig(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	campaigns, err := GetCampaigns()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	mostRecentCampaign, err := GetMostRecentlyEndedCampaign(campaigns)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// if no campaigns ended in the last 24 hours, exit
	if mostRecentCampaign.EndDate < time.Now().Add(-24*time.Hour).Unix() {
		fmt.Println("No campaigns ended in the last 24 hours")
		os.Exit(0)
	}

	payments, err := GetPayments()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	fmt.Printf("Number of payments: %d\n", len(payments))

	salesByType := make(map[string]TicketTypeSales)
	for ticketType := range ticketTypes {
		salesByType[normalizeTicketType(ticketType)] = TicketTypeSales{0, 0}
	}
	salesByType[doorString] = TicketTypeSales{0, 0}

	for _, payment := range payments {
		// handle door payments in the last 24 hours of the campaign
		if payment.CampaignID == doorCampaignID && payment.Created >= time.Now().Add(-24*time.Hour).Unix() {
			sales := salesByType[doorString]
			sales.TicketsSold++
			sales.NetRevenue += payment.Amount
			salesByType[doorString] = sales
		}

		if payment.CampaignID != mostRecentCampaign.ID {
			continue
		}
		for _, paymentItem := range payment.Items {
			// fmt.Printf("Payment: %d, Item: %s\n", payment.Amount, paymentItem.RateTitle)
			ticketType := normalizeTicketType(paymentItem.RateTitle)
			if _, exists := ticketTypes[ticketType]; exists {
				sales := salesByType[ticketType]
				sales.TicketsSold++
				sales.NetRevenue += paymentItem.Amount
				salesByType[ticketType] = sales
			}
		}
	}

	fmt.Println("Last ended campaign:", mostRecentCampaign.Title)

	fmt.Println("Sales by ticket type:")
	for ticketType, sales := range salesByType {
		fmt.Printf("%s: %d tickets sold, $%.2f net revenue\n", ticketType, sales.TicketsSold, float64(sales.NetRevenue)/100)
	}

	report, err := GenerateSalesReportPDF(mostRecentCampaign, salesByType)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	filename := fmt.Sprintf("Sales report - %s - %s.pdf", mostRecentCampaign.Title, time.Now().Format("2006-01-02 15-04-05"))
	fileID, err := UploadReportToGoogleDrive(context.Background(), filename, report)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	fmt.Printf("Uploaded sales report to Google Drive: https://drive.google.com/file/d/%s/view\n", fileID)
}

func normalizeTicketType(ticketType string) string {
	decomposed := norm.NFD.String(strings.ToLower(ticketType))
	var normalized strings.Builder
	for _, r := range decomposed {
		if !unicode.Is(unicode.Mn, r) {
			normalized.WriteRune(r)
		}
	}
	return normalized.String()
}

type TicketTypeSales struct {
	TicketsSold int
	NetRevenue  int64
}
