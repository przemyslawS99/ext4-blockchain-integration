package main

import (
    "fmt"
    "encoding/json"
    "log"
    "github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
    contractapi.Contract
}

type Time struct {
    Sec  string `json:"sec"`
    Nsec string `json:"nsec"`
}

type Asset struct {
    Uid   string `json:"uid"`
    Gid   string `json:"gid"`
    Atime Time   `json:"atime"`
    Mtime Time   `json:"mtime"`
    Ctime Time   `json:"ctime"`
    Mode  string `json:"mode"`
    Ino   string `json:"ino"`
}

func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, uid, gid, atimeSec, atimeNsec, mtimeSec, mtimeNsec, ctimeSec, ctimeNsec, mode, ino string) error {
    exists, err := s.AssetExists(ctx, ino)
    
    if err != nil {
        return err
    }
    
    if exists {
        return fmt.Errorf("the asset %s already exists", ino)
    }
    
    asset := Asset{
        Uid: uid,
        Gid: gid,
        Atime: Time{
            Sec:  atimeSec,
            Nsec: atimeNsec,
        },
        Mtime: Time{
            Sec:  mtimeSec,
            Nsec: mtimeNsec,
        },
        Ctime: Time{
            Sec:  ctimeSec,
            Nsec: ctimeNsec,
        },
        Mode: mode,
        Ino:  ino,
    }
  
    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return err
    }

    assetKey := fmt.Sprintf("asset_%s", ino)
    return ctx.GetStub().PutState(assetKey, assetJSON)
}

func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, uid, gid, atimeSec, atimeNsec, mtimeSec, mtimeNsec, ctimeSec, ctimeNsec, mode, ino string) error {
    exists, err := s.AssetExists(ctx, ino)
    
    if err != nil {
        return err
    }
    
    if !exists {
        return fmt.Errorf("the asset %s does not exist", ino)
    }

    asset, err := s.ReadAsset(ctx, ino)
    if err != nil {
        return err
    }
   
    if uid != "" {
        asset.Uid = uid
    }
    if gid != "" {
        asset.Gid = gid
    }
    if atimeSec != "" {
        asset.Atime.Sec = atimeSec
    }
    if atimeNsec != "" {
        asset.Atime.Nsec = atimeNsec
    }
    if mtimeSec != "" {
        asset.Mtime.Sec = mtimeSec
    }
    if mtimeNsec != "" {
        asset.Mtime.Nsec = mtimeNsec
    }
    if ctimeSec != "" {
        asset.Ctime.Sec = ctimeSec
    }
    if ctimeNsec != "" {
        asset.Ctime.Nsec = ctimeNsec
    }
    if mode != "" {
        asset.Mode = mode
    }

    assetJSON, err := json.Marshal(asset)
    if err != nil {
        return err
    }

    assetKey := fmt.Sprintf("asset_%s", ino)
    return ctx.GetStub().PutState(assetKey, assetJSON)
}

func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, ino string) (*Asset, error) {
    assetKey := fmt.Sprintf("asset_%s", ino)
    assetJSON, err := ctx.GetStub().GetState(assetKey)
    
    if err != nil {
        return nil, fmt.Errorf("failed to read asset: %v", err)
    }
    
    if assetJSON == nil {
        return nil, fmt.Errorf("asset %s does not exist", ino)
    }

    var asset Asset
    err = json.Unmarshal(assetJSON, &asset)
    if err != nil {
        return nil, fmt.Errorf("failed to unmarshal asset: %v", err)
    }

    return &asset, nil
}

func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, ino string) (bool, error) {
    assetKey := fmt.Sprintf("asset_%s", ino)
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
