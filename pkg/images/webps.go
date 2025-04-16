package images

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func GenerateWebpWithFFmpeg(dir string, fps int, out string) error {
	input := filepath.Join(dir, "frame_%04d.png")
	cmd := exec.Command("ffmpeg",
		"-y",
		"-framerate", fmt.Sprint(fps),
		"-i", input,
		"-loop", "0", // loop forever
		"-vf", "scale=480:-1",
		"-c:v", "libwebp",
		"-preset", "picture", // or "default", "drawing", "photo", etc.
		"-an",            // no audio
		"-lossless", "0", // set to "1" if you want lossless
		"-qscale", "60", // adjust quality: lower = better quality
		out,
	)

	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("WebP generation failed: %w", err)
	}

	return nil
}
