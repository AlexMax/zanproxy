package main

import "testing"

func TestGetScore(t *testing.T) {
	// New IpIntel
	intel := NewIPIntel()

	// Test request
	score, cached, err := intel.GetScore("82.183.48.49", "alexmax@example.com")
	if err != nil {
		t.Errorf("Unexpected error %#v.", err)
	}
	if cached != false {
		t.Errorf("Unexpected cached %f.", score)
	}

	// Test cached request
	cscore, cached, err := intel.GetScore("82.183.48.49", "alexmax2742@gmail.com")
	if err != nil {
		t.Errorf("Unexpected error %#v.", err)
	}
	if cached != true {
		t.Errorf("Unexpected cached %f.", score)
	}
	if cscore != score {
		t.Errorf("Cached score is incorrect (%f vs %f).", cscore, score)
	}
}
