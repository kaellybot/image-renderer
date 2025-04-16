package images

import (
	"fmt"
	"image"
	"os"
	"os/exec"
	"time"
)

func RecordScreen(dir string, duration time.Duration, fps int, bounds image.Rectangle) error {
	cmd := exec.Command("ffmpeg",
		"-f", "gdigrab",
		"-draw_mouse", "0",
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
