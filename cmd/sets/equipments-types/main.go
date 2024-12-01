package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/dofusdude/dodugo"
)

func main() {
	config := dodugo.NewConfiguration()
	apiClient := dodugo.NewAPIClient(config)

	respSet, rSet, errSet := apiClient.SetsAPI.
		GetAllSetsList(context.Background(), "fr", "dofus3beta").Execute()
	if errSet != nil && (rSet == nil || rSet.StatusCode != http.StatusNotFound) {
		log.Fatalf("failed to fetch sets: %v", errSet)
	}
	defer rSet.Body.Close()

	respItems, rItems, errItems := apiClient.EquipmentAPI.
		GetAllItemsEquipmentList(context.Background(), "fr", "dofus3beta").Execute()
	if errItems != nil && (rItems == nil || rItems.StatusCode != http.StatusNotFound) {
		log.Fatalf("failed to fetch equipments: %v", errItems)
	}
	defer rItems.Body.Close()

	respCosms, rCosms, errCosms := apiClient.CosmeticsAPI.
		GetAllCosmeticsList(context.Background(), "fr", "dofus3beta").Execute()
	if errCosms != nil && (rCosms == nil || rCosms.StatusCode != http.StatusNotFound) {
		log.Fatalf("failed to fetch equipments: %v", errCosms)
	}
	defer rCosms.Body.Close()

	equipmentToType := make(map[int32]int32)
	for _, equipment := range respItems.GetItems() {
		equipmentToType[equipment.GetAnkamaId()] = equipment.Type.GetId()
	}

	cosmeticToType := make(map[int32]int32)
	for _, cosmetic := range respCosms.GetItems() {
		cosmeticToType[cosmetic.GetAnkamaId()] = cosmetic.Type.GetId()
	}

	maxEquipmentType := make(map[string]int)
	maxCosmeticType := make(map[string]int)
	for _, set := range respSet.GetSets() {
		localEquipmentType := make(map[string]int)
		localCosmeticType := make(map[string]int)

		for _, itemID := range set.GetEquipmentIds() {
			if set.GetContainsCosmeticsOnly() {
				itemType, found := cosmeticToType[itemID]
				if !found {
					log.Fatalf("Cannot find cosmetic item type for item ID %v", itemID)
				}
				itemEnum := MapTypeToEnum(itemType)

				localCosmeticType[itemEnum] = localCosmeticType[itemEnum] + 1
			} else {
				itemType, found := equipmentToType[itemID]
				if !found {
					log.Fatalf("Cannot find equipment item type for item ID %v", itemID)
				}
				itemEnum := MapTypeToEnum(itemType)

				localEquipmentType[itemEnum] = localEquipmentType[itemEnum] + 1
			}
		}

		for equipType, number := range localEquipmentType {
			if number > maxEquipmentType[equipType] {
				maxEquipmentType[equipType] = number
			}
		}

		for cosmType, number := range localCosmeticType {
			if number > maxCosmeticType[cosmType] {
				maxCosmeticType[cosmType] = number
			}
		}
	}

	fmt.Println("Max for equipments:")
	for equipType, number := range maxEquipmentType {
		fmt.Printf("%v=%v\n", equipType, number)
	}

	fmt.Println("\nMax for cosmetics:")
	for cosmType, number := range maxCosmeticType {
		fmt.Printf("%v=%v\n", cosmType, number)
	}
}

func MapTypeToEnum(itemType int32) string {
	enums := map[int32]string{
		5:   "living object",
		6:   "ceremonial hat",
		16:  "costume",
		18:  "ceremonial cape",
		19:  "ceremonial shield",
		24:  "ceremonial pet",
		34:  "ceremonial weapon",
		59:  "ceremonial petmount",
		97:  "misc",
		112: "seemyhol harnachment",
		130: "dragoturkey harnachment",
		166: "rhineetle harnachment",
		217: "shoulder",
		218: "wings",
		27:  "hat",
		43:  "cape",
		73:  "weapon",
		39:  "weapon",
		93:  "weapon",
		42:  "weapon",
		111: "weapon",
		99:  "weapon",
		163: "weapon",
		52:  "weapon",
		125: "weapon",
		80:  "weapon",
		65:  "weapon",
		87:  "shield",
		1:   "pet",
		33:  "amulet",
		17:  "ring",
		58:  "belt",
		45:  "boots",
	}

	if _, found := enums[itemType]; !found {
		return fmt.Sprintf("%v", itemType)
	}
	return enums[itemType]
}
