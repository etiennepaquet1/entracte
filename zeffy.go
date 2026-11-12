package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const paymentsURL = "https://api.zeffy.com/api/v1/payments"
const campaignsURL = "https://api.zeffy.com/api/v1/campaigns"
const paymentsLimit = "100"
const doorCampaignID = "c391ebd1-b8ba-484e-a60d-d046f02ebf54"

type Campaign struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	StartDate int64  `json:"start_date"`
	EndDate   int64  `json:"end_date"`
}

func GetCampaigns() ([]Campaign, error) {
	if ZeffyAPI == "" {
		return nil, fmt.Errorf("ZEFFY_API is not set; call LoadConfig first")
	}

	var campaigns []Campaign
	cursor := ""
	seenCursors := make(map[string]struct{})

	for {
		page, err := getCampaignsPage(cursor)
		if err != nil {
			return nil, err
		}
		campaigns = append(campaigns, page.Data...)
		if !page.HasMore {
			return campaigns, nil
		}
		if page.NextCursor == "" {
			return nil, fmt.Errorf("campaigns response indicates more results but has no next_cursor")
		}
		if _, exists := seenCursors[page.NextCursor]; exists || page.NextCursor == cursor {
			return nil, fmt.Errorf("campaigns pagination returned a repeated cursor")
		}
		seenCursors[page.NextCursor] = struct{}{}
		cursor = page.NextCursor
	}
}

func GetMostRecentlyEndedCampaign(campaigns []Campaign) (Campaign, error) {
	now := time.Now().Unix()
	var mostRecent Campaign
	found := false
	for _, campaign := range campaigns {
		if campaign.EndDate <= 0 || campaign.EndDate > now {
			continue
		}
		if !found || campaign.EndDate > mostRecent.EndDate {
			mostRecent = campaign
			found = true
		}
	}
	if !found {
		return Campaign{}, fmt.Errorf("no ended campaigns found")
	}
	return mostRecent, nil
}

func getCampaignsPage(cursor string) (campaignsPage, error) {
	req, err := http.NewRequest(http.MethodGet, campaignsURL, nil)
	if err != nil {
		return campaignsPage{}, fmt.Errorf("create campaigns request: %w", err)
	}
	query := url.Values{}
	query.Set("limit", paymentsLimit)
	if cursor != "" {
		query.Set("starting_after", cursor)
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+ZeffyAPI)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return campaignsPage{}, fmt.Errorf("get campaigns: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return campaignsPage{}, fmt.Errorf("read campaigns response: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return campaignsPage{}, fmt.Errorf("get campaigns: unexpected status %s", resp.Status)
	}

	var page campaignsPage
	if err := json.Unmarshal(body, &page); err != nil {
		return campaignsPage{}, fmt.Errorf("decode campaigns response: %w", err)
	}
	return page, nil
}

func GetPayments() ([]Payment, error) {
	if ZeffyAPI == "" {
		return nil, fmt.Errorf("ZEFFY_API is not set; call LoadConfig first")
	}

	var allPayments []Payment
	cursor := ""
	seenCursors := make(map[string]struct{})

	for {
		page, err := getPaymentsPage(cursor)
		if err != nil {
			return nil, err
		}
		allPayments = append(allPayments, page.Data...)
		if !page.HasMore {
			break
		}

		if page.NextCursor == "" {
			return nil, fmt.Errorf("payments response indicates more results but has no next_cursor")
		}
		if _, exists := seenCursors[page.NextCursor]; exists || page.NextCursor == cursor {
			return nil, fmt.Errorf("payments pagination returned a repeated cursor")
		}
		seenCursors[page.NextCursor] = struct{}{}
		cursor = page.NextCursor
	}

	return allPayments, nil
}

func getPaymentsPage(cursor string) (paymentsPage, error) {
	req, err := http.NewRequest(http.MethodGet, paymentsURL, nil)
	if err != nil {
		return paymentsPage{}, fmt.Errorf("create payments request: %w", err)
	}
	query := url.Values{}
	query.Set("limit", paymentsLimit)
	if cursor != "" {
		query.Set("starting_after", cursor)
	}
	req.URL.RawQuery = query.Encode()
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+ZeffyAPI)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return paymentsPage{}, fmt.Errorf("get payments: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return paymentsPage{}, fmt.Errorf("read payments response: %w", err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
		return paymentsPage{}, fmt.Errorf("get payments: unexpected status %s", resp.Status)
	}

	var payments paymentsPage
	if err := json.Unmarshal(body, &payments); err != nil {
		return paymentsPage{}, fmt.Errorf("decode payments response: %w", err)
	}
	return payments, nil
}

type paymentsPage struct {
	Data       []Payment `json:"data"`
	HasMore    bool      `json:"has_more"`
	NextCursor string    `json:"next_cursor"`
}

type campaignsPage struct {
	Data       []Campaign `json:"data"`
	HasMore    bool       `json:"has_more"`
	NextCursor string     `json:"next_cursor"`
}

type Payment struct {
	Amount      int64         `json:"amount"`
	Created     int64         `json:"created"`
	Description string        `json:"description"`
	CampaignID  string        `json:"campaign_id"`
	Type        string        `json:"type"`
	Items       []PaymentItem `json:"items"`
}

type PaymentItem struct {
	Type      string `json:"type"`
	RateTitle string `json:"rate_title"`
	Amount    int64  `json:"amount"`
}
