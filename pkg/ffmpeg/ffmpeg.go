package ffmpeg

import (
	"fmt"
	"image"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

func RecordScreen(dir string, duration time.Duration, fps int, bounds image.Rectangle) error {
	cmd := exec.Command("ffmpeg",
		"-f", "gdigrab",
		"-framerate", fmt.Sprintf("%v", fps),
		"-video_size", fmt.Sprintf("%dx%d", bounds.Dx(), bounds.Dy()),
		"-offset_x", fmt.Sprintf("%d", bounds.Min.X),
		"-offset_y", fmt.Sprintf("%d", bounds.Min.Y),
		"-i", "desktop",
		"-t", fmt.Sprintf("%.0f", duration.Seconds()),
		"-f", "image2", // ensure ffmpeg treats this as image sequence
		fmt.Sprintf("%s/frame_%%04d.png", dir),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
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
