package models

import (
	"image"

	"io"

	"github.com/signintech/gopdf"
)

func Image2Pdf(img image.Image, writer io.Writer) error {
	rec := img.Bounds()
	originalWidth := rec.Max.X - rec.Min.X
	originalHeight := rec.Max.Y - rec.Min.Y
	rect := gopdf.Rect{W: float64(originalWidth), H: float64(originalHeight)}

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{PageSize: rect})
	pdf.AddPage()
	if err := pdf.ImageFrom(img, 10, 10, &rect); err != nil {
		return err
	}
	if err := pdf.Write(writer); err != nil {
		return err
	}

	return nil
}
