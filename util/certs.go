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
	template, err := newCertificateTemplate()
	if err != nil {
		return nil, nil, err
	}

	template.Subject.CommonName = commonName
	template.IsCA = true
	template.KeyUsage |= x509.KeyUsageCertSign

	return newCertificate(template, nil, 2048)
}

func CreateCertificate(commonName string, parent *x509.Certificate, hosts ...string) (cert []byte, key []byte, err error) {
	template, err := newCertificateTemplate(hosts...)
	if err != nil {
		return nil, nil, err
	}

	template.Subject.CommonName = commonName

	return newCertificate(template, parent, 2048)
}

func ParseCertificate(data []byte) (*x509.Certificate, error) {
	block, _ := pem.Decode(data)
	if block == nil {
		return nil, errors.New("input does not contain a PEM block")
	}

	if block.Type != "CERTIFICATE" {
		return nil, errors.New("input does not contain a x509 certificate")
	}

	return x509.ParseCertificate(block.Bytes)
}

func newCertificate(template, parent *x509.Certificate, rsaBits int) (cert []byte, key []byte, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, rsaBits)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to generate RSA private key")
	}

	if parent == nil {
		parent = template
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, parent, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, errors.Wrapf(err, "failed to create certificate")
	}

	cert = pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	key = pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(privateKey)})

	return cert, key, nil
}

func newCertificateTemplate(hosts ...string) (*x509.Certificate, error) {
	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to generate serial number")
	}

	template := &x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"kcm"},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
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
