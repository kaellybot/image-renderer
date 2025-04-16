package main

import (
	"fmt"
	"image"
	"kaellybot/image-renderer/pkg/automations"
	"kaellybot/image-renderer/pkg/images"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	amqp "github.com/kaellybot/kaelly-amqp"
)

func main() {
	bounds := image.Rect(350, 200, 1000, 1010)
	duration := 10 * time.Second
	fps := 15
	outputDir := filepath.Join(os.TempDir(), "kaelly_frames")
	outputWEBP := "out.webp"

	os.Mkdir(outputDir, 0755)
	defer os.RemoveAll(outputDir)

	if err := automations.SetupDiscordTutorial(); err != nil {
		log.Printf("discord setup failed: %v\n", err)
		return
	}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		log.Println("Recording...")
		if err := images.RecordScreen(outputDir, duration, fps, bounds); err != nil {
			log.Printf("ffmpeg recording failed: %v\n", err)
		}
		wg.Done()
	}()

	go func() {
		log.Println("Automating...")
		if err := automations.RunCommandTutorial("about", amqp.Language_FR); err != nil {
			log.Printf("ffmpeg recording failed: %v\n", err)
		}
		wg.Done()
	}()

	wg.Wait()

	if err := images.GenerateWebpWithFFmpeg(outputDir, fps, outputWEBP); err != nil {
		log.Printf("ffmpeg webp generation failed: %v\n", err)
		return
	}

	fmt.Println("WEBP saved to", outputWEBP)
}
