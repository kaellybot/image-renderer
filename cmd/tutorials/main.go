package main

import (
	"fmt"
	"image"
	"kaellybot/image-renderer/pkg/images"
	"log"
	"os"
	"time"
)

func main() {
	bounds := image.Rect(350, 200, 1000, 1010)
	duration := 10 * time.Second
	fps := 20
	outputDir := "gif_frames"
	outputGIF := "out.gif"

	os.Mkdir(outputDir, 0755)
	defer os.RemoveAll(outputDir)

	log.Println("Recording...")
	if err := images.RecordScreen(outputDir, duration, fps, bounds); err != nil {
		log.Printf("ffmpeg recording failed: %v\n", err)
		return
	}

	if err := images.GenerateGIFWithFFmpeg(outputDir, fps, outputGIF); err != nil {
		log.Printf("ffmpeg gif generation failed: %v\n", err)
		return
	}

	fmt.Println("GIF saved to", outputGIF)
}
