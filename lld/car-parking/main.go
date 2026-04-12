package main

import (
	"fmt"
	"time"

	cp "github.com/lld/car-parking/parking"
	"github.com/lld/car-parking/payment"
)

func main() {

	parkingLot := cp.NewParkingLot()

	v1 := &cp.Vehicle{
		License: "TN 59 BQ 4257",
		Type:    cp.BikeSpot,
		IsEV:    false,
	}
	v2 := &cp.Vehicle{
		License: "PY01 CV0060",
		Type:    cp.CarSpot,
		IsEV:    true,
	}

	t1, err := parkingLot.Park(v1, cp.Hourly, payment.UPI)
	if err != nil {
		fmt.Println(err)
		return
	}
	t1.Print()

	t2, err := parkingLot.Park(v2, cp.Flat, payment.Cash)
	if err != nil {
		fmt.Println(err)
		return
	}
	t2.Print()
	time.Sleep(20 * time.Second)
	if err := parkingLot.UnPark(t1); err != nil {
		fmt.Println(err)
	}
	if err := parkingLot.UnPark(t2); err != nil {
		fmt.Println(err)
	}
}
