package gipher

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/gif"
	"os"
	"time"

	"github.com/fogleman/gg"
)

// CreateTimeStampGIF inserts current time as text on top of background GIF and saves it to output file. Returns an error.
func CreateTimeStampGIF(backgroundFilename, outputFilename, fontFilename string) error {
	// Background file is opened.
	backgroundFile, err := os.Open(backgroundFilename)
	if err != nil {
		return fmt.Errorf("error opening background file: %s\n", err)
	}
	defer backgroundFile.Close()

	// Background file is decoded to GIF.
	backgroundGIF, err := gif.DecodeAll(backgroundFile)
	if err != nil {
		return fmt.Errorf("error decoding background gif: %s\n", err)
	}

	const fontSize = 16 // Font size of the time stamp.

	// This GIF will hold all the frames in the end.
	outGif := &gif.GIF{
		Image:     make([]*image.Paletted, len(backgroundGIF.Image)),
		Delay:     backgroundGIF.Delay,
		LoopCount: backgroundGIF.LoopCount,
	}

	for i := 0; i < len(backgroundGIF.Image); i++ {
		// Current background frame is converted to RGBA.
		backgroundFrame := backgroundGIF.Image[i]
		rgbaFrame := image.NewRGBA(backgroundFrame.Bounds())
		draw.Draw(rgbaFrame, rgbaFrame.Bounds(), backgroundFrame, image.Point{}, draw.Over)

		// gg context is created from the RGBA frame.
		dc := gg.NewContextForRGBA(rgbaFrame)

		// Font face is loaded with the specified font size.
		if err := dc.LoadFontFace(fontFilename, fontSize); err != nil {
			return fmt.Errorf("error loading font face: %s\n", err)
		}

		// Timestamp is drawn.
		dc.SetRGB(0, 0, 0)
		timestamp := time.Now().Format("02.01. 15:04")
		dc.DrawStringAnchored(timestamp, float64(backgroundFrame.Bounds().Dx())/2, float64(backgroundFrame.Bounds().Dy())/2, 0.5, 0.5)
		dc.Fill()

		palette := append(backgroundFrame.Palette, color.Black)
		fmt.Println(palette)
		// Frame is converted back to paletted image.
		palettedFrame := imageToPaletted(dc.Image(), palette)

		// New frame is inserted to the final GIF.
		outGif.Image[i] = palettedFrame
	}

	// Final GIF is saved to output file.
	outputFile, err := os.OpenFile(outputFilename, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	err = gif.EncodeAll(outputFile, outGif)
	if err != nil {
		return fmt.Errorf("error encoding output file: %s\n", err)
	}

	// If all has gone to plan, function returns nil.
	return nil
}

// imageToPaletted converts image.Image to image.Paletted. Returns a pointer to result image.
func imageToPaletted(img image.Image, palette color.Palette) *image.Paletted {
	bounds := img.Bounds()

	paletted := image.NewPaletted(bounds, palette)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			paletted.Set(x, y, img.At(x, y))
		}
	}

	return paletted
}
