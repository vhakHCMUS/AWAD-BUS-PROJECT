package infrastructure

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/jung-kurt/gofpdf"
	"github.com/skip2/go-qrcode"
)

// PDFGenerator generates PDF tickets
type PDFGenerator struct {
	outputDir string
}

// NewPDFGenerator creates a new PDF generator
func NewPDFGenerator() *PDFGenerator {
	outputDir := os.Getenv("TICKET_PDF_DIR")
	if outputDir == "" {
		outputDir = "./uploads/tickets"
	}
	os.MkdirAll(outputDir, 0755)

	return &PDFGenerator{
		outputDir: outputDir,
	}
}

// GenerateTicket generates a PDF ticket with QR code
func (g *PDFGenerator) GenerateTicket(ticketCode, passengerName, fromCity, toCity, seatNumber, departureTime string) (string, error) {
	// Generate QR code
	qrPath, err := g.generateQRCode(ticketCode)
	if err != nil {
		return "", err
	}

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add header
	pdf.SetFont("Arial", "B", 20)
	pdf.Cell(190, 10, "BUS E-TICKET")
	pdf.Ln(15)

	// Ticket details
	pdf.SetFont("Arial", "", 12)
	pdf.Cell(60, 8, "Booking Code:")
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(130, 8, ticketCode)
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(60, 8, "Passenger Name:")
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(130, 8, passengerName)
	pdf.Ln(10)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(60, 8, "From:")
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(130, 8, fromCity)
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(60, 8, "To:")
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(130, 8, toCity)
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(60, 8, "Seat Number:")
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(130, 8, seatNumber)
	pdf.Ln(8)

	pdf.SetFont("Arial", "", 12)
	pdf.Cell(60, 8, "Departure:")
	pdf.SetFont("Arial", "B", 12)
	pdf.Cell(130, 8, departureTime)
	pdf.Ln(15)

	// Add QR code
	pdf.Image(qrPath, 80, 120, 50, 50, false, "", 0, "")

	pdf.Ln(55)
	pdf.SetFont("Arial", "I", 10)
	pdf.Cell(190, 5, "Please present this QR code at the boarding point")

	// Save PDF
	filename := fmt.Sprintf("ticket_%s.pdf", ticketCode)
	outputPath := filepath.Join(g.outputDir, filename)

	err = pdf.OutputFileAndClose(outputPath)
	if err != nil {
		return "", err
	}

	return outputPath, nil
}

func (g *PDFGenerator) generateQRCode(data string) (string, error) {
	qrDir := os.Getenv("QR_CODE_DIR")
	if qrDir == "" {
		qrDir = "./uploads/qrcodes"
	}
	os.MkdirAll(qrDir, 0755)

	filename := fmt.Sprintf("qr_%s.png", data)
	outputPath := filepath.Join(qrDir, filename)

	err := qrcode.WriteFile(data, qrcode.Medium, 256, outputPath)
	return outputPath, err
}
