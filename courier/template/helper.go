package template

import "os"

// GetTenantFromContext tries to extract the tenant name from identity traits or transient payload.
// Falls back to TENANT_NAME env var or "Unknown".
func GetTenantFromContext(identity map[string]interface{}, transientPayload map[string]interface{}) string {
	// Prefer to get tenant from identity traits
	if traits, ok := identity["traits"].(map[string]interface{}); ok {
		if tenant, ok := traits["tenant"].(string); ok {
			return tenant
		}
	}

	// Fallback to transient payload
	if tenant, ok := transientPayload["tenant"].(string); ok {
		return tenant
	}

	// Fallback to environment
	if fallback := os.Getenv("TENANT_NAME"); fallback != "" {
		return fallback
	}

	return "Unknown"
}
