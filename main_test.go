package main

import (
	"testing"
	"time"
)

func TestNormalizeTicketType(t *testing.T) {
	tests := map[string]string{
		"Étudiant":           "etudiant",
		"Admission Générale": "admission generale",
		"PRÉVENTE":          "prevente",
	}

	for input, want := range tests {
		if got := normalizeTicketType(input); got != want {
			t.Errorf("normalizeTicketType(%q) = %q, want %q", input, got, want)
		}
	}
}

func TestGetMostRecentlyEndedCampaign(t *testing.T) {
	now := time.Now().Unix()
	campaigns := []Campaign{
		{ID: "older", EndDate: now - 200},
		{ID: "latest", EndDate: now - 100},
		{ID: "upcoming", EndDate: now + 100},
	}

	got, err := GetMostRecentlyEndedCampaign(campaigns)
	if err != nil {
		t.Fatalf("GetMostRecentlyEndedCampaign() error = %v", err)
	}
	if got.ID != "latest" {
		t.Errorf("GetMostRecentlyEndedCampaign() = %q, want %q", got.ID, "latest")
	}
}

func TestGetMostRecentlyEndedCampaignNoneEnded(t *testing.T) {
	_, err := GetMostRecentlyEndedCampaign([]Campaign{{ID: "upcoming", EndDate: time.Now().Unix() + 100}})
	if err == nil {
		t.Fatal("GetMostRecentlyEndedCampaign() expected an error when no campaign has ended")
	}
}
