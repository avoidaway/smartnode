package migration

import (
	"fmt"

	"github.com/rocket-pool/smartnode/shared/types/config"
)

func upgradeFromV191(serializedConfig map[string]map[string]string) error {
	// v1.9.1 had the BN API port mode as a boolean
	if err := updateRPCPortConfig(serializedConfig, "consensusCommon", "openApiPort"); err != nil {
		return err
	}
	if err := updateRPCPortConfig(serializedConfig, "executionCommon", "openRpcPorts"); err != nil {
		return err
	}
	if err := updateRPCPortConfig(serializedConfig, "mevBoost", "openRpcPort"); err != nil {
		return err
	}
	return nil
}

func updateRPCPortConfig(serializedConfig map[string]map[string]string, configKeyString string, keyOpenPorts string) error {
	// v1.9.1 had the EC API ports mode as a boolean
	configSection, exists := serializedConfig[configKeyString]
	if !exists {
		return fmt.Errorf("expected a section called `%s` but it didn't exist", configKeyString)
	}
	openRPCPorts, exists := configSection[keyOpenPorts]
	if !exists {
		return fmt.Errorf("expected a executionCommon setting named `%s` but it didn't exist", keyOpenPorts)
	}

	// Update the config
	if openRPCPorts == "true" {
		configSection[keyOpenPorts] = string(config.RPC_OpenLocalhost)
	} else {
		configSection[keyOpenPorts] = string(config.RPC_Closed)
	}
	serializedConfig[configKeyString] = configSection
	return nil
}
