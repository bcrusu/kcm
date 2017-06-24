package util

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"time"

	"github.com/pkg/errors"
)

func CreateCACertificate(commonName string) (cert []byte, key []byte, err error) {
	template, err := newCertificateTemplate(commonName)
	if err != nil {
		return nil, nil, err
	}

	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign | x509.KeyUsageCRLSign

	return newCertificate(template, nil, nil, 2048)
}

func CreateServerCertificate(commonName string, signer *x509.Certificate, signerKey *rsa.PrivateKey, hosts ...string) (cert []byte, key []byte, err error) {
	template, err := newCertificateTemplate(commonName, hosts...)
	if err != nil {
		return nil, nil, err
	}

	template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth}

	return newCertificate(template, signer, signerKey, 2048)
}

func CreateClientCertificate(commonName string, signer *x509.Certificate, signerKey *rsa.PrivateKey, hosts ...string) (cert []byte, key []byte, err error) {
	template, err := newCertificateTemplate(commonName, hosts...)
	if err != nil {
		return nil, nil, err
	}

	template.ExtKeyUsage = []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth}

	return newCertificate(template, signer, signerKey, 2048)
}

func ParseCertificate(data []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("certs: input does not contain a PEM block")
	}

	if block.Type != "CERTIFICATE" {
		return nil, errors.New("certs: input does not contain a x509 certificate")
	}

	return x509.ParseCertificate(block.Bytes)
}

func ParsePrivateKey(data []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("certs: input does not contain a PEM block")
	}

	if block.Type != "RSA PRIVATE KEY" {
		return nil, errors.New("certs: input does not contain a RSA private key")
	}

	return x509.ParsePKCS1PrivateKey(block.Bytes)
}

func newCertificate(template, signer *x509.Certificate, signerKey *rsa.PrivateKey, rsaBits int) (cert []byte, key []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "certs: failed to generate RSA private key")
	}

	if signer == nil {
		signer = template
		signerKey = privateKey
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, signer, &privateKey.PublicKey, signerKey)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "certs: failed to create certificate")
	}

	cert = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	key = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	return cert, key, nil
}

func newCertificateTemplate(commonName string, hosts ...string) (*x509.Certificate, error) {
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.Wrapf(err, "certs: failed to generate serial number")
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"kcm"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		BasicConstraintsValid: true,
	}

	for _, n := range hosts {
		if ip := net.ParseIP(n); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, n)
		}
	}

	return template, nil
}
