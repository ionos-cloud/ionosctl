package jwt

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
)

func Valid(token string) bool {
	err := client.TestCreds("", "", token)
	if err != nil {
		return false
	}

	return true
}

func Claims(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("incorrect format of JWT token")
	}

	payload, err := base64.RawURLEncoding.DecodeString(pad(parts[1]))
	if err != nil {
		return nil, fmt.Errorf("could not decode JWT payload: %v", err)
	}

	var claims map[string]interface{}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, fmt.Errorf("could not parse JWT payload: %v", err)
	}
	return claims, nil
}

func Username(token string) (string, error) {
	claims, err := Claims(token)
	if err != nil {
		return "", fmt.Errorf("failed getting claims of JWT token: %w", err)
	}
	userid, err := uuid(claims)
	if err != nil {
		return "", fmt.Errorf("failed getting UUID via JWT Claims: %w", err)
	}
	ls, _, err := client.Must().CloudClient.UserManagementApi.UmUsersFindById(context.Background(), userid).Depth(1).Execute()
	if err != nil {
		return "", err
	}
	return *ls.Properties.Email, nil
}

func uuid(claims map[string]interface{}) (string, error) {
	identityInterface, ok := claims["identity"]
	if !ok {
		return "", fmt.Errorf("could not find identity in JWT payload")
	}

	identity, ok := identityInterface.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("could not parse identity in JWT payload")
	}

	uuidInterface, ok := identity["uuid"]
	if !ok {
		return "", fmt.Errorf("could not find uuid in JWT payload identity")
	}

	id, ok := uuidInterface.(string)
	if !ok {
		return "", fmt.Errorf("uuid in JWT payload identity is not a string")
	}

	return id, nil
}

// pad adds padding to base64-encoded string, if needed
func pad(base64string string) string {
	switch len(base64string) % 4 {
	case 2:
		base64string += "=="
	case 3:
		base64string += "="
	}
	return base64string
}
