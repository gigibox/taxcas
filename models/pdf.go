package models

import (
	"github.com/jung-kurt/gofpdf"
)

func Image2PDF(dst, image string) (error) {
	pdf := gofpdf.New("P", "pt", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "", 11)
	pdf.Image(image, 0, 0, 595.28, 841.89, false, "png", 0, "")
	return pdf.OutputFileAndClose(dst)
}