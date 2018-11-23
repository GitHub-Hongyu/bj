package BLC

import (
	"fmt"
	"testing"
)

func TestNewWallet(t *testing.T) {
	fmt.Printf("wallet : %v\n", NewWallet())
}