package main

import (
	"fmt"

	"github.com/Ilhom5005/wallet/v1/pkg/wallet"
)

func main() {
	svc := &wallet.Service{}
	account, err := svc.RegisterAccount("+992987360003")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(account)

	err = svc.Deposit(account.ID, -100)
	if err != nil {
		switch err{
		case wallet.ErrAccountNotFound:
			fmt.Println("not found")
		case wallet.ErrAmountMustBePossitive:
			fmt.Println("must be positive")
		case wallet.ErrPhoneregistered:
			fmt.Println("phone registered")
		}
	}
	fmt.Println(account.Balance)

	//10
}
