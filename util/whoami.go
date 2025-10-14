package util

import (
	"context"
	"fmt"

	"github.com/gotd/td/telegram"
)

func Whoami(ctx context.Context, client *telegram.Client) {
	if client == nil {
		fmt.Println("func Whoami: client is nil")
		return
	}

	me1, _ := client.Self(ctx)
	fmt.Println("-------------------------------------------------")
	fmt.Printf("Login as: %s(@%s)", me1.FirstName, me1.Username)
	fmt.Println()
	fmt.Println("-------------------------------------------------")
}
