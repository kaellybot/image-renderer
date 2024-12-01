package main

import (
	"context"
	"fmt"
	"image"
	"kaellybot/image-renderer/pkg/constants"
	"kaellybot/image-renderer/pkg/images"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/chai2010/webp"
	"github.com/disintegration/imaging"
	"github.com/dofusdude/dodugo"
)

type maxEquipmentType struct {
	equipmentType     string
	nbCurrentEquipped int
	nbCanEquip        int
}

func main() {
	var id int32 = 514 // Set ID
	if len(os.Args) >= 2 {
		newId, err := strconv.Atoi(os.Args[1])
		if err == nil {
			id = int32(newId)
		}
	}

	var classicGrid bool = true
	config := dodugo.NewConfiguration()
	apiClient := dodugo.NewAPIClient(config)

	resp, r, err := apiClient.SetsAPI.
		GetSetsSingle(context.Background(), "fr", id, "dofus3beta").Execute()
	if err != nil && (r == nil || r.StatusCode != http.StatusNotFound) {
		log.Fatalf("failed to fetch set %v: %v", id, err)
	}
	defer r.Body.Close()

	// check set isCosmetic
	if resp.GetContainsCosmetics() {
		classicGrid = false
	}

	if classicGrid {
		slotGrid := placeClassicGrid(apiClient, resp)
		writeOnDisk(id, &slotGrid)
	} else {
		customSlotGrid := placeCustomGrid(apiClient, resp)
		writeOnDisk(id, &customSlotGrid)
	}
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

func writeOnDisk(setID int32, img *image.Image) error {
	path := constants.SetImageFolderPath
	filename := fmt.Sprintf("%v.webp", setID)
	out, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return err
	}
	defer out.Close()
	return webp.Encode(out, *img, &webp.Options{Lossless: true})
}

