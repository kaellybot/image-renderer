package main

import (
	"flag"
	"fmt"
	"kaellybot/image-renderer/pkg/automations"
	"kaellybot/image-renderer/pkg/constants"
	"kaellybot/image-renderer/pkg/images"
	"log"
	"os"
	"path/filepath"
	"slices"
	"sync"

	amqp "github.com/kaellybot/kaelly-amqp"
)

func main() {
	commandArg := flag.String("command", "", "Command name to execute")
	localeArg := flag.String("locale", "FR", "Locale (2 chars, uppercase, default 'FR')")
	flag.Parse()
	fps := 15
	locale := amqp.Language_EN
	if amqpLocale, found := amqp.Language_value[*localeArg]; found {
		locale = amqp.Language(amqpLocale)
	}

	command, found := constants.GetCommands()[*commandArg]
	if !found {
		commandNames := make([]string, 0)
		for name := range constants.GetCommands() {
			commandNames = append(commandNames, name)
		}
		slices.Sort(commandNames)
		log.Printf("command %v not found\nList available: %v\n", *commandArg, commandNames)
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
		if err := automations.RunCommandTutorial(command, locale); err != nil {
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
