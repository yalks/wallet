package logic

import (
	"encoding/json"
	"net"
	"regexp"
	"strings"
)

// TransactionValidator provides validation for transaction fields
type TransactionValidator struct {
	// Allowed sources
	allowedSources map[string]bool
	
	// IP validation pattern
	ipPattern *regexp.Regexp
}

// NewTransactionValidator creates a new transaction validator
func NewTransactionValidator() *TransactionValidator {
	return &TransactionValidator{
		allowedSources: map[string]bool{
			"telegram": true,
			"web":      true,
			"api":      true,
			"admin":    true,
			"unknown":  true,
		},
		ipPattern: regexp.MustCompile(`^([0-9]{1,3}\.){3}[0-9]{1,3}$`),
	}
}

// ValidateRequestSource validates the request source
func (v *TransactionValidator) ValidateRequestSource(source string) bool {
	if source == "" {
		return true // Allow empty source
	}
	
	source = strings.ToLower(strings.TrimSpace(source))
	return v.allowedSources[source]
}

// ValidateIP validates an IP address
func (v *TransactionValidator) ValidateIP(ip string) bool {
	if ip == "" {
		return true // Allow empty IP
	}
	
	// Check if it's a valid IP
	parsedIP := net.ParseIP(ip)
	return parsedIP != nil
}

// ValidateFeeType validates the fee type
func (v *TransactionValidator) ValidateFeeType(feeType string) bool {
	if feeType == "" {
		return true // Allow empty fee type
	}
	
	feeType = strings.ToLower(strings.TrimSpace(feeType))
	return feeType == "fixed" || feeType == "percentage"
}

// ValidateUserAgent validates the user agent
func (v *TransactionValidator) ValidateUserAgent(userAgent string) bool {
	if userAgent == "" {
		return true // Allow empty user agent
	}
	
	// Basic validation - ensure it's not too long
	return len(userAgent) <= 1000
}

// ValidateMetadataJSON validates JSON metadata
func (v *TransactionValidator) ValidateMetadataJSON(metadataJSON string) bool {
	if metadataJSON == "" || metadataJSON == "{}" {
		return true
	}
	
	// Try to parse as JSON
	var data map[string]interface{}
	err := json.Unmarshal([]byte(metadataJSON), &data)
	return err == nil
}

// SanitizeRequestSource sanitizes and normalizes the request source
func (v *TransactionValidator) SanitizeRequestSource(source string) string {
	if source == "" {
		return "unknown"
	}
	
	source = strings.ToLower(strings.TrimSpace(source))
	if !v.allowedSources[source] {
		return "unknown"
	}
	
	return source
}

// GetAllowedSources returns all allowed request sources
func (v *TransactionValidator) GetAllowedSources() []string {
	sources := make([]string, 0, len(v.allowedSources))
	for source := range v.allowedSources {
		sources = append(sources, source)
	}
	return sources
}