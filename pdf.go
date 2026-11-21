package main

import (
	"bytes"
	"fmt"
	"math"
	"time"

	"github.com/phpdave11/gofpdf"
)

func GenerateSalesReportPDF(campaign Campaign, salesByType map[string]TicketTypeSales) ([]byte, error) {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.SetMargins(20, 18, 20)
	pdf.SetAutoPageBreak(true, 18)
	pdf.AddPage()
	pdf.SetTitle("Rapport de Vente", true)
	pdf.SetAuthor("Cabaret Jazz L'Entracte O.S.B.L.", true)
	translate := pdf.UnicodeTranslatorFromDescriptor("")

	pdf.SetFont("Arial", "B", 15)
	pdf.CellFormat(0, 8, translate("Cabaret Jazz L’Entracte O.S.B.L."), "", 1, "C", false, 0, "")
	pdf.SetFont("Arial", "B", 13)
	pdf.CellFormat(0, 8, translate("Rapport de Vente / Sales Report"), "", 1, "C", false, 0, "")
	pdf.Ln(4)

	pdf.SetFont("Arial", "", 11)
	showDate := "Non spécifiée"
	if campaign.StartDate > 0 {
		showDate = time.Unix(campaign.StartDate, 0).Format("02/01/2006")
	}
	showTitle := campaign.Title
	if showTitle == "" {
		showTitle = "Non spécifié"
	}
	pdf.MultiCell(0, 6, translate("Date du spectacle / Show date: "+showDate), "", "L", false)
	pdf.MultiCell(0, 6, translate("Titre du spectacle / Show title: "+showTitle), "", "L", false)
	pdf.Ln(3)

	writeParagraph := func(text string) {
		pdf.MultiCell(0, 6, translate(text), "", "L", false)
		pdf.Ln(2)
	}
	writeParagraph("Merci encore d’avoir joué chez nous! Vous contribuez directement à bâtir quelque chose de durable pour la communauté musicale de Montréal. Ça compte énormément pour nous.\n\nThank you again for playing with us! You are directly helping to build something lasting for Montreal’s music community. It means a lot to us.")
	writeParagraph("Voici votre résumé / Here is your summary:")

	ticketTypes := []struct {
		key   string
		label string
	}{
		{key: "prevente", label: "Prévente / Presale"},
		{key: "etudiant", label: "Étudiant / Student"},
		{key: "admission generale", label: "Admission Générale / General Admission"},
		{key: "etudiant - 2e set", label: "Étudiant - 2e set / Student - 2nd set"},
		{key: "admission generale - 2e set", label: "Admission Générale - 2e set / General Admission - 2nd set"},
	}

	var subtotal int64
	for _, ticketType := range ticketTypes {
		sales := salesByType[ticketType.key]
		subtotal += sales.NetRevenue
		unitPrice := 0.0
		if sales.TicketsSold > 0 {
			unitPrice = float64(sales.NetRevenue) / float64(sales.TicketsSold) / 100
		}
		line := fmt.Sprintf("%s: %d @ $%.2f = $%.2f", ticketType.label, sales.TicketsSold, unitPrice, float64(sales.NetRevenue)/100)
		pdf.MultiCell(0, 6, translate(line), "", "L", false)
	}

	amountToInvoice := math.Round(float64(subtotal) * 0.8)
	pdf.Ln(2)
	pdf.SetFont("Arial", "B", 11)
	pdf.MultiCell(0, 6, translate(fmt.Sprintf("Sous-total / Subtotal: $%.2f", float64(subtotal)/100)), "", "L", false)
	pdf.MultiCell(0, 6, translate(fmt.Sprintf("Montant à nous facturer / Amount to invoice us (80%%) = $%.2f", amountToInvoice/100)), "", "L", false)
	pdf.SetFont("Arial", "", 11)
	pdf.Ln(4)

	writeParagraph("S’il vous plaît adresser la facture à: / Please address the invoice to:")
	writeParagraph("Cabaret Jazz L’Entracte O.S.B.L.\n1112 rue Sainte-Catherine Ouest\nMontréal, Qc, H3B 1H4\n514-772-6044\nraphael@entractejazz.org")
	writeParagraph("Veuillez noter qu’à cause de notre situation de financement en tant qu’OSBL, nos paiements peuvent être temporairement délayés. Nous vous remercions de votre compréhension, et nous vous prions de bien vouloir refléter nos plus sincères excuses aux membres individuels de vos ensembles. Nous vous ferons suivre votre paiement dès que possible.\n\nPlease note that, due to our funding situation as a non-profit organization, our payments may be temporarily delayed. We thank you for your understanding, and kindly ask that you extend our sincerest apologies to the individual members of your ensembles. We will send your payment as soon as possible.")
	writeParagraph("Merci d’avoir été des nôtres, et au plaisir de vous recevoir de nouveau. À très bientôt!\nThank you for joining us, and we look forward to welcoming you again. See you very soon!")
	writeParagraph("Chaleureusement, / Warmly,\nL'équipe de L'Entracte / The L'Entracte Team")

	var output bytes.Buffer
	if err := pdf.Output(&output); err != nil {
		return nil, fmt.Errorf("generate sales report PDF: %w", err)
	}
	return output.Bytes(), nil
}
