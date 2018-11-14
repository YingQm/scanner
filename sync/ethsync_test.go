package sync

import (
	"fmt"
	"testing"
)

func TestEthSync(t *testing.T) {
	ServiceAddr := []string{"http://39.107.57.133:8545", "http://39.107.60.254:8545", "http://47.106.117.142:8801"}
	fmt.Println(EthSync(ServiceAddr))
}

func TestDcrSync(t *testing.T) {
	ServiceAddr := []string{"https://47.106.117.142:8804"}
	fmt.Println(DcrSync(ServiceAddr))
}
