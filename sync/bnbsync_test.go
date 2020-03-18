package sync

import "testing"

func TestGetBnbBlockHeight(t *testing.T) {
	node := "http://119.3.17.220:27147"
	height, err := GetBnbBlockHeight(node)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log("GetBnBBlockHeight return height:", height)
}
