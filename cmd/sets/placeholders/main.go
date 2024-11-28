package main

import (
	"fmt"
	"image"
	"kaellybot/image-renderer/pkg/constants"
	"kaellybot/image-renderer/pkg/images"
	"log"
	"os"
	"path"

	"github.com/disintegration/imaging"
)

func main() {
	basePath := "./outputs"
	if len(os.Args) >= 2 {
		basePath = os.Args[1]
	}

	emptySlots := make(map[constants.HorizontalAlign]image.Image)
	for _, emptySlot := range constants.GetEmptySlots() {
		icon, err := imaging.Open(emptySlot.Icon)
		if err != nil {
			log.Fatalf("failed to open %v icon: %v", emptySlot.Icon, err)
		}
		emptySlots[emptySlot.HorizontalAlign] = icon
	}

	for _, equipment := range constants.GetEquipments() {
		itemPath := path.Join(constants.EquipmentBasePath, equipment.Icon)
		icon, err := imaging.Open(itemPath)
		if err != nil {
			log.Fatalf("failed to open %v icon: %v", itemPath, err)
		}

		resizedIcon := images.CoverResize(icon, constants.SlotCoverWidth, constants.SlotCoverHeight)
		result := images.OverlayImages(emptySlots[equipment.HorizontalAlign], resizedIcon)

		dest := path.Join(basePath, fmt.Sprintf("placeholder_%v", equipment.Icon))
		outFile, err := os.Create(dest)
		if err != nil {
			log.Fatalf("failed to create output file: %v", err)
		}
		defer outFile.Close()

		err = imaging.Encode(outFile, result, imaging.PNG)
		if err != nil {
			log.Fatalf("failed to encode %v icon result: %v", equipment.Icon, err)
		}

		log.Printf("Image successfully written to %v\n", dest)
	}
}
