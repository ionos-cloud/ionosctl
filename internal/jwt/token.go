package jwt

/*
 * This package is meant to find more details about JWT tokens.
 * - Claims:			Retrieve Claims payload off a JWT token
 * - Headers:			Retrieve Headers off a JWT token
 * - Username: 			Given a JWT, retrieve user email using the identity found in the JWT Claims; and using CloudAPI
 * 						User Management API to query the found UUID. Note that the UUID can only be queried if
 * 						its respective user is managed by (or is) the user with that JWT
 * - Kid:				Extract TokenId from JWT token Headers
 * - Uuid:				Extract UserId from JWT token Claims
 * - ContractNumber:	Extract ContractNumber from JWT token Claims
 * - Role:				Extract Role from JWT token Claims
 */

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"strings"

	"github.com/ionos-cloud/ionosctl/v6/internal/client"
)

func Claims(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("incorrect format of JWT token")
	}

	payload := strings.NewReader(parts[1])
	payloadDecoded := base64.NewDecoder(base64.StdEncoding, payload)
	decBytes, err := io.ReadAll(payloadDecoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	var claims map[string]interface{}
	if err = json.Unmarshal(decBytes, &claims); err != nil {
		return nil, fmt.Errorf("could not parse JWT payload: %w", err)
	}

	return claims, nil
}

func Headers(token string) (map[string]interface{}, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, fmt.Errorf("incorrect format of JWT token")
	}

	payload := strings.NewReader(parts[0])
	payloadDecoded := base64.NewDecoder(base64.StdEncoding, payload)
	decBytes, err := io.ReadAll(payloadDecoded)
	if err != nil {
		return nil, fmt.Errorf("failed to decode payload: %w", err)
	}

	var headers map[string]interface{}
	if err = json.Unmarshal(decBytes, &headers); err != nil {
		return nil, fmt.Errorf("could not parse JWT headers: %w", err)
	}

	return headers, nil
}

func Username(token string) (string, error) {
	claims, err := Claims(token)
	if err != nil {
		return "", fmt.Errorf("failed getting claims of JWT token: %w", err)
	}
	userid, err := Uuid(claims)
	if err != nil {
		return "", fmt.Errorf("failed getting UUID via JWT Claims: %w", err)
	}
	ls, _, err := client.Must().CloudClient.UserManagementApi.UmUsersFindById(context.Background(), userid).Depth(1).Execute()
	if err != nil {
		return "", err
	}
	return *ls.Properties.Email, nil
}

func Uuid(claims map[string]interface{}) (string, error) {
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

func Kid(headers map[string]interface{}) (string, error) {
	kidInterface, ok := headers["kid"]
	if !ok {
		return "", fmt.Errorf("could not find TokenId")
	}

	kid, ok := kidInterface.(string)
	if !ok {
		return "", fmt.Errorf("tokenId is not a string")
	}

	return kid, nil
}

func ContractNumber(claims map[string]interface{}) (int64, error) {
	identityInterface, ok := claims["identity"]
	if !ok {
		return -1, fmt.Errorf("could not find identity in JWT payload")
	}

	identity, ok := identityInterface.(map[string]interface{})
	if !ok {
		return -1, fmt.Errorf("could not parse identity in JWT payload")
	}

	contractNumberInterface, ok := identity["contractNumber"]
	if !ok {
		return -1, fmt.Errorf("could not find ContractNumber in JWT payload identity")
	}

	contractNumberFloat, ok := contractNumberInterface.(float64)
	if !ok {
		return -1, fmt.Errorf("contractNumber in JWT payload identity is not a float64")
	}

	return int64(contractNumberFloat), nil
}

func Role(claims map[string]interface{}) (string, error) {
	identityInterface, ok := claims["identity"]
	if !ok {
		return "", fmt.Errorf("could not find identity in JWT payload")
	}

	identity, ok := identityInterface.(map[string]interface{})
	if !ok {
		return "", fmt.Errorf("could not parse identity in JWT payload")
	}

	roleInterface, ok := identity["role"]
	if !ok {
		return "", fmt.Errorf("could not find uuid in JWT payload identity")
	}

	role, ok := roleInterface.(string)
	if !ok {
		return "", fmt.Errorf("uuid in JWT payload identity is not a string")
	}

	return role, nil
}

func Privileges(claims map[string]interface{}) ([]string, error) {
	identityInterface, ok := claims["identity"]
	if !ok {
		return nil, fmt.Errorf("could not find identity in JWT payload")
	}

	identity, ok := identityInterface.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("could not parse identity in JWT payload")
	}

	privilegesInterface, ok := identity["privileges"]
	if !ok {
		return nil, fmt.Errorf("could not find privileges in JWT payload identity")
	}

	privilegesInterfaceArray, ok := privilegesInterface.([]interface{})
	if !ok {
		return nil, fmt.Errorf("could not parse privileges in JWT payoad identity")
	}

	var privileges = make([]string, 0)

	for _, privInterface := range privilegesInterfaceArray {
		priv, ok := privInterface.(string)
		if !ok {
			return nil, fmt.Errorf("could not parse individual privileges")
		}

		privileges = append(privileges, priv)
	}

	return privileges, nil
}
