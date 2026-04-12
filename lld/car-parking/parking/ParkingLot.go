package parking

import (
	"container/heap"
	"fmt"
	"time"

	spotheap "github.com/lld/car-parking/heap"
	"github.com/lld/car-parking/payment"
	"github.com/lld/car-parking/pricing"
	"github.com/lld/car-parking/ticket"
)

func newEmptyAvailableSlots() map[SpotType]map[SlotKind]*spotheap.SpotMinHeap {
	sizes := []SpotType{BikeSpot, CarSpot, TruckSpot}
	kinds := []SlotKind{NormalBay, EVBay, SolarBay}
	m := make(map[SpotType]map[SlotKind]*spotheap.SpotMinHeap, len(sizes))
	for _, sz := range sizes {
		byKind := make(map[SlotKind]*spotheap.SpotMinHeap, len(kinds))
		for _, k := range kinds {
			byKind[k] = &spotheap.SpotMinHeap{}
		}
		m[sz] = byKind
	}
	return m
}

func buildFloorSlots(floorNum int, c FloorConfig, pl *ParkingLot) *ParkingFloors {
	total := c.totalSlots()
	floor := &ParkingFloors{
		FloorNumber: floorNum,
		Slots:       make([]*ParkingSlot, 0, total),
	}
	id := 1
	type plan struct {
		spotType SpotType
		kind     SlotKind
		n        int
	}
	plans := []plan{
		{BikeSpot, NormalBay, c.Bikes},
		{BikeSpot, EVBay, c.BikesEV},
		{BikeSpot, SolarBay, c.BikesSolar},
		{CarSpot, NormalBay, c.Cars},
		{CarSpot, EVBay, c.CarsEV},
		{CarSpot, SolarBay, c.CarsSolar},
		{TruckSpot, NormalBay, c.Trucks},
		{TruckSpot, EVBay, c.TrucksEV},
		{TruckSpot, SolarBay, c.TrucksSolar},
	}
	for _, p := range plans {
		for range p.n {
			slot := &ParkingSlot{Id: id, FloorNumber: floorNum, Type: p.spotType, Kind: p.kind, IsAvailable: true, Vehicle: nil}
			id++
			floor.Slots = append(floor.Slots, slot)
			heap.Push(pl.AvailableSlots[p.spotType][p.kind], slot)
		}
	}
	return floor
}

func initializeFloors(pl *ParkingLot) {
	config := []FloorConfig{
		{Bikes: 2, BikesEV: 2, Cars: 2, CarsEV: 2, Trucks: 0, TrucksEV: 0},
		{Bikes: 2, BikesEV: 1, Cars: 3, CarsEV: 2, Trucks: 1, TrucksEV: 0},
		{Bikes: 3, BikesEV: 2, Cars: 2, CarsEV: 2, Trucks: 1, TrucksEV: 1},
		{Bikes: 4, BikesEV: 2, Cars: 3, CarsEV: 3, Trucks: 1, TrucksEV: 1},
		{Bikes: 3, BikesEV: 1, Cars: 3, CarsEV: 2, Trucks: 1, TrucksEV: 0},
	}
	for i, c := range config {
		pl.Floors = append(pl.Floors, buildFloorSlots(i+1, c, pl))
	}
}

func NewParkingLot() *ParkingLot {
	pl := &ParkingLot{
		AvailableSlots: newEmptyAvailableSlots(),
		Floors:         make([]*ParkingFloors, 0, 10),
	}

	initializeFloors(pl)
	return pl
}

func (pl *ParkingLot) nextFloorNumber() int {
	maxN := 0
	for _, f := range pl.Floors {
		if f.FloorNumber > maxN {
			maxN = f.FloorNumber
		}
	}
	return maxN + 1
}

// AddFloor appends one new deck with the next floor number (max existing + 1).
func (pl *ParkingLot) AddFloor(config FloorConfig) error {
	if pl == nil {
		return fmt.Errorf("nil parking lot")
	}
	n := pl.nextFloorNumber()
	pl.Floors = append(pl.Floors, buildFloorSlots(n, config, pl))
	return nil
}

