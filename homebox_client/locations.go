package homeboxclient

import (
	"fmt"
	"net/url"
)

type LocationsService struct {
	client *Client
}

type Location struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	CreatedAt   string     `json:"createdAt"`
	UpdatedAt   string     `json:"updatedAt"`
	Parent      *Location  `json:"parent,omitempty"`
	Children    []Location `json:"children,omitempty"`
	TotalPrice  float64    `json:"totalPrice,omitempty"`
}

type LocationCreate struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ParentID    string `json:"parentId,omitempty"`
}

type LocationUpdate struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	ParentID    string `json:"parentId,omitempty"`
}

func (s *LocationsService) List(filterChildren bool) ([]Location, error) {
	u := url.Values{}
	if filterChildren {
		u.Set("filterChildren", "true")
	}

	req, err := s.client.newRequest("GET", "/v1/locations?"+u.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var locations []Location
	if err := s.client.do(req, &locations); err != nil {
		return nil, err
	}

	return locations, nil
}

func (s *LocationsService) Get(id string) (*Location, error) {
	req, err := s.client.newRequest("GET", fmt.Sprintf("/v1/locations/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var location Location
	if err := s.client.do(req, &location); err != nil {
		return nil, err
	}

	return &location, nil
}

func (s *LocationsService) Create(location *LocationCreate) (*Location, error) {
	req, err := s.client.newRequest("POST", "/v1/locations", location)
	if err != nil {
		return nil, err
	}

	var created Location
	if err := s.client.do(req, &created); err != nil {
		return nil, err
	}

	return &created, nil
}

func (s *LocationsService) Update(id string, location *LocationUpdate) (*Location, error) {
	req, err := s.client.newRequest("PUT", fmt.Sprintf("/v1/locations/%s", id), location)
	if err != nil {
		return nil, err
	}

	var updated Location
	if err := s.client.do(req, &updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

func (s *LocationsService) Delete(id string) error {
	req, err := s.client.newRequest("DELETE", fmt.Sprintf("/v1/locations/%s", id), nil)
	if err != nil {
		return err
	}

	return s.client.do(req, nil)
}

func (s *LocationsService) GetTree(withItems bool) ([]Location, error) {
	u := url.Values{}
	if withItems {
		u.Set("withItems", "true")
	}

	req, err := s.client.newRequest("GET", "/v1/locations/tree?"+u.Encode(), nil)
	if err != nil {
		return nil, err
	}

	var tree []Location
	if err := s.client.do(req, &tree); err != nil {
		return nil, err
	}

	return tree, nil
}
