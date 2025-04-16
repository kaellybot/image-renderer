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
	fps := 15
	outputDir := "gif_frames"
	outputWEBP := "out.webp"

	os.Mkdir(outputDir, 0755)
	defer os.RemoveAll(outputDir)

	log.Println("Recording...")
	if err := images.RecordScreen(outputDir, duration, fps, bounds); err != nil {
		log.Printf("ffmpeg recording failed: %v\n", err)
		return
	}

	if err := images.GenerateWebpWithFFmpeg(outputDir, fps, outputWEBP); err != nil {
		log.Printf("ffmpeg webp generation failed: %v\n", err)
		return
	}

	fmt.Println("WEBP saved to", outputWEBP)
}
