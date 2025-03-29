package hlaservice

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

const (
	loginPath = "/auth/login"
)

type loginPayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type user struct {
	Token      string `json:"token"`
	CustomerID string `json:"customer_id"`
}

// Login logs in to the HLA service and returns a token and customer ID
func (h *HLA) loginRequest(username, password string) (*user, error) {
	loginURL := h.baseURL + loginPath

	loginCredentials := loginPayload{
		Username: username,
		Password: password,
	}

	payload, err := json.Marshal(loginCredentials)
	if err != nil {
		return nil, fmt.Errorf("Error marshalling credentials: %v", err)
	}

	req, err := http.NewRequest("POST", loginURL, io.NopCloser(strings.NewReader(string(payload))))
	if err != nil {
		return nil, fmt.Errorf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Error sending request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Login failed: %v: %s", resp.Status, string(body))
	}

	var user *user
	if err = json.Unmarshal(body, &user); err != nil {
		return nil, fmt.Errorf("Error unmarshalling login response: %v", err)
	}

	return user, nil
}
