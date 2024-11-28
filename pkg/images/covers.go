package images

import (
	"image"

	"github.com/disintegration/imaging"
)

// Resize and crop the image to cover the given dimensions while keeping it centered
func CoverResize(src image.Image, targetWidth, targetHeight int) image.Image {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	// Calculate aspect ratios
	targetAspect := float64(targetWidth) / float64(targetHeight)
	srcAspect := float64(srcWidth) / float64(srcHeight)

	// Determine cropping dimensions
	var cropWidth, cropHeight int
	if srcAspect > targetAspect {
		// Source is wider: crop width
		cropHeight = srcHeight
		cropWidth = int(float64(cropHeight) * targetAspect)
	} else {
		// Source is taller: crop height
		cropWidth = srcWidth
		cropHeight = int(float64(cropWidth) / targetAspect)
	}

	// Calculate cropping rectangle (centered)
	cropX := (srcWidth - cropWidth) / 2
	cropY := (srcHeight - cropHeight) / 2
	cropRect := image.Rect(cropX, cropY, cropX+cropWidth, cropY+cropHeight)

	// Crop and resize
	cropped := imaging.Crop(src, cropRect)
	resized := imaging.Resize(cropped, targetWidth, targetHeight, imaging.Lanczos)

	return resized
}