func placeClassicGrid(apiClient *dodugo.APIClient, equipmentSet *dodugo.EquipmentSet) image.Image {
	maxEquipmentTypes := initMaxEquipmentTypes()
	slotGrid, errSlotGrid := imaging.Open("pkg/resources/classic-grid.png")
	if errSlotGrid != nil {
		log.Fatal(errSlotGrid)
	}

	leftFilled, errleftFilled := imaging.Open("pkg/resources/left-filled-slot.png")
	if errleftFilled != nil {
		log.Fatal(errleftFilled)
	}

	rightFilled, errrightFilled := imaging.Open("pkg/resources/right-filled-slot.png")
	if errrightFilled != nil {
		log.Fatal(errrightFilled)
	}

	for _, itemId := range equipmentSet.GetEquipmentIds() {
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
		for _, maxEquipmentType := range maxEquipmentTypes {
			test := fmt.Sprintf("'%d'", *respEquip.GetType().Id)
			if strings.Contains(maxEquipmentType.equipmentType, test) {
				index += maxEquipmentType.nbCurrentEquipped
				maxEquipmentType.nbCurrentEquipped++
				if maxEquipmentType.nbCurrentEquipped > maxEquipmentType.nbCanEquip {
					return placeCustomGrid(apiClient, equipmentSet)
				}
				break
			}
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

	return slotGrid
}

func placeCustomGrid(apiClient *dodugo.APIClient, equipmentSet *dodugo.EquipmentSet) image.Image {
	var moduloPos int = 1
	var imageError error
	var customSlotGrid image.Image

	filled, errFilled := imaging.Open("pkg/resources/filled-slot.png")
	if errFilled != nil {
		log.Fatal(errFilled)
	}

	if len(equipmentSet.GetEquipmentIds()) < 2 {
		customSlotGrid, imageError = imaging.Open("pkg/resources/1-grid-kaelly.png")
		if imageError != nil {
			log.Fatal(imageError)
		}
	} else if len(equipmentSet.GetEquipmentIds()) < 5 {
		moduloPos = 2
		customSlotGrid, imageError = imaging.Open("pkg/resources/4-grid-kaelly.png")
		if imageError != nil {
			log.Fatal(imageError)
		}
	} else if len(equipmentSet.GetEquipmentIds()) < 10 {
		moduloPos = 3
		customSlotGrid, imageError = imaging.Open("pkg/resources/9-grid-kaelly.png")
		if imageError != nil {
			log.Fatal(imageError)
		}
	} else if len(equipmentSet.GetEquipmentIds()) < 17 {
		moduloPos = 4
		customSlotGrid, imageError = imaging.Open("pkg/resources/16-grid-kaelly.png")
		if imageError != nil {
			log.Fatal(imageError)
		}
	} else {
		moduloPos = 5
		customSlotGrid, imageError = imaging.Open("pkg/resources/25-grid-kaelly.png")
		if imageError != nil {
			log.Fatal(imageError)
		}
	}

	for i, itemId := range equipmentSet.GetEquipmentIds() {
		respEquip, r, err := apiClient.EquipmentAPI.
			GetItemsEquipmentSingle(context.Background(), "fr", int32(itemId), "dofus3beta").Execute()
		if err != nil && (r == nil || r.StatusCode != http.StatusNotFound) {
			log.Fatalf("failed to fetch equipment %v: %v", itemId, err)
		}
		defer r.Body.Close()
		imageItem := getImageFromItem(context.Background(), respEquip)

		// Resize sd image to HD
		imageItem = images.CoverResize(imageItem, constants.SlotCoverWidth, constants.SlotCoverHeight)

		// Overlay image on filled slot
		imageItem = images.OverlayImages(filled, imageItem)

		// Overlay filled slot to grid
		var posX int = constants.SetItemMarginPx + i%moduloPos*(constants.SetItemMarginPx+constants.SetItemSizePx)
		var posY int = constants.SetItemMarginPx + int(math.Floor(float64(i)/float64(moduloPos)))*(constants.SetItemMarginPx+constants.SetItemSizePx)
		position := image.Pt(posX, posY)
		customSlotGrid = imaging.Overlay(customSlotGrid, imageItem, position, 1)
	}

	return customSlotGrid
}

func initMaxEquipmentTypes() []*maxEquipmentType {
	maxEquipmentTypes := make([]*maxEquipmentType, 0)

	maxEquipmentTypes = append(maxEquipmentTypes, &maxEquipmentType{
		equipmentType:     "'27'",
		nbCurrentEquipped: 0,
		nbCanEquip:        1,
	})
	maxEquipmentTypes = append(maxEquipmentTypes, &maxEquipmentType{
		equipmentType:     "'43'",
		nbCurrentEquipped: 0,
		nbCanEquip:        1,
	})
	maxEquipmentTypes = append(maxEquipmentTypes, &maxEquipmentType{
		equipmentType:     "'73' '39' '93' '42' '111' '99' '163' '52' '125' '80' '65'",
		nbCurrentEquipped: 0,
		nbCanEquip:        1,
	})
	maxEquipmentTypes = append(maxEquipmentTypes, &maxEquipmentType{
		equipmentType:     "'87'",
		nbCurrentEquipped: 0,
		nbCanEquip:        1,
	})
	maxEquipmentTypes = append(maxEquipmentTypes, &maxEquipmentType{
		equipmentType:     "'1'",
		nbCurrentEquipped: 0,
		nbCanEquip:        1,
	})
	maxEquipmentTypes = append(maxEquipmentTypes, &maxEquipmentType{
		equipmentType:     "'33'",
		nbCurrentEquipped: 0,
		nbCanEquip:        1,
	})
	maxEquipmentTypes = append(maxEquipmentTypes, &maxEquipmentType{
		equipmentType:     "'17'",
		nbCurrentEquipped: 0,
		nbCanEquip:        2,
	})
	maxEquipmentTypes = append(maxEquipmentTypes, &maxEquipmentType{
		equipmentType:     "'58'",
		nbCurrentEquipped: 0,
		nbCanEquip:        1,
	})
	maxEquipmentTypes = append(maxEquipmentTypes, &maxEquipmentType{
		equipmentType:     "'45'",
		nbCurrentEquipped: 0,
		nbCanEquip:        1,
	})

	return maxEquipmentTypes
}
