package BLC

import "fmt"

func (cli *CLI) getAddresses(nodeID string) {
	fmt.Println("打印所有钱包地址...")

	wallets := NewWallets(nodeID)
	for address, _ := range wallets.Wallets {
		fmt.Printf("address [%s]\n", address)
	}
}
