package main

import (
	"fmt"
	"math"
	"time"

	"github.com/phpdave11/gofpdf"
)

func GenerateSalesReportPDF(filename string, campaign Campaign, salesByType map[string]TicketTypeSales) error {
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
	pdf.CellFormat(0, 8, translate("Rapport de Vente"), "", 1, "C", false, 0, "")
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
	pdf.MultiCell(0, 6, translate("Date du spectacle: "+showDate), "", "L", false)
	pdf.MultiCell(0, 6, translate("Titre du spectacle: "+showTitle), "", "L", false)
	pdf.Ln(3)

	writeParagraph := func(text string) {
		pdf.MultiCell(0, 6, translate(text), "", "L", false)
		pdf.Ln(2)
	}
	writeParagraph("Merci encore d’avoir joué chez nous! Vous contribuez directement à bâtir quelque chose de durable pour la communauté musicale de Montréal. Ça compte énormément pour nous.")
	writeParagraph("Voici votre résumé:")

	ticketTypes := []struct {
		key   string
		label string
	}{
		{key: "prevente", label: "Prévente"},
		{key: "etudiant", label: "Étudiant"},
		{key: "admission generale", label: "Admission Générale"},
		{key: doorString, label: "Porte"},
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
	pdf.MultiCell(0, 6, translate(fmt.Sprintf("Sous-total: $%.2f", float64(subtotal)/100)), "", "L", false)
	pdf.MultiCell(0, 6, translate(fmt.Sprintf("Montant à nous facturer (80%%) = $%.2f", amountToInvoice/100)), "", "L", false)
	pdf.SetFont("Arial", "", 11)
	pdf.Ln(4)

	writeParagraph("S’il vous plaît adresser la facture à:")
	writeParagraph("Cabaret Jazz L’Entracte O.S.B.L.\n1112 rue Sainte-Catherine Ouest\nMontréal, Qc, H3B 1H4\n514-772-6044\nraphael@entractejazz.org")
	writeParagraph("Veuillez noter qu’à cause de notre situation de financement en tant qu’OSBL, nos paiements peuvent être temporairement délayés. Nous vous remercions de votre compréhension, et nous vous prions de bien vouloir refléter nos plus sincères excuses aux membres individuels de vos ensembles. Nous vous ferons suivre votre paiement dès que possible.")
	writeParagraph("Merci d’avoir été des nôtres, et au plaisir de vous recevoir de nouveau. À très bientôt!")
	writeParagraph("Chaleureusement,\nL'équipe de L'Entracte")

	if err := pdf.OutputFileAndClose(filename); err != nil {
		return fmt.Errorf("write sales report PDF %q: %w", filename, err)
	}
	return nil
}
