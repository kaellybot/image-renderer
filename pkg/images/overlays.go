package images

import (
	"image"
	"image/draw"
)

// Overlay a child image onto a parent image while keeping transparency
func OverlayImages(parent image.Image, child image.Image) *image.RGBA {
	// Create an RGBA canvas with the same size as the parent image
	parentBounds := parent.Bounds()
	canvas := image.NewRGBA(parentBounds)

	// Draw the parent image onto the canvas
	draw.Draw(canvas, parentBounds, parent, image.Point{}, draw.Src)

	// Get the bounds of the child image
	childBounds := child.Bounds()
	childWidth := childBounds.Dx()
	childHeight := childBounds.Dy()

	// Calculate the position to center the child image on the parent
	offsetX := (parentBounds.Dx() - childWidth) / 2
	offsetY := (parentBounds.Dy() - childHeight) / 2

	// Draw the child image onto the canvas
	draw.Draw(canvas, image.Rect(offsetX, offsetY, offsetX+childWidth, offsetY+childHeight), child, image.Point{}, draw.Over)

	return canvas
}
