package ticket

import (
	"fmt"
	"time"

	"github.com/lld/car-parking/payment"
	"github.com/lld/car-parking/pricing"
	"github.com/lld/car-parking/slot"
)

type Ticket struct {
	FloorNumber     int
	SlotId          int
	Vehicle         *slot.Vehicle
	EntryTime       time.Time
	ExitTime        *time.Time
	PricingStrategy pricing.PricingStrategy
	Price           float64
	PaymentType     payment.PaymentType
}

func NewTicket(floorNumber int, slotId int, vehicle *slot.Vehicle, entryTime time.Time, strategy pricing.PricingStrategy, paymentType payment.PaymentType) *Ticket {
	return &Ticket{
		FloorNumber:     floorNumber,
		SlotId:          slotId,
		Vehicle:         vehicle,
		EntryTime:       entryTime,
		ExitTime:        nil,
		PricingStrategy: strategy,
		Price:           0,
		PaymentType:     paymentType,
	}
}

func (t *Ticket) CalculatePrice() float64 {
	if t.PricingStrategy == nil {
		return 0
	}
	exit := time.Time{}
	if t.ExitTime != nil {
		exit = *t.ExitTime
	}
	return t.PricingStrategy.Compute(t.EntryTime, exit)
}

func (t *Ticket) Print() {
	fmt.Println("--------------------------------")
	fmt.Println("Ticket:")
	fmt.Println("Floor Number: ", t.FloorNumber)
	fmt.Println("Slot ID: ", t.SlotId)
	fmt.Println("Vehicle: ", t.Vehicle.License)
	fmt.Println("Entry Time: ", t.EntryTime.Format("Jan 2, 2006 03:04 PM"))
	fmt.Println("--------------------------------")
}
