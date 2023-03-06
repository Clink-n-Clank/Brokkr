package fsm

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const (
	stateInitial byte = iota
	stateSkip
	stateCheckout
	stateArchive
)

type Order struct {
	m Machine[byte, *Order]

	ID         string
	StockPrice uint16
	SalePrice  uint16

	IsReceiptPrinted  bool
	IsPendingToCancel bool
	IsCanceled        bool
}

func SetSalePrice[S State, O *Order](entity O) (S, O) {
	const discount = float32(0.15) // %
	var e *Order
	e = entity

	if e.IsPendingToCancel {
		return S(stateSkip), e
	}

	effect := uint16(float32(e.StockPrice) * discount)
	e.SalePrice = e.StockPrice - effect

	return S(stateCheckout), e
}

func Checkout[S State, O *Order](entity O) (S, O) {
	var e *Order
	e = entity

	if e.IsPendingToCancel {
		return S(stateSkip), e
	}

	fmt.Printf("-- Checkout Receipt [ORDER: %s Amount: %d] ---\n", e.ID, e.SalePrice)

	e.IsReceiptPrinted = true

	return S(stateArchive), e
}

func Cancel[S State, O *Order](entity O) (S, O) {
	var e *Order
	e = entity

	if e.IsPendingToCancel {
		e.IsCanceled = true
	}

	return S(stateArchive), e
}

func TestMachineTransition(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	order := &Order{
		m:                 Machine[byte, *Order]{},
		ID:                "unit-test",
		StockPrice:        42,
		IsPendingToCancel: rand.Intn(2) == 0,
	}

	order.m.AddState(stateInitial, SetSalePrice[byte, *Order])
	order.m.AddState(stateCheckout, Checkout[byte, *Order])
	order.m.AddState(stateSkip, Cancel[byte, *Order])
	order.m.AddEndState(stateArchive)

	// assert
	order.m.MakeTransition(order)

	if !order.IsPendingToCancel {
		assert.True(t, uint16(36) == order.SalePrice, "discount must be applied")
		assert.True(t, order.IsReceiptPrinted, "expected receipt to be printed")
	} else {
		assert.True(t, order.IsCanceled, "expected to be canceled")
	}
}
