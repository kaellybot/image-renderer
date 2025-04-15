package main

import (
	"fmt"
	"image"
	"image/png"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kbinani/screenshot"
	// for resizing (optional)
)

func main() {
	bounds := image.Rect(350, 200, 1000, 1010)
	duration := 10 * time.Second
	fps := 20

	delay := 100 / fps
	frameDelay := time.Duration(delay*10) * time.Millisecond
	frameCount := int(duration / frameDelay)

	outputDir := "gif_frames"
	outputGIF := "out.gif"

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

	log.Printf("record ended in %v\n", time.Since(beforeRecording))
	if err := SaveFrames(images, outputDir); err != nil {
		log.Fatalf("saving frames: %v", err)
	}
	defer os.RemoveAll(outputDir)

	if err := GenerateGIFWithFFmpeg(outputDir, fps, outputGIF); err != nil {
		log.Fatalf("ffmpeg failed: %v", err)
	}

	fmt.Println("GIF saved to", outputGIF)
}

func SaveFrames(frames []*image.RGBA, dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	for i, img := range frames {
		filename := filepath.Join(dir, fmt.Sprintf("frame_%04d.png", i))
		f, err := os.Create(filename)
		if err != nil {
			return err
		}
		if err := png.Encode(f, img); err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

func GenerateGIFWithFFmpeg(dir string, fps int, out string) error {
	input := filepath.Join(dir, "frame_%04d.png")
	palette := filepath.Join(dir, "palette.png")

	// Step 1: generate palette
	cmd1 := exec.Command("ffmpeg", "-y", "-framerate", fmt.Sprint(fps), "-i", input,
		"-vf", "palettegen", palette)
	cmd1.Stderr = os.Stderr
	cmd1.Stdout = os.Stdout
	if err := cmd1.Run(); err != nil {
		return fmt.Errorf("palettegen failed: %w", err)
	}

	// Step 2: generate GIF using palette
	cmd2 := exec.Command("ffmpeg", "-y", "-framerate", fmt.Sprint(fps), "-i", input,
		"-i", palette,
		"-filter_complex", fmt.Sprintf("fps=%d[x];[x][1:v]paletteuse", fps),
		out)
	cmd2.Stderr = os.Stderr
	cmd2.Stdout = os.Stdout
	if err := cmd2.Run(); err != nil {
		return fmt.Errorf("GIF generation failed: %w", err)
	}
	return nil
}
