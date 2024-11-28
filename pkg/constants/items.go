package constants

type EquipmentType int

const (
	Amulet EquipmentType = iota + 1
	Belt
	Boots
	Cape
	Hat
	Pet
	Ring
	Shield
	Weapon

	EquipmentBasePath = "pkg/resources"
)

type Item struct {
	EquipmentType
	HorizontalAlign
	Icon      string
	Positions []int
}

func GetEquipments() []Item {
	return []Item{
		{
			EquipmentType:   Amulet,
			HorizontalAlign: Right,
			Positions:       []int{1},
			Icon:            "amulet.png",
		},
		{
			EquipmentType:   Belt,
			HorizontalAlign: Right,
			Positions:       []int{4},
			Icon:            "belt.png",
		},
		{
			EquipmentType:   Boots,
			HorizontalAlign: Right,
			Positions:       []int{5},
			Icon:            "boots.png",
		},
		{
			EquipmentType:   Cape,
			HorizontalAlign: Left,
			Positions:       []int{2},
			Icon:            "cape.png",
		},
		{
			EquipmentType:   Hat,
			HorizontalAlign: Left,
			Positions:       []int{1},
			Icon:            "hat.png",
		},
		{
			EquipmentType:   Pet,
			HorizontalAlign: Left,
			Positions:       []int{5},
			Icon:            "pet.png",
		},
		{
			EquipmentType:   Ring,
			HorizontalAlign: Right,
			Positions:       []int{2, 3},
			Icon:            "ring.png",
		},
		{
			EquipmentType:   Shield,
			HorizontalAlign: Left,
			Positions:       []int{4},
			Icon:            "shield.png",
		},
		{
			EquipmentType:   Weapon,
			HorizontalAlign: Left,
			Positions:       []int{3},
			Icon:            "weapon.png",
		},
	}
}
