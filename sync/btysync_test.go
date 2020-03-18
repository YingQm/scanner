package sync

import "testing"

func TestGetBtyBlockHeight(t *testing.T) {
	url := "http://39.107.60.254:8801"
	height, err := GetBtyBlockHeight(url)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log("GetBtyBlockHeight return height:", height)
}

func TestGetBtyParallelHeight(t *testing.T) {
	height, err := GetBtyParallelHeight()
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log("insight-api return height:", height)
}
