package slot

type SpotType int

const (
	BikeSpot SpotType = iota + 1
	CarSpot
	TruckSpot
)

type SlotKind int

const (
	NormalBay SlotKind = iota
	EVBay
	SolarBay
)

type Vehicle struct {
	License string
	Type    SpotType
	IsEV    bool
}

type ParkingSlot struct {
	Id           int
	FloorNumber  int
	Type         SpotType
	Kind         SlotKind
	OccupiedType SpotType
	Vehicle      *Vehicle
	IsAvailable  bool
}
