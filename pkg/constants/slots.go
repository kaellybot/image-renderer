package constants

const (
	SlotWidth       = 128
	SlotHeight      = 128
	SlotBorder      = 5
	SlotCoverWidth  = SlotWidth - 2*SlotBorder
	SlotCoverHeight = SlotHeight - 2*SlotBorder
)

type Slot struct {
	HorizontalAlign
	Icon string
}

func GetEmptySlots() []Slot {
	return []Slot{
		{
			HorizontalAlign: Left,
			Icon:            "outputs/left-empty-slot.png",
		},
		{
			HorizontalAlign: Right,
			Icon:            "outputs/right-empty-slot.png",
		},
	}
}
