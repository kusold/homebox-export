package homeboxclient

import "fmt"

type Label struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CreatedAt   string `json:"createdAt"`
	UpdatedAt   string `json:"updatedAt"`
}

type LabelCreate struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Color       string `json:"color"`
}

type LabelsService struct {
	client *Client
}

func (s *LabelsService) List() ([]Label, error) {
	req, err := s.client.newRequest("GET", "/v1/labels", nil)
	if err != nil {
		return nil, err
	}

	var labels []Label
	if err := s.client.do(req, &labels); err != nil {
		return nil, err
	}

	return labels, nil
}

func (s *LabelsService) Get(id string) (*Label, error) {
	req, err := s.client.newRequest("GET", fmt.Sprintf("/v1/labels/%s", id), nil)
	if err != nil {
		return nil, err
	}

	var label Label
	if err := s.client.do(req, &label); err != nil {
		return nil, err
	}

	return &label, nil
}

func (s *LabelsService) Create(label *LabelCreate) (*Label, error) {
	req, err := s.client.newRequest("POST", "/v1/labels", label)
	if err != nil {
		return nil, err
	}

	var created Label
	if err := s.client.do(req, &created); err != nil {
		return nil, err
	}

	return &created, nil
}

func (s *LabelsService) Update(id string, label *Label) (*Label, error) {
	req, err := s.client.newRequest("PUT", fmt.Sprintf("/v1/labels/%s", id), label)
	if err != nil {
		return nil, err
	}

	var updated Label
	if err := s.client.do(req, &updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

func (s *LabelsService) Delete(id string) error {
	req, err := s.client.newRequest("DELETE", fmt.Sprintf("/v1/labels/%s", id), nil)
	if err != nil {
		return err
	}

	return s.client.do(req, nil)
}
