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
	"slices"
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

type equipImageSlot struct {
	image    image.Image
	itemType int32
}

func main() {
	var id int32 = 1 // Set ID
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
		writeOnDisk(id, slotGrid)
	} else {
		customSlotGrid := placeCustomGrid(apiClient, resp)
		writeOnDisk(id, customSlotGrid)
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

func writeOnDisk(setID int32, img image.Image) error {
	path := constants.SetImageFolderPath
	filename := fmt.Sprintf("%v.webp", setID)
	out, err := os.Create(filepath.Join(path, filename))
	if err != nil {
		return err
	}
	defer out.Close()
	return webp.Encode(out, img, &webp.Options{Lossless: true})
}

func placeClassicGrid(apiClient *dodugo.APIClient, equipmentSet *dodugo.EquipmentSet) image.Image {
	equipmentsImage := make([]equipImageSlot, 0)
	maxEquipmentTypes := initMaxEquipmentTypes()

	leftFilled, errleftFilled := imaging.Open("pkg/resources/left-filled-slot.png")
	if errleftFilled != nil {
		log.Fatal(errleftFilled)
	}

	rightFilled, errrightFilled := imaging.Open("pkg/resources/right-filled-slot.png")
	if errrightFilled != nil {
		log.Fatal(errrightFilled)
	}

	filled, errFilled := imaging.Open("pkg/resources/filled-slot.png")
	if errFilled != nil {
		log.Fatal(errFilled)
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

		equipmentsImage = append(equipmentsImage, equipImageSlot{
			image:    imageItem,
			itemType: *respEquip.GetType().Id,
		})
	}

	var maxColToLeft, maxColToRight int = 1, 1
	mapTypes := make(map[int32]int)
	for _, equipment := range equipmentsImage {
		itemType := equipment.itemType
		if slices.Contains([]int32{73, 39, 93, 42, 111, 99, 163, 52, 125, 80, 65}, itemType) {
			itemType = 73
		}
		mapTypes[itemType] = mapTypes[itemType] + 1
	}

	for mapType, nb := range mapTypes {
		point := constants.GetSetPoints()[mapType][0]
		if point.X == constants.SetItemMarginPx {
			if nb > maxColToLeft {
				maxColToLeft = nb
			}
		} else {
			if mapType == 17 {
				nb = nb/2 + nb%2
			}
			if nb > maxColToRight {
				maxColToRight = nb
			}
		}
	}

	var globalShift int = (maxColToLeft - 1) * (constants.SetItemMarginPx + constants.SetItemSizePx)
	gridPath := fmt.Sprintf("pkg/resources/classic-grid-%v-%v.png", maxColToLeft, maxColToRight)
	slotGrid, errSlotGrid := imaging.Open(gridPath)
	if errSlotGrid != nil {
		log.Fatal(errSlotGrid)
	}

	for _, equipment := range equipmentsImage {
		index := 0
		for _, maxEquipmentType := range maxEquipmentTypes {
			if strings.Contains(maxEquipmentType.equipmentType, fmt.Sprintf("'%d'", equipment.itemType)) {
				index += maxEquipmentType.nbCurrentEquipped

				points, pointFound := constants.GetSetPoints()[equipment.itemType]
				if !pointFound {
					log.Fatalf("item type have not equivalent point: %v",
						equipment.itemType)
				}

				// Overlay image on filled slot
				var shift int
				position := points[index%maxEquipmentType.nbCanEquip]
				position = image.Pt(position.X+globalShift, position.Y)
				if index >= maxEquipmentType.nbCanEquip {
					equipment.image = images.OverlayImages(filled, equipment.image)

					if points[0].X == constants.SetItemMarginPx {
						shift = -int(math.Floor(float64(maxEquipmentType.nbCurrentEquipped) / float64(maxEquipmentType.nbCanEquip)))
					} else {
						shift = int(math.Floor(float64(maxEquipmentType.nbCurrentEquipped) / float64(maxEquipmentType.nbCanEquip)))
					}

				} else if points[0].X == constants.SetItemMarginPx {
					equipment.image = images.OverlayImages(leftFilled, equipment.image)
				} else {
					equipment.image = images.OverlayImages(rightFilled, equipment.image)
				}

				position = image.Pt(position.X+shift*(constants.SetItemMarginPx+constants.SetItemSizePx), position.Y)
				// Overlay filled slot to grid
				slotGrid = imaging.Overlay(slotGrid, equipment.image, position, 1)

				maxEquipmentType.nbCurrentEquipped++
				break
			}
		}
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
