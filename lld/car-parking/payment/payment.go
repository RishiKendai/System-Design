package payment

import "fmt"

type PaymentType string

const (
	UPI        PaymentType = "UPI"
	CreditCard PaymentType = "CreditCard"
	DebitCard  PaymentType = "DebitCard"
	Cash       PaymentType = "Cash"
)

type Payment interface {
	Pay(amount float64) bool
}

type UPIPayment struct {
	UPIID string
}

func (p *UPIPayment) Pay(amount float64) bool {
	fmt.Println("Paying with UPI: ", p.UPIID)
	return true
}

type CreditCardPayment struct {
	CardNumber string
	ExpiryDate string
	CVV        string
}

func (p *CreditCardPayment) Pay(amount float64) bool {
	fmt.Println("Paying with Credit Card: ", p.CardNumber)
	return true
}

type DebitCardPayment struct {
	CardNumber string
	ExpiryDate string
	CVV        string
}

func (p *DebitCardPayment) Pay(amount float64) bool {
	fmt.Println("Paying with Debit Card: ", p.CardNumber)
	return true
}

type CashPayment struct {
	Cash float64
}

func (p *CashPayment) Pay(amount float64) bool {
	fmt.Println("Paying with Cash: ", p.Cash)
	return true
}

func PaymentStrategy(paymentType PaymentType) Payment {
	switch paymentType {
	case UPI:
		return &UPIPayment{UPIID: "1234567890"}
	case CreditCard:
		return &CreditCardPayment{CardNumber: "1234567890", ExpiryDate: "12/2026", CVV: "123"}
	case DebitCard:
		return &DebitCardPayment{CardNumber: "1234567890", ExpiryDate: "12/2026", CVV: "123"}
	case Cash:
		return &CashPayment{Cash: 100}
	default:
		return invalidPayment{}
	}
}
