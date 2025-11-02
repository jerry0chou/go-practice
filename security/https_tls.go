package security

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"net/http"
	"time"
)

// TLSServerConfig holds TLS server configuration
type TLSServerConfig struct {
	CertFile string
	KeyFile  string
	MinTLS   uint16
	MaxTLS   uint16
}

// TLSClientConfig holds TLS client configuration
type TLSClientConfig struct {
	InsecureSkipVerify bool
	MinTLS             uint16
	MaxTLS             uint16
	ServerName         string
}

// TLSSecurity handles TLS/HTTPS security operations
type TLSSecurity struct {
	serverConfig *TLSServerConfig
	clientConfig *TLSClientConfig
}

// NewTLSSecurity creates a new TLS security instance
func NewTLSSecurity() *TLSSecurity {
	return &TLSSecurity{
		serverConfig: &TLSServerConfig{
			MinTLS: tls.VersionTLS12,
			MaxTLS: tls.VersionTLS13,
		},
		clientConfig: &TLSClientConfig{
			InsecureSkipVerify: false,
			MinTLS:             tls.VersionTLS12,
			MaxTLS:             tls.VersionTLS13,
		},
	}
}

// SetServerConfig sets the TLS server configuration
func (t *TLSSecurity) SetServerConfig(config *TLSServerConfig) {
	t.serverConfig = config
}

// SetClientConfig sets the TLS client configuration
func (t *TLSSecurity) SetClientConfig(config *TLSClientConfig) {
	t.clientConfig = config
}

// CreateServerTLSConfig creates a TLS configuration for servers
func (t *TLSSecurity) CreateServerTLSConfig() (*tls.Config, error) {
	config := &tls.Config{
		MinVersion:               t.serverConfig.MinTLS,
		MaxVersion:               t.serverConfig.MaxTLS,
		InsecureSkipVerify:       false,
		PreferServerCipherSuites: true,
		CipherSuites: []uint16{
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_RSA_WITH_AES_128_GCM_SHA256,
		},
		CurvePreferences: []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		},
	}

	// Load certificate and key if provided
	if t.serverConfig.CertFile != "" && t.serverConfig.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(t.serverConfig.CertFile, t.serverConfig.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load certificate: %v", err)
		}
		config.Certificates = []tls.Certificate{cert}
	}

	return config, nil
}

// CreateClientTLSConfig creates a TLS configuration for clients
func (t *TLSSecurity) CreateClientTLSConfig() *tls.Config {
	return &tls.Config{
		MinVersion:         t.clientConfig.MinTLS,
		MaxVersion:         t.clientConfig.MaxTLS,
		InsecureSkipVerify: t.clientConfig.InsecureSkipVerify,
		ServerName:         t.clientConfig.ServerName,
	}
}

// CreateHTTPSClient creates an HTTP client with TLS configuration
func (t *TLSSecurity) CreateHTTPSClient() *http.Client {
	transport := &http.Transport{
		TLSClientConfig: t.CreateClientTLSConfig(),
	}

	return &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}
}

// GenerateSelfSignedCert generates a self-signed certificate for development
func (t *TLSSecurity) GenerateSelfSignedCert(host string) ([]byte, []byte, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, nil, err
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization:  []string{"Go Practice App"},
			Country:       []string{"US"},
			Province:      []string{""},
			Locality:      []string{"San Francisco"},
			StreetAddress: []string{""},
			PostalCode:    []string{""},
		},
		NotBefore:   time.Now(),
		NotAfter:    time.Now().Add(365 * 24 * time.Hour), // 1 year
		KeyUsage:    x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage: []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		IPAddresses: []net.IP{net.IPv4(127, 0, 0, 1), net.IPv6loopback},
		DNSNames:    []string{host, "localhost"},
	}

	// Create certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, err
	}

	// Encode certificate
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	// Encode private key
	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	return certPEM, keyPEM, nil
}

// ValidateCertificate validates a certificate
func (t *TLSSecurity) ValidateCertificate(certPEM []byte) error {
	block, _ := pem.Decode(certPEM)
	if block == nil {
		return fmt.Errorf("failed to decode PEM block")
	}

	cert, err := x509.ParseCertificate(block.Bytes)
	if err != nil {
		return fmt.Errorf("failed to parse certificate: %v", err)
	}

	// Check if certificate is expired
	if time.Now().After(cert.NotAfter) {
		return fmt.Errorf("certificate has expired")
	}

	// Check if certificate is not yet valid
	if time.Now().Before(cert.NotBefore) {
		return fmt.Errorf("certificate is not yet valid")
	}

	return nil
}

// GetTLSVersionString returns a human-readable TLS version
func (t *TLSSecurity) GetTLSVersionString(version uint16) string {
	switch version {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return "Unknown"
	}
}

// IsSecureTLSVersion checks if the TLS version is secure
func (t *TLSSecurity) IsSecureTLSVersion(version uint16) bool {
	// Only TLS 1.2 and 1.3 are considered secure
	return version >= tls.VersionTLS12
}

// CreateSecureServer creates a secure HTTP server with TLS
func (t *TLSSecurity) CreateSecureServer(addr string, handler http.Handler) (*http.Server, error) {
	tlsConfig, err := t.CreateServerTLSConfig()
	if err != nil {
		return nil, err
	}

	server := &http.Server{
		Addr:         addr,
		Handler:      handler,
		TLSConfig:    tlsConfig,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return server, nil
}

// AddSecurityHeaders adds security headers to HTTP responses
func (t *TLSSecurity) AddSecurityHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add security headers
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "DENY")
		w.Header().Set("X-XSS-Protection", "1; mode=block")
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
		w.Header().Set("Content-Security-Policy", "default-src 'self'")

		next.ServeHTTP(w, r)
	})
}
