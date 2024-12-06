package homeboxclient

import "fmt"

type LoginForm struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	StayLoggedIn bool   `json:"stayLoggedIn"`
}

type TokenResponse struct {
	Token           string `json:"token"`
	AttachmentToken string `json:"attachmentToken"`
	ExpiresAt       string `json:"expiresAt"`
}

func (c *Client) Login(username, password string) (*TokenResponse, error) {
	form := LoginForm{
		Username: username,
		Password: password,
	}

	req, err := c.newRequest("POST", "/v1/users/login", form)
	if err != nil {
		return nil, err
	}

	var resp TokenResponse
	if err := c.do(req, &resp); err != nil {
		return nil, fmt.Errorf("login failed: %w", err)
	}

	c.token = resp.Token
	return &resp, nil
}

func (c *Client) Logout() error {
	req, err := c.newRequest("POST", "/v1/users/logout", nil)
	if err != nil {
		return err
	}

	return c.do(req, nil)
}
