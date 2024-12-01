package constants

import (
	"image"
)

type HorizontalAlign int

const (
	Left HorizontalAlign = iota + 1
	Right
	SetImageFolderPath = "outputs"
	SetItemMarginPx    = 10
	SetItemSizePx      = 128
	setFirstCell       = SetItemMarginPx
	setSecondCell      = setFirstCell + SetItemSizePx + SetItemMarginPx
	setThirdCell       = setSecondCell + SetItemSizePx + SetItemMarginPx
	setFourthCell      = setThirdCell + SetItemSizePx + SetItemMarginPx
	setFifthCell       = setFourthCell + SetItemSizePx + SetItemMarginPx
)

//nolint:exhaustive // No other types needed.
func GetSetPoints() map[int32][]image.Point {
	weaponPoints := []image.Point{image.Pt(setFirstCell, setThirdCell)}
	return map[int32][]image.Point{

		27:  {image.Pt(setFirstCell, setFirstCell)}, // Hat
		43:  {image.Pt(setFirstCell, setSecondCell)}, // Cloak
		73:  weaponPoints,
		39:  weaponPoints,
		93:  weaponPoints,
		42:  weaponPoints,
		111: weaponPoints,
		99:  weaponPoints,
		163: weaponPoints,
		52:  weaponPoints,
		125: weaponPoints,
		80:  weaponPoints,
		65:  weaponPoints,
		87:  {image.Pt(setFirstCell, setFourthCell)}, // Shield
		1:   {image.Pt(setFirstCell, setFifthCell)}, // Pet

		33: {image.Pt(setFifthCell, setFirstCell)}, // Amulet
		17: {
			image.Pt(setFifthCell, setSecondCell), // Ring1
			image.Pt(setFifthCell, setThirdCell), // Ring2
		},
		58: {image.Pt(setFifthCell, setFourthCell)}, // Belt
		45: {image.Pt(setFifthCell, setFifthCell)}, // Boots
	}
}
