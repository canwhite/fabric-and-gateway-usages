package service

import (
	"fmt"

	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// AssetService handles interactions with the asset-transfer-basic chaincode
type AssetService struct {
	contract *client.Contract
}

// NewAssetService creates a new asset service instance
func NewAssetService(gateway *client.Gateway) *AssetService {
	network := gateway.GetNetwork("mychannel")
	contract := network.GetContract("basic")

	return &AssetService{
		contract: contract,
	}
}

// InitLedger initializes the ledger with sample assets
func (s *AssetService) InitLedger() error {
	fmt.Println("Submitting InitLedger transaction...")
	_, err := s.contract.SubmitTransaction("InitLedger")
	if err != nil {
		return fmt.Errorf("failed to init ledger: %w", err)
	}
	fmt.Println("✓ Ledger initialized successfully")
	return nil
}

// CreateAsset creates a new asset on the ledger
func (s *AssetService) CreateAsset(id, color, size, owner, value string) error {
	fmt.Printf("Creating asset %s...\n", id)
	_, err := s.contract.SubmitTransaction("CreateAsset", id, color, size, owner, value)
	if err != nil {
		return fmt.Errorf("failed to create asset %s: %w", id, err)
	}
	fmt.Printf("✓ Asset %s created successfully\n", id)
	return nil
}

// GetAllAssets returns all assets from the ledger
func (s *AssetService) GetAllAssets() (string, error) {
	fmt.Println("Querying all assets...")
	result, err := s.contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		return "", fmt.Errorf("failed to get all assets: %w", err)
	}
	return string(result), nil
}

// ReadAsset returns a specific asset by ID
func (s *AssetService) ReadAsset(id string) (string, error) {
	fmt.Printf("Reading asset %s...\n", id)
	result, err := s.contract.EvaluateTransaction("ReadAsset", id)
	if err != nil {
		return "", fmt.Errorf("failed to read asset %s: %w", id, err)
	}
	return string(result), nil
}

// UpdateAsset updates an existing asset
func (s *AssetService) UpdateAsset(id, color, size, owner, value string) error {
	fmt.Printf("Updating asset %s...\n", id)
	_, err := s.contract.SubmitTransaction("UpdateAsset", id, color, size, owner, value)
	if err != nil {
		return fmt.Errorf("failed to update asset %s: %w", id, err)
	}
	fmt.Printf("✓ Asset %s updated successfully\n", id)
	return nil
}

// DeleteAsset deletes an asset from the ledger
func (s *AssetService) DeleteAsset(id string) error {
	fmt.Printf("Deleting asset %s...\n", id)
	_, err := s.contract.SubmitTransaction("DeleteAsset", id)
	if err != nil {
		return fmt.Errorf("failed to delete asset %s: %w", id, err)
	}
	fmt.Printf("✓ Asset %s deleted successfully\n", id)
	return nil
}