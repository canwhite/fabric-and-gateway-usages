package network

import (
	"crypto/x509"
	"fmt"
	"os"

	"github.com/hyperledger/fabric-gateway/pkg/identity"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// NewGrpcConnection creates a new gRPC client connection to the gateway
func NewGrpcConnection() (*grpc.ClientConn, error) {
	// Load TLS certificate for peer0.org1.example.com
	tlsCertificatePEM, err := os.ReadFile("../test-network/organizations/peerOrganizations/org1.example.com/tlsca/tlsca.org1.example.com-cert.pem")
	if err != nil {
		return nil, fmt.Errorf("failed to read TLS certificate: %w", err)
	}

	tlsCertificate, err := identity.CertificateFromPEM(tlsCertificatePEM)
	if err != nil {
		return nil, fmt.Errorf("failed to parse TLS certificate: %w", err)
	}

	certPool := x509.NewCertPool()
	certPool.AddCert(tlsCertificate)
	transportCredentials := credentials.NewClientTLSFromCert(certPool, "peer0.org1.example.com")

	// Connect to the gateway
	return grpc.NewClient("dns:///localhost:7051", grpc.WithTransportCredentials(transportCredentials))
}

// NewIdentity creates a client identity for this Gateway connection using an X.509 certificate
func NewIdentity() *identity.X509Identity {
	// Load client certificate for User1@org1.example.com
	certificatePEM, err := os.ReadFile("../test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/signcerts/User1@org1.example.com-cert.pem")
	if err != nil {
		panic(fmt.Errorf("failed to read certificate: %w", err))
	}

	certificate, err := identity.CertificateFromPEM(certificatePEM)
	if err != nil {
		panic(fmt.Errorf("failed to parse certificate: %w", err))
	}

	id, err := identity.NewX509Identity("Org1MSP", certificate)
	if err != nil {
		panic(fmt.Errorf("failed to create identity: %w", err))
	}

	return id
}

// NewSign creates a function that generates a digital signature from a message digest using a private key
func NewSign() identity.Sign {
	// Load private key for User1@org1.example.com
	privateKeyPEM, err := os.ReadFile("../test-network/organizations/peerOrganizations/org1.example.com/users/User1@org1.example.com/msp/keystore/priv_sk")
	if err != nil {
		panic(fmt.Errorf("failed to read private key: %w", err))
	}

	privateKey, err := identity.PrivateKeyFromPEM(privateKeyPEM)
	if err != nil {
		panic(fmt.Errorf("failed to parse private key: %w", err))
	}

	sign, err := identity.NewPrivateKeySign(privateKey)
	if err != nil {
		panic(fmt.Errorf("failed to create sign function: %w", err))
	}

	return sign
}