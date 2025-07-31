package template

import (
	"os"
	"strings"
)

const (
	TenantLifeAI   = "LIFE AI"
	TenantGenetica = "GENETICA"
)

// getTenantFromTraits tries to extract the tenant name from traits or transient payload.
// Falls back to TENANT_NAME env var or "Unknown".
func getTenantFromTraits(traits map[string]interface{}, transientPayload map[string]interface{}) string {
	var raw string

	// 1. From traits.tenant
	if t, ok := traits["tenant"].(string); ok {
		raw = t
	}

	// 2. From transient payload
	if raw == "" {
		if t, ok := transientPayload["tenant"].(string); ok {
			raw = t
		}
	}

	// 3. From env
	if raw == "" {
		raw = os.Getenv("TENANT_NAME")
	}

	if raw == "" {
		raw = "Unknown"
	}

	return raw
}

// normalizeTenant standardizes tenant display names.
func normalizeTenant(t string) string {
	switch strings.ToLower(t) {
	case "life_ai", "lifeai", "life ai":
		return TenantLifeAI
	case "genetica":
		return TenantGenetica
	default:
		return t
	}
}

// GetNormalizedTenantFromTraits extracts and normalizes the tenant name from traits or transient payload.
func GetNormalizedTenantFromTraits(traits map[string]interface{}, transientPayload map[string]interface{}) string {
	tenant := getTenantFromTraits(traits, transientPayload)
	return normalizeTenant(tenant)
}
