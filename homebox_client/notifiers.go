package homeboxclient

type Notifier struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	URL       string `json:"url"`
	IsActive  bool   `json:"isActive"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
	GroupID   string `json:"groupId"`
	UserID    string `json:"userId"`
}

type NotifierCreate struct {
	Name     string `json:"name"`
	URL      string `json:"url"`
	IsActive bool   `json:"isActive"`
}

type NotifierUpdate struct {
	Name     string `json:"name"`
	URL      string `json:"url,omitempty"`
	IsActive bool   `json:"isActive"`
}

type NotifiersService struct {
	client *Client
}

func (s *NotifiersService) List() ([]Notifier, error) {
	req, err := s.client.newRequest("GET", "/v1/notifiers", nil)
	if err != nil {
		return nil, err
	}

	var notifiers []Notifier
	if err := s.client.do(req, &notifiers); err != nil {
		return nil, err
	}

	return notifiers, nil
}

func (s *NotifiersService) Create(notifier *NotifierCreate) (*Notifier, error) {
	req, err := s.client.newRequest("POST", "/v1/notifiers", notifier)
	if err != nil {
		return nil, err
	}

	var created Notifier
	if err := s.client.do(req, &created); err != nil {
		return nil, err
	}

	return &created, nil
}

func (s *NotifiersService) Update(id string, notifier *NotifierUpdate) (*Notifier, error) {
	req, err := s.client.newRequest("PUT", "/v1/notifiers/"+id, notifier)
	if err != nil {
		return nil, err
	}

	var updated Notifier
	if err := s.client.do(req, &updated); err != nil {
		return nil, err
	}

	return &updated, nil
}

func (s *NotifiersService) Delete(id string) error {
	req, err := s.client.newRequest("DELETE", "/v1/notifiers/"+id, nil)
	if err != nil {
		return err
	}

	return s.client.do(req, nil)
}

func (s *NotifiersService) Test(id string, testURL string) error {
	req, err := s.client.newRequest("POST", "/v1/notifiers/test?url="+testURL, nil)
	if err != nil {
		return err
	}

	return s.client.do(req, nil)
}