func getPriority(vType SpotType) []SpotType {
	switch vType {
	case BikeSpot:
		return []SpotType{BikeSpot, CarSpot, TruckSpot}
	case CarSpot:
		return []SpotType{CarSpot, TruckSpot}
	case TruckSpot:
		return []SpotType{TruckSpot}
	}
	return nil
}

func slotKindsForVehicle(v *Vehicle) []SlotKind {
	if v.IsEV {
		return []SlotKind{EVBay, SolarBay, NormalBay}
	}
	return []SlotKind{NormalBay}
}

func (pl *ParkingLot) takeSlot(spotType SpotType, kinds []SlotKind) *ParkingSlot {
	for _, k := range kinds {
		h := pl.AvailableSlots[spotType][k]
		if h.Len() > 0 {
			return heap.Pop(h).(*ParkingSlot)
		}
	}
	return nil
}

func (pl *ParkingLot) floorByNumber(floorNumber int) *ParkingFloors {
	for _, f := range pl.Floors {
		if f.FloorNumber == floorNumber {
			return f
		}
	}
	return nil
}

func slotByID(floor *ParkingFloors, slotID int) *ParkingSlot {
	for _, s := range floor.Slots {
		if s.Id == slotID {
			return s
		}
	}
	return nil
}

func (pl *ParkingLot) Park(v *Vehicle, strategy pricing.PricingStrategy, paymentType payment.PaymentType) (*ticket.Ticket, error) {
	if pl == nil {
		return nil, fmt.Errorf("nil parking lot")
	}
	if v == nil {
		return nil, fmt.Errorf("nil vehicle")
	}
	if strategy == nil {
		return nil, fmt.Errorf("nil pricing strategy")
	}
	priorities := getPriority(v.Type)
	kindOrder := slotKindsForVehicle(v)

	for _, spotType := range priorities {
		spot := pl.takeSlot(spotType, kindOrder)

		if spot != nil {
			spot.Vehicle = v
			spot.IsAvailable = false
			spot.OccupiedType = v.Type
			t := ticket.NewTicket(spot.FloorNumber, spot.Id, v, time.Now(), strategy, paymentType)
			return t, nil
		}
	}
	return nil, fmt.Errorf("❌ Spot not available.")

}

func formatTime(t time.Time) string {
	return t.Format("Jan 2, 2006 03:04 PM")
}

func (pl *ParkingLot) UnPark(tk *ticket.Ticket) error {
	if pl == nil {
		return fmt.Errorf("nil parking lot")
	}
	if tk == nil {
		return fmt.Errorf("nil ticket")
	}
	floor := pl.floorByNumber(tk.FloorNumber)
	if floor == nil {
		return fmt.Errorf("floor %d not found", tk.FloorNumber)
	}
	spot := slotByID(floor, tk.SlotId)
	if spot == nil {
		return fmt.Errorf("slot id %d not found on floor %d", tk.SlotId, tk.FloorNumber)
	}
	if spot.Vehicle == nil {
		return fmt.Errorf("slot %d on floor %d is already empty", tk.SlotId, tk.FloorNumber)
	}
	if tk.Vehicle != nil && spot.Vehicle != tk.Vehicle {
		return fmt.Errorf("ticket does not match vehicle in slot")
	}

	exitTime := time.Now()
	tk.ExitTime = &exitTime
	tk.Price = tk.CalculatePrice()

	fmt.Println("---------------------")
	fmt.Println(" RK Parking Bay")
	fmt.Println("---------------------")
	fmt.Println("Vehicle: ", tk.Vehicle.License)
	fmt.Println("Entry Time: ", formatTime(tk.EntryTime))
	fmt.Println("Exit Time: ", formatTime(*tk.ExitTime))
	fmt.Println("Price: ", tk.Price)

	fmt.Println("Paying for the parking...")
	payment := payment.PaymentStrategy(tk.PaymentType)
	if !payment.Pay(tk.Price) {
		return fmt.Errorf("payment failed")
	}
	fmt.Println("Payment successful")
	spot.Vehicle = nil
	spot.IsAvailable = true
	spot.OccupiedType = 0
	heap.Push(pl.AvailableSlots[spot.Type][spot.Kind], spot)

	fmt.Println("---------------------")

	return nil
}
