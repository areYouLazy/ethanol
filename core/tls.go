package core

import (
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"net"
	"os"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	hosts    = []string{"ethanol", "127.0.0.1"}
	validFor = 3650 * 24 * time.Hour
)

func publicKey(priv any) any {
	switch k := priv.(type) {
	case *rsa.PrivateKey:
		return &k.PublicKey
	case *ecdsa.PrivateKey:
		return &k.PublicKey
	case ed25519.PrivateKey:
		return k.Public().(ed25519.PublicKey)
	default:
		return nil
	}
}

func checkEthanolTLSCertificateKeyPair() error {
	_, err := os.ReadFile(viper.GetString("ethanol.server.tls.certificate"))
	if err != nil {
		return err
	}

	_, err = os.ReadFile(viper.GetString("ethanol.server.tls.key"))
	if err != nil {
		return err
	}

	_, err = tls.LoadX509KeyPair(viper.GetString("ethanol.server.tls.certificate"), viper.GetString("ethanol.server.tls.key"))
	if err != nil {
		return err
	}

	return nil
}

func getEthanolTLSCertificateKeyPair() tls.Certificate {
	certKeyPair, err := tls.LoadX509KeyPair(viper.GetString("ethanol.server.tls.certificate"), viper.GetString("ethanol.server.tls.key"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("error getting ethanol tls certificate key pair")
	}

	return certKeyPair
}

func getEthanolTLSCertificateKeyPairAsSlice() []tls.Certificate {
	certKeyPair, err := tls.LoadX509KeyPair(viper.GetString("ethanol.server.tls.certificate"), viper.GetString("ethanol.server.tls.key"))
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("error getting ethanol tls certificate key pair")
	}

	return []tls.Certificate{certKeyPair}
}

func generateEthanolTLSCertificateKeyPair() {
	var priv any
	var err error

	// // try to generate ED25519 key
	// _, priv, err = ed25519.GenerateKey(rand.Reader)
	// if err != nil {
	// 	// fallback to 4096bit RSA key
	// 	priv, err = rsa.GenerateKey(rand.Reader, 4096)
	// 	if err != nil {
	// 		logrus.WithFields(logrus.Fields{
	// 			"error": err.Error(),
	// 		}).Fatal("Failed to generate private key")
	// 	}
	// }

	// generate RSA
	priv, err = rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Failed to generate private key")
	}

	// ECDSA, ED25519 and RSA subject keys should have the DigitalSignature
	// KeyUsage bits set in the x509.Certificate template
	keyUsage := x509.KeyUsageDigitalSignature
	// Only RSA subject keys should have the KeyEncipherment KeyUsage bits set. In
	// the context of TLS this KeyUsage is particular to RSA key exchange and
	// authentication.
	if _, isRSA := priv.(*rsa.PrivateKey); isRSA {
		keyUsage |= x509.KeyUsageKeyEncipherment
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(validFor)

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("failed to generate serial number")
	}

	template := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			Organization: []string{"Ethanol Inc."},
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              keyUsage,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	for _, h := range hosts {
		if ip := net.ParseIP(h); ip != nil {
			template.IPAddresses = append(template.IPAddresses, ip)
		} else {
			template.DNSNames = append(template.DNSNames, h)
		}
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, publicKey(priv), priv)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("failed to create certificate")
	}

	certOut, err := os.Create("ssl/server.pem")
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("failed to open server.pem for writing")
	}

	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: derBytes}); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Failed to write data to server.pem")
	}

	if err := certOut.Close(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("Error closing server.pem")
	}

	logrus.Debug("certificate ssl/server.pem generated")

	keyOut, err := os.OpenFile("ssl/server.key", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("failed to open server.key for writing")
	}

	privBytes, err := x509.MarshalPKCS8PrivateKey(priv)
	if err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("unable to marshal private key")
	}

	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: privBytes}); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("failed to write data to server.key")
	}

	if err := keyOut.Close(); err != nil {
		logrus.WithFields(logrus.Fields{
			"error": err.Error(),
		}).Fatal("error closing server.key")
	}

	logrus.Debug("private key server.key generated")
}
