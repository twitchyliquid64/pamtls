package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/donpark/pam"
)

// Modes
const (
	modeVerifyFullTLS int = iota
	modeVerifySignAndExpiry
	modeInsecure
)

var tlsMode int
var caCertParsed *x509.Certificate
var tlsConfig *tls.Config

func tlsInit(args pam.Args) error {
	tlsConfig = &tls.Config{
		MinVersion:       tls.VersionTLS12,
		CurvePreferences: []tls.CurveID{tls.CurveP521, tls.CurveP384, tls.CurveP256},
	}

	certPath, certSpecified := args["cert"]
	keyPath, keySpecified := args["key"]
	if certSpecified && keySpecified {
		mainCert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			return err
		}
		tlsConfig.Certificates = []tls.Certificate{mainCert}
		tlsConfig.BuildNameToCertificate()
	} else if certSpecified || keySpecified {
		info("TLS", "Unable to init client certificates: either cert or key missing")
	}

	if v, verifySpecified := args["verify"]; verifySpecified {
		verifySpec := strings.ToLower(v)
		switch verifySpec {
		case "full", "system", "tls":
			tlsMode = modeVerifyFullTLS

		case "insecure":
			tlsMode = modeInsecure
			tlsConfig.InsecureSkipVerify = true

		case "pinned", "custom", "custom-root":
			tlsMode = modeVerifySignAndExpiry
			if args["root"] == "" {
				return errors.New("no root specified")
			}
			pemBytes, err := ioutil.ReadFile(args["root"])
			if err != nil {
				return err
			}
			certDERBlock, _ := pem.Decode(pemBytes)
			if certDERBlock == nil {
				return errors.New("no certificate data read from PEM")
			}
			caCertParsed, err = x509.ParseCertificate(certDERBlock.Bytes)
			if err != nil {
				return err
			}
			tlsConfig.VerifyPeerCertificate = verifyCertSignAndExpiry
		default:
			return errors.New("unrecognised 'verify' setting")
		}
	}

	if isDebugMode {
		info("DEBUG-TLS", fmt.Sprintf("mode=%d,client_cert=%v,insecureSkipVerify=%v", tlsMode, len(tlsConfig.Certificates) > 0, tlsConfig.InsecureSkipVerify))
	}
	transport = &http.Transport{TLSClientConfig: tlsConfig}
	return nil
}

func verifyCertSignAndExpiry(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
	for _, cert := range rawCerts {
		parsedCert, err := x509.ParseCertificate(cert)
		if err != nil {
			return err
		}
		certErr := parsedCert.CheckSignatureFrom(caCertParsed)
		if parsedCert.NotAfter.Before(time.Now()) || parsedCert.NotBefore.After(time.Now()) {
			certErr = errors.New("Certificate expired or used too soon")
		}
		return certErr
	}
	return errors.New("Expected certificate which would pass, none presented")
}
