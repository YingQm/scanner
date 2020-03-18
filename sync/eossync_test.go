package sync

import "testing"

func TestGetEosBlockHeight(t *testing.T) {
	url := "http://183.129.226.77:8888"
	height, err := GetEosBlockHeight(url)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log("GetEosBlockHeight return height:", height)

}
