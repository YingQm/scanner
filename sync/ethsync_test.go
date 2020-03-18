package sync

import (
	"fmt"
	"testing"
)

func TestEthSync(t *testing.T) {
	url := "http://183.129.226.77:8545"

	height, err := getEthBlockHeight(url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	fmt.Println("resp height:", height)
	t.Log("getEthBlockHeight return height:", height)
}
