package main

import (
	"fmt"
	"reflect"
)

// go has duck typing and polymorphism can be achieved by using interfaces
type Payer interface {
	Pay(int) error
}


type Wallet struct {
	Id int
	Cash int
}

type Cash struct {
	Amount int
}

type UnknownPayment struct {
	Amount int
}

type Card struct {
	Cash int
	Cardholder string
	IsValid bool
}

func (w *Wallet) Pay(amount int) error {
	if w.Cash < amount {
		return fmt.Errorf("insufficient amount in wallet: left: %d, required: %d", w.Cash, amount)
	}
	w.Cash -= amount
	return nil
}

func (c *Cash) Pay(amount int) error {
	if c.Amount < amount {
		return fmt.Errorf("insufficient amount in wallet: left: %d, required: %d", c.Amount, amount)
	}
	c.Amount -= amount
	return nil
}

func (c *UnknownPayment) Pay(amount int) error {
	if c.Amount < amount {
		return fmt.Errorf("insufficient amount in wallet: left: %d, required: %d", c.Amount, amount)
	}
	c.Amount -= amount
	return nil
}

func (c *Card) Pay(amount int) error {
	if !c.IsValid {
		return fmt.Errorf("invalid card")
	}
	if c.Cash < amount {
		return fmt.Errorf("insufficient amount in wallet: left: %d, required: %d", c.Cash, amount)
	}
	c.Cash -= amount
	return nil

}

func Buy(p Payer, amount int) bool {
	switch p.(type) {
	case *Wallet:
		fmt.Println("Wallet payment")
	case *Card:
		fmt.Println("Card payment")
	case *Cash: fmt.Println("Cash payment")
	default:
		fmt.Println("Unavailable payment method, exit...")
		return false
	}
	err := p.Pay(amount)
	if err != nil {
		fmt.Println(err)
		return false
	}
	fmt.Printf("Payment with %s \n", reflect.TypeOf(p))
	return true
}

// Empty interfaces - Can be used when there is NO requirement for realisation
// for example fmt.Printf takes empty interface


// Composition of interfaces
type Phone interface {
	Ring(string) error
}

type PhoneWithPay interface {
	Phone
	Payer
}


func interfaces() {
	w := Wallet{1, 25}
	firstPayment := Buy(&w, 20)
	secondPayment := Buy(&w, 6)
	fmt.Println(firstPayment, secondPayment)
	cardOne := Card{20, "Cardholder Name", true}
	cardTwo := Card{Cash: 10000, Cardholder: "Second Cardholder Name"}
	Buy(&cardOne, 5)
	Buy(&cardTwo, 50)
	Buy(&Cash{55}, 40)
	Buy(&UnknownPayment{55}, 40)

}