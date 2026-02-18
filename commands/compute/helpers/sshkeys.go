package helpers

import (
	utils2 "github.com/ionos-cloud/ionosctl/v6/internal/utils"
)

func GetSshKeysFromPaths(paths []string) ([]string, error) {
	sshKeys := make([]string, 0)
	if len(paths) != 0 {
		for _, sshKeyPath := range paths {
			publicKey, err := utils2.ReadPublicKey(sshKeyPath)
			if err != nil {
				return sshKeys, err
			}
			sshKeys = append(sshKeys, publicKey)
		}
	}
	return sshKeys, nil
}
