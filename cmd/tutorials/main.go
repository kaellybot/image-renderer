package main

import (
	"fmt"
	"kaellybot/image-renderer/pkg/automations"
	"kaellybot/image-renderer/pkg/constants"
	"kaellybot/image-renderer/pkg/images"
	"log"
	"os"
	"path/filepath"
	"sync"

	amqp "github.com/kaellybot/kaelly-amqp"
)

func main() {
	commandName := "job_get"
	locale := amqp.Language_FR
	fps := 15

	command, found := constants.GetCommands()[commandName]
	if !found {
		log.Printf("command %v not found\n", commandName)
		return
	}

	outputDir := filepath.Join(os.TempDir(), "kaelly_frames")
	outputWEBP := fmt.Sprintf("%v_%v.webp", command.Name, locale)

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
		if err := images.RecordScreen(outputDir, command.Duration, fps, command.Bounds); err != nil {
			log.Printf("ffmpeg recording failed: %v\n", err)
		}
		wg.Done()
	}()

	go func() {
		log.Println("Automating...")
		if err := automations.RunCommandTutorial(command, amqp.Language_FR); err != nil {
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
