package main

import (
    "fmt"
    "encoding/json"
    "log"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
    "github.com/przemyslawS99/ext4-blockchain-integration/ext4-blockchain-daemon/internal/common"
)

type SmartContract struct {
    contractapi.Contract
}

type Time struct {
    Sec  uint64 `json:"sec"`
    Nsec uint32 `json:"nsec"`
}

type Asset struct {
    Uid   uint32 `json:"uid"`
    Gid   uint32 `json:"gid"`
    Atime Time   `json:"atime"`
    Mtime Time   `json:"mtime"`
    Ctime Time   `json:"ctime"`
    Mode  uint32 `json:"mode"`
    Ino   uint64 `json:"ino"`
}

func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, attrs common.Attrs) error {
    exists, err := s.AssetExists(ctx, attrs.Ino)
  if err != nil {
    return err
  }
  if exists {
    return fmt.Errorf("the asset %d already exists", attrs.Ino)
  }
    asset := Asset{
        Uid:   attrs.Uid,
        Gid:   attrs.Gid,
        Atime: Time{Sec: attrs.Atime.Sec, Nsec: attrs.Atime.Nsec},
        Mtime: Time{Sec: attrs.Mtime.Sec, Nsec: attrs.Mtime.Nsec},
        Ctime: Time{Sec: attrs.Ctime.Sec, Nsec: attrs.Ctime.Nsec},
        Mode:  attrs.Mode,
        Ino:   attrs.Ino,
    } 
  assetJSON, err := json.Marshal(asset)
  if err != nil {
    return err
  }
     assetKey := fmt.Sprintf("asset_%d", asset.Ino)
  return ctx.GetStub().PutState(assetKey, assetJSON)
}

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, ino uint64) (*Asset, error) {
    assetKey := fmt.Sprintf("asset_%d", ino)
    assetJSON, err := ctx.GetStub().GetState(assetKey)
    if err != nil {
        return nil, fmt.Errorf("failed to read asset: %v", err)
    }
    if assetJSON == nil {
        return nil, fmt.Errorf("asset %d does not exist", ino)
    }

    var asset Asset
    err = json.Unmarshal(assetJSON, &asset)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal asset: %v", err)
    }

    return &asset, nil
}

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, ino uint64) (*Asset, error) {
    assetKey := fmt.Sprintf("asset_%d", ino)
    assetJSON, err := ctx.GetStub().GetState(assetKey)
    if err != nil {
        return nil, fmt.Errorf("failed to read asset: %v", err)
    }
    if assetJSON == nil {
        return nil, fmt.Errorf("asset %d does not exist", ino)
    }

    var asset Asset
    err = json.Unmarshal(assetJSON, &asset)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal asset: %v", err)
    }

    return &asset, nil
}

func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, ino uint64) (bool, error) {
    assetKey := fmt.Sprintf("asset_%d", ino)
    assetJSON, err := ctx.GetStub().GetState(assetKey)
    if err != nil {
        return false, fmt.Errorf("failed to read asset: %v", err)
    }
    return assetJSON != nil, nil
}

func main() {
    assetChaincode, err := contractapi.NewChaincode(&SmartContract{})
  if err != nil {
    log.Panicf("Error creating ext4-blockchain chaincode: %v", err)
  }

  if err := assetChaincode.Start(); err != nil {
    log.Panicf("Error starting ext4-blockchain chaincode: %v", err)
  }
}
