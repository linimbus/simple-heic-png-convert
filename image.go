package main

import (
	"fmt"
	"image/jpeg"
	"image/png"
	"os"

	"github.com/adrium/goheif"
)

func ConvertHeic2Png(fileIn string, fileOut string, level string) error {
	fd, err := os.Open(fileIn)
	if err != nil {
		return fmt.Errorf("unable to open %s: %w", fileIn, err)
	}
	defer fd.Close()

	img, err := goheif.Decode(fd)
	if err != nil {
		return fmt.Errorf("unable to decode %s: %w", fileIn, err)
	}

	fOut, err := os.Create(fileOut)
	if err != nil {
		return fmt.Errorf("unable to create %s: %v", fileOut, err)
	}
	defer fOut.Close()

	var cmpLevel png.CompressionLevel

	switch level {
	case "Low":
		cmpLevel = png.NoCompression
	case "Middle":
		cmpLevel = png.BestSpeed
	case "High:":
		cmpLevel = png.BestCompression
	default:
		cmpLevel = png.DefaultCompression
	}

	pngenc := png.Encoder{CompressionLevel: cmpLevel}
	err = pngenc.Encode(fOut, img)
	if err != nil {
		return fmt.Errorf("unable to encode %s: %w", fileOut, err)
	}

	return nil
}

func ConvertHeic2Jpeg(fileIn string, fileOut string, quality int) error {
	fd, err := os.Open(fileIn)
	if err != nil {
		return fmt.Errorf("unable to open %s: %w", fileIn, err)
	}
	defer fd.Close()

	img, err := goheif.Decode(fd)
	if err != nil {
		return fmt.Errorf("unable to decode %s: %w", fileIn, err)
	}

	fOut, err := os.Create(fileOut)
	if err != nil {
		return fmt.Errorf("unable to create %s: %v", fileOut, err)
	}
	defer fOut.Close()

	err = jpeg.Encode(fOut, img, &jpeg.Options{Quality: quality})
	if err != nil {
		return fmt.Errorf("unable to encode %s: %w", fileOut, err)
	}

	return nil
}
