package main

import (
	"image"
	"image/color/palette"
	"image/gif"
	"log"
	"os"
	"sync"
	"time"

	"github.com/kbinani/screenshot"
	"golang.org/x/image/draw" // for resizing (optional)
)

func main() {
	bounds := screenshot.GetDisplayBounds(0)
	duration := 10 * time.Second
	delay := 5 // in 100ths of a second → 50ms
	frameDelay := time.Duration(delay*10) * time.Millisecond
	frameCount := int(duration / frameDelay)

	var images []*image.RGBA

	log.Printf("recording %v frames...\n", frameCount)
	beforeRecording := time.Now()
	for range frameCount {
		frameStart := time.Now()
		img, err := screenshot.CaptureRect(bounds)
		if err != nil {
			log.Fatal(err)
		}

		images = append(images, img)
		time.Sleep(time.Duration(delay*10)*time.Millisecond - time.Since(frameStart))
	}

	log.Printf("record ended in %v, preprocessing images...\n", time.Since(beforeRecording))
	beforeProcessing := time.Now()

	palettedImages := make([]*image.Paletted, len(images))
	delays := make([]int, len(images))
	var wg sync.WaitGroup
	for i, img := range images {
		wg.Add(1)
		go func(i int, img image.Image) {
			defer wg.Done()
			bounds := img.Bounds()
			paletted := image.NewPaletted(bounds, palette.Plan9)
			draw.Draw(paletted, bounds, img, image.Point{}, draw.Src)
			palettedImages[i] = paletted
			delays[i] = delay
		}(i, img)
	}
	wg.Wait()

	log.Printf("preprocessing ended in %v, creating GIF...\n", time.Since(beforeProcessing))
	beforeGIF := time.Now()
	// Write to file
	f, err := os.Create("screen_record.gif")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	err = gif.EncodeAll(f, &gif.GIF{
		Image: palettedImages,
		Delay: delays,
	})
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("GIF created in %v: screen_record.gif\n", time.Since(beforeGIF))
}
