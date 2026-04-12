package payment

import "fmt"

type invalidPayment struct{}

func (invalidPayment) Pay(_ float64) bool {
	fmt.Println("Invalid payment type")
	return false
}
