package parking

import (
	spotheap "github.com/lld/car-parking/heap"
	"github.com/lld/car-parking/pricing"
	"github.com/lld/car-parking/slot"
)

type SpotType = slot.SpotType
type SlotKind = slot.SlotKind
type Vehicle = slot.Vehicle
type ParkingSlot = slot.ParkingSlot

const (
	BikeSpot  SpotType = slot.BikeSpot
	CarSpot   SpotType = slot.CarSpot
	TruckSpot SpotType = slot.TruckSpot
	NormalBay SlotKind = slot.NormalBay
	EVBay     SlotKind = slot.EVBay
	SolarBay  SlotKind = slot.SolarBay
)

type ParkingFloors struct {
	FloorNumber int
	Slots       []*ParkingSlot
}

type ParkingLot struct {
	Floors         []*ParkingFloors
	AvailableSlots map[SpotType]map[SlotKind]*spotheap.SpotMinHeap
}

type FloorConfig struct {
	Bikes       int
	BikesEV     int
	BikesSolar  int
	Cars        int
	CarsEV      int
	CarsSolar   int
	Trucks      int
	TrucksEV    int
	TrucksSolar int
}

func (c FloorConfig) totalSlots() int {
	return c.Bikes + c.BikesEV + c.BikesSolar +
		c.Cars + c.CarsEV + c.CarsSolar +
		c.Trucks + c.TrucksEV + c.TrucksSolar
}

// Pre-built strategies for callers (e.g. main) that want cp.Hourly-style names.
var (
	Hourly  = pricing.NewHourly()
	Flat    = pricing.NewFlat()
	Dynamic = pricing.NewDynamic()
)
