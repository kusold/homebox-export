package homeboxclient

type UsersService struct {
	client *Client
}

type UserOut struct {
	ID          string `json:"id"`
	Email       string `json:"email"`
	Name        string `json:"name"`
	GroupID     string `json:"groupId"`
	GroupName   string `json:"groupName"`
	IsOwner     bool   `json:"isOwner"`
	IsSuperuser bool   `json:"isSuperuser"`
}

type UserUpdate struct {
	Email string `json:"email"`
	Name  string `json:"name"`
}

type ChangePassword struct {
	Current string `json:"current"`
	New     string `json:"new"`
}

func (s *UsersService) GetSelf() (*UserOut, error) {
	req, err := s.client.newRequest("GET", "/v1/users/self", nil)
	if err != nil {
		return nil, err
	}

	var wrapped struct {
		Item UserOut `json:"item"`
	}
	if err := s.client.do(req, &wrapped); err != nil {
		return nil, err
	}

	return &wrapped.Item, nil
}

func (s *UsersService) UpdateSelf(update UserUpdate) (*UserUpdate, error) {
	req, err := s.client.newRequest("PUT", "/v1/users/self", update)
	if err != nil {
		return nil, err
	}

	var wrapped struct {
		Item UserUpdate `json:"item"`
	}
	if err := s.client.do(req, &wrapped); err != nil {
		return nil, err
	}

	return &wrapped.Item, nil
}

func (s *UsersService) DeleteSelf() error {
	req, err := s.client.newRequest("DELETE", "/v1/users/self", nil)
	if err != nil {
		return err
	}

	return s.client.do(req, nil)
}

func (s *UsersService) ChangePassword(current, new string) error {
	change := ChangePassword{
		Current: current,
		New:     new,
	}

	req, err := s.client.newRequest("PUT", "/v1/users/change-password", change)
	if err != nil {
		return err
	}

	return s.client.do(req, nil)
}
