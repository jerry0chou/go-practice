package security

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/scrypt"
)

// PasswordHasher interface for different hashing algorithms
type PasswordHasher interface {
	Hash(password string) (string, error)
	Verify(password, hash string) bool
}

// BcryptHasher implements bcrypt password hashing
type BcryptHasher struct {
	Cost int
}

// NewBcryptHasher creates a new bcrypt hasher with specified cost
func NewBcryptHasher(cost int) *BcryptHasher {
	if cost < bcrypt.MinCost {
		cost = bcrypt.DefaultCost
	}
	if cost > bcrypt.MaxCost {
		cost = bcrypt.MaxCost
	}
	return &BcryptHasher{Cost: cost}
}

// Hash hashes a password using bcrypt
func (b *BcryptHasher) Hash(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), b.Cost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// Verify verifies a password against a bcrypt hash
func (b *BcryptHasher) Verify(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// ScryptHasher implements scrypt password hashing
type ScryptHasher struct {
	N       int // CPU/memory cost parameter
	R       int // Block size parameter
	P       int // Parallelization parameter
	KeyLen  int // Key length
	SaltLen int // Salt length
}

// NewScryptHasher creates a new scrypt hasher with specified parameters
func NewScryptHasher() *ScryptHasher {
	return &ScryptHasher{
		N:       32768, // 2^15
		R:       8,     // 8 bytes
		P:       1,     // 1 parallelization
		KeyLen:  32,    // 32 bytes
		SaltLen: 16,    // 16 bytes
	}
}

// Hash hashes a password using scrypt
func (s *ScryptHasher) Hash(password string) (string, error) {
	// Generate random salt
	salt := make([]byte, s.SaltLen)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	// Hash password with scrypt
	hash, err := scrypt.Key([]byte(password), salt, s.N, s.R, s.P, s.KeyLen)
	if err != nil {
		return "", err
	}

	// Encode parameters and hash
	encoded := s.encodeHash(hash, salt)
	return encoded, nil
}

// Verify verifies a password against a scrypt hash
func (s *ScryptHasher) Verify(password, hash string) bool {
	// Decode hash to get salt and parameters
	decodedHash, salt, err := s.decodeHash(hash)
	if err != nil {
		return false
	}

	// Hash the provided password with the same salt and parameters
	computedHash, err := scrypt.Key([]byte(password), salt, s.N, s.R, s.P, s.KeyLen)
	if err != nil {
		return false
	}

	// Compare hashes using constant time comparison
	return subtle.ConstantTimeCompare(decodedHash, computedHash) == 1
}

// encodeHash encodes hash with parameters for storage
func (s *ScryptHasher) encodeHash(hash, salt []byte) string {
	// Format: $scrypt$N$r$p$salt$hash
	encodedSalt := base64.StdEncoding.EncodeToString(salt)
	encodedHash := base64.StdEncoding.EncodeToString(hash)

	return fmt.Sprintf("$scrypt$%d$%d$%d$%s$%s",
		s.N, s.R, s.P, encodedSalt, encodedHash)
}

// decodeHash decodes stored hash to extract salt and hash
func (s *ScryptHasher) decodeHash(encoded string) ([]byte, []byte, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "scrypt" {
		return nil, nil, fmt.Errorf("invalid scrypt hash format")
	}

	// Parse parameters
	n, err := strconv.Atoi(parts[2])
	if err != nil {
		return nil, nil, err
	}
	r, err := strconv.Atoi(parts[3])
	if err != nil {
		return nil, nil, err
	}
	p, err := strconv.Atoi(parts[4])
	if err != nil {
		return nil, nil, err
	}

	// Decode salt and hash
	salt, err := base64.StdEncoding.DecodeString(parts[5])
	if err != nil {
		return nil, nil, err
	}
	hash, err := base64.StdEncoding.DecodeString(parts[6])
	if err != nil {
		return nil, nil, err
	}

	// Update parameters from stored hash
	s.N = n
	s.R = r
	s.P = p

	return hash, salt, nil
}

// PasswordManager manages password operations with different hashers
type PasswordManager struct {
	hasher PasswordHasher
}

// NewPasswordManager creates a new password manager
func NewPasswordManager(hasher PasswordHasher) *PasswordManager {
	return &PasswordManager{hasher: hasher}
}

// HashPassword hashes a password using the configured hasher
func (p *PasswordManager) HashPassword(password string) (string, error) {
	return p.hasher.Hash(password)
}

// VerifyPassword verifies a password against a hash
func (p *PasswordManager) VerifyPassword(password, hash string) bool {
	return p.hasher.Verify(password, hash)
}

// ValidatePasswordStrength validates password strength
func (p *PasswordManager) ValidatePasswordStrength(password string) error {
	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	hasUpper := false
	hasLower := false
	hasDigit := false
	hasSpecial := false

	for _, char := range password {
		switch {
		case char >= 'A' && char <= 'Z':
			hasUpper = true
		case char >= 'a' && char <= 'z':
			hasLower = true
		case char >= '0' && char <= '9':
			hasDigit = true
		case char >= 33 && char <= 126: // Special characters
			hasSpecial = true
		}
	}

	if !hasUpper {
		return fmt.Errorf("password must contain at least one uppercase letter")
	}
	if !hasLower {
		return fmt.Errorf("password must contain at least one lowercase letter")
	}
	if !hasDigit {
		return fmt.Errorf("password must contain at least one digit")
	}
	if !hasSpecial {
		return fmt.Errorf("password must contain at least one special character")
	}

	return nil
}

// GenerateSecurePassword generates a cryptographically secure random password
func (p *PasswordManager) GenerateSecurePassword(length int) (string, error) {
	if length < 8 {
		length = 12
	}

	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*"
	password := make([]byte, length)

	for i := range password {
		randomBytes := make([]byte, 1)
		if _, err := rand.Read(randomBytes); err != nil {
			return "", err
		}
		password[i] = charset[randomBytes[0]%byte(len(charset))]
	}

	return string(password), nil
}
