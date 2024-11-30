package main

import (
	"context"
	"fmt"
	"image"
	"kaellybot/image-renderer/pkg/constants"
	"kaellybot/image-renderer/pkg/images"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/dofusdude/dodugo"
	"github.com/spf13/viper"
)

func main() {
	var id int32 = 1 // Set ID
	if len(os.Args) >= 2 {
		newId, err := strconv.Atoi(os.Args[1])
		if err == nil {
			id = int32(newId)
		}
	}

	config := dodugo.NewConfiguration()
	apiClient := dodugo.NewAPIClient(config)

	slotGrid, errSlotGrid := imaging.Open("pkg/resources/classic-grid.png")
	if errSlotGrid != nil {
		log.Fatal(errSlotGrid)
	}
	leftFilled, errleftFilled := imaging.Open("outputs/left-filled-slot.png")
	if errleftFilled != nil {
		log.Fatal(errleftFilled)
	}
	rightFilled, errrightFilled := imaging.Open("outputs/right-filled-slot.png")
	if errrightFilled != nil {
		log.Fatal(errrightFilled)
	}

	resp, r, err := apiClient.SetsAPI.
		GetSetsSingle(context.Background(), "fr", id, "dofus3beta").Execute()
	if err != nil && (r == nil || r.StatusCode != http.StatusNotFound) {
		log.Fatalf("failed to fetch set %v: %v", id, err)
	}
	defer r.Body.Close()

	var ringNumber int
	for _, itemId := range resp.EquipmentIds {
		respEquip, r, err := apiClient.EquipmentAPI.
			GetItemsEquipmentSingle(context.Background(), "fr", int32(itemId), "dofus3beta").Execute()
		if err != nil && (r == nil || r.StatusCode != http.StatusNotFound) {
			log.Fatalf("failed to fetch equipment %v: %v", itemId, err)
		}
		defer r.Body.Close()
		imageItem := getImageFromItem(context.Background(), respEquip)

		// Resize sd image to HD
		imageItem = images.CoverResize(imageItem, constants.SlotCoverWidth, constants.SlotCoverHeight)

		index := 0
		if *respEquip.GetType().Id == int32(17) {
			index += ringNumber
			ringNumber++
		}
		points, pointFound := constants.GetSetPoints()[*respEquip.GetType().Id]
		if !pointFound {
			log.Fatalf("item %v type have not equivalent point: %v",
				respEquip.GetAnkamaId(), *respEquip.GetType().Id)
		}
		// Overlay image on filled slot
		if points[index].X == constants.SetItemMarginPx {
			imageItem = images.OverlayImages(leftFilled, imageItem)
		} else {
			imageItem = images.OverlayImages(rightFilled, imageItem)
		}

		// Overlay filled slot to grid
		slotGrid = imaging.Overlay(slotGrid, imageItem, points[index], 1)
	}

	writeOnDisk(id, slotGrid)
}

func getImageFromURL(ctx context.Context, rawURL string,
) (image.Image, error) {
	parsedURL, err := url.ParseRequestURI(rawURL)
	if err != nil {
		return nil, err
	}

	req, errReq := http.NewRequestWithContext(ctx, http.MethodGet, parsedURL.String(), nil)
	if errReq != nil {
		return nil, errReq
	}

	client := &http.Client{}
	resp, errDo := client.Do(req)
	if errDo != nil {
		return nil, errDo
	}
	defer resp.Body.Close()

	image, errDecode := imaging.Decode(resp.Body)
	if errDecode != nil {
		return nil, errDecode
	}

	return image, nil
}

func getImageFromItem(ctx context.Context, item *dodugo.Weapon) image.Image {
	if item.GetImageUrls().Sd.IsSet() {
		itemImage, errGetImg := getImageFromURL(ctx, *item.GetImageUrls().Sd.Get())
		if errGetImg != nil {
			log.Fatalf("failed to fetch image equipment: %v", errGetImg)
		}

		return itemImage
	}

	log.Fatalf("failed to fetch image equipment in sd")
	return nil
}

func writeOnDisk(setID int32, img image.Image) error {
	path := viper.GetString(constants.SetImageFolderPath)
	filename := fmt.Sprintf("%v.webp", setID)
	out, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return err
	}
	defer out.Close()
	return webp.Encode(out, img, &webp.Options{Lossless: true})
}
