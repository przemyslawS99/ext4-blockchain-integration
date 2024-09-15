package fabric

import (
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"

	"github.com/hyperledger/fabric-gateway/pkg/client"
	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"github.com/przemyslawS99/ext4-blockchain-integration/ext4-blockchain-daemon/internal/common"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
)

const (
	mspID        = "Org1MSP"
	cryptoPath   = "/home/przem/go/src/github.com/przemyslawS99/fabric-samples/test-network/organizations/peerOrganizations/org1.example.com"
	certPath     = cryptoPath + "/users/User1@org1.example.com/msp/signcerts"
	keyPath      = cryptoPath + "/users/User1@org1.example.com/msp/keystore"
	tlsCertPath  = cryptoPath + "/peers/peer0.org1.example.com/tls/ca.crt"
	peerEndpoint = "dns:///localhost:7051"
	gatewayPeer  = "peer0.org1.example.com"
)

func NewGrpcConnection() *grpc.ClientConn {
	certificatePEM, err := os.ReadFile(tlsCertPath)
	if err != nil {
		panic(fmt.Errorf("failed to read TLS certifcate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(certificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, gatewayPeer)

	connection, err := grpc.NewClient(peerEndpoint, grpc.WithTransportCredentials(transportCredentials))
	if err != nil {
		panic(fmt.Errorf("failed to create gRPC connection: %w", err))
	}

	return connection
}

func NewIdentity() *identity.X509Identity {
	certificatePEM, err := readFirstFile(certPath)
	if err != nil {
		panic(fmt.Errorf("failed to read certificate file: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(err)
	}

	id, err := identity.NewX509Identity(mspID, certificate)
	if err != nil {
		panic(err)
	}

	return id
}

func NewSign() identity.Sign {
	privateKeyPEM, err := readFirstFile(keyPath)
	if err != nil {
		panic(fmt.Errorf("failed to read private key file: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(err)
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(err)
	}

	return sign
}

func readFirstFile(dirPath string) ([]byte, error) {
	dir, err := os.Open(dirPath)
	if err != nil {
		return nil, err
	}

	fileNames, err := dir.Readdirnames(1)
	if err != nil {
		return nil, err
	}

	return os.ReadFile(path.Join(dirPath, fileNames[0]))
}

func NewInode(attrs *common.Attrs, contract *client.Contract) uint16 {
	log.Printf("fabric: NewInode %v", attrs.Ino)
	args := convertAttrs(attrs)
	log.Printf("args: %v", args)
	_, err := contract.SubmitTransaction("CreateAsset", args...)
	if err != nil {
		log.Printf("failed to submit transaction: %v", err)
		return common.EXT4BD_STATUS_FAIL
	}

	log.Printf("transaction committed successfully")
	return common.EXT4BD_STATUS_SUCCESS
}

func SetAttributes(attrs *common.Attrs, contract *client.Contract) uint16 {
	log.Printf("fabric: SetAttributes %v", attrs.Ino)
	args := convertAttrs(attrs)
	_, err := contract.SubmitTransaction("UpdateAsset", args...)
	if err != nil {
		return handleError(err)
	}

	log.Printf("transaction committed successfully")
	return common.EXT4BD_STATUS_SUCCESS
}

func GetAttributes(ino uint64, contract *client.Contract) (uint16, *common.Attrs) {
	log.Printf("fabric: GetAttributes %v", ino)

	var ret uint16
	var attrs *common.Attrs
	var result string

	evaluateResult, err := contract.EvaluateTransaction("ReadAsset", fmt.Sprintf("%d", ino))
	if err != nil {
		ret = handleError(err)
		goto err_out
	}

	attrs, err = parseAttrs(evaluateResult)
	if err != nil {
		ret = common.EXT4BD_STATUS_FAIL
		goto err_out
	}

	log.Printf("transaction committed successfully: %s", result)
	ret = common.EXT4BD_STATUS_SUCCESS

	return ret, attrs

err_out:
	return ret, &common.Attrs{Ino: ino}
}

func handleError(err error) uint16 {
	var endorseErr *client.EndorseError
	if errors.As(err, &endorseErr) {
		log.Printf("Endorse error for transaction %s with gRPC status %v: %s", endorseErr.TransactionID, status.Code(endorseErr), endorseErr)
		return common.EXT4BD_STATUS_INODE_NOT_FOUND
	} else {
		log.Printf("failed to submit transaction: %v", err)
		return common.EXT4BD_STATUS_FAIL
	}
}

func formatUint32(value uint32) string {
	if value == 0 {
		return ""
	}
	return strconv.FormatUint(uint64(value), 10)
}

func formatUint64(value uint64) string {
	if value == 0 {
		return ""
	}
	return strconv.FormatUint(value, 10)
}

func formatTime(t common.Time) (string, string) {
	return formatUint64(t.Sec), formatUint32(t.Nsec)
}

func convertAttrs(attrs *common.Attrs) []string {
	atimeSec, atimeNsec := formatTime(attrs.Atime)
	mtimeSec, mtimeNsec := formatTime(attrs.Mtime)
	ctimeSec, ctimeNsec := formatTime(attrs.Ctime)

	return []string{
		formatUint32(attrs.Uid),
		formatUint32(attrs.Gid),
		atimeSec,
		atimeNsec,
		mtimeSec,
		mtimeNsec,
		ctimeSec,
		ctimeNsec,
		formatUint32(attrs.Mode),
		formatUint64(attrs.Ino),
	}
}

func parseAttrs(data []byte) (*common.Attrs, error) {
	var asset struct {
		Uid   string `json:"uid"`
		Gid   string `json:"gid"`
		Atime struct {
			Sec  string `json:"sec"`
			Nsec string `json:"nsec"`
		} `json:"atime"`
		Mtime struct {
			Sec  string `json:"sec"`
			Nsec string `json:"nsec"`
		} `json:"mtime"`
		Ctime struct {
			Sec  string `json:"sec"`
			Nsec string `json:"nsec"`
		} `json:"ctime"`
		Mode string `json:"mode"`
		Ino  string `json:"ino"`
	}

	err := json.Unmarshal(data, &asset)
	if err != nil {
		log.Printf("Failed to unmarshal asset: %v", err)
		return nil, err
	}

	uid, _ := strconv.ParseUint(asset.Uid, 10, 32)
	gid, _ := strconv.ParseUint(asset.Gid, 10, 32)
	atimeSec, _ := strconv.ParseUint(asset.Atime.Sec, 10, 64)
	atimeNsec, _ := strconv.ParseUint(asset.Atime.Nsec, 10, 32)
	mtimeSec, _ := strconv.ParseUint(asset.Mtime.Sec, 10, 64)
	mtimeNsec, _ := strconv.ParseUint(asset.Mtime.Nsec, 10, 32)
	ctimeSec, _ := strconv.ParseUint(asset.Ctime.Sec, 10, 64)
	ctimeNsec, _ := strconv.ParseUint(asset.Ctime.Nsec, 10, 32)
	mode, _ := strconv.ParseUint(asset.Mode, 10, 32)
	ino, _ := strconv.ParseUint(asset.Ino, 10, 64)

	attrs := common.Attrs{
		Uid: uint32(uid),
		Gid: uint32(gid),
		Atime: common.Time{
			Sec:  atimeSec,
			Nsec: uint32(atimeNsec),
		},
		Mtime: common.Time{
			Sec:  mtimeSec,
			Nsec: uint32(mtimeNsec),
		},
		Ctime: common.Time{
			Sec:  ctimeSec,
			Nsec: uint32(ctimeNsec),
		},
		Mode: uint32(mode),
		Ino:  ino,
	}

	return &attrs, nil
}
