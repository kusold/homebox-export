package homeboxclient

import (
	"fmt"
	"time"
)

type MaintenanceEntry struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Cost          string    `json:"cost"`
	ScheduledDate time.Time `json:"scheduledDate"`
	CompletedDate time.Time `json:"completedDate"`
}

type MaintenanceEntryWithDetails struct {
	MaintenanceEntry
	ItemID   string `json:"itemID"`
	ItemName string `json:"itemName"`
}

type MaintenanceFilterStatus string

const (
	MaintenanceFilterStatusScheduled MaintenanceFilterStatus = "scheduled"
	MaintenanceFilterStatusCompleted MaintenanceFilterStatus = "completed"
	MaintenanceFilterStatusBoth      MaintenanceFilterStatus = "both"
)

type MaintenanceService struct {
	client *Client
}

func (s *MaintenanceService) List(status MaintenanceFilterStatus) ([]MaintenanceEntryWithDetails, error) {
	req, err := s.client.newRequest("GET", fmt.Sprintf("/v1/maintenance?status=%s", status), nil)
	if err != nil {
		return nil, err
	}

	var entries []MaintenanceEntryWithDetails
	if err := s.client.do(req, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func (s *MaintenanceService) GetItemMaintenance(itemID string, status MaintenanceFilterStatus) ([]MaintenanceEntryWithDetails, error) {
	req, err := s.client.newRequest("GET", fmt.Sprintf("/v1/items/%s/maintenance?status=%s", itemID, status), nil)
	if err != nil {
		return nil, err
	}

	var entries []MaintenanceEntryWithDetails
	if err := s.client.do(req, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}

func (s *MaintenanceService) Create(itemID string, entry *MaintenanceEntry) (*MaintenanceEntry, error) {
	req, err := s.client.newRequest("POST", fmt.Sprintf("/v1/items/%s/maintenance", itemID), entry)
	if err != nil {
		return nil, err
	}

	var created MaintenanceEntry
	if err := s.client.do(req, &created); err != nil {
		return nil, err
	}

	return &created, nil
}

func (s *MaintenanceService) Update(id string, entry *MaintenanceEntry) (*MaintenanceEntry, error) {
	req, err := s.client.newRequest("PUT", fmt.Sprintf("/v1/maintenance/%s", id), entry)
	if err != nil {
		return nil, err
	}

	var updated MaintenanceEntry
	if err := s.client.do(req, &updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

func (s *MaintenanceService) Delete(id string) error {
	req, err := s.client.newRequest("DELETE", fmt.Sprintf("/v1/maintenance/%s", id), nil)
	if err != nil {
		return err
	}

	return s.client.do(req, nil)
}
