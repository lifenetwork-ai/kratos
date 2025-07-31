package template

import (
	"fmt"
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

// getLangFromTraits extracts the language from traits, defaulting to "en" if not found.
func getLangFromTraits(traits map[string]interface{}) string {
	if lang, ok := traits["lang"].(string); ok && lang != "" {
		return strings.ToLower(lang)
	}
	return "en"
}

// GetTemplatePathAndGlob constructs the template path and glob pattern based on traits, transient payload, action, status, and template name.
func GetTemplatePathAndGlob(
	traits map[string]interface{},
	transientPayload map[string]interface{},
	action string, // e.g. "registration_code"
	status string, // e.g. "valid" or "invalid"
	templateName string, // e.g. "email.subject"
) (string, string) {
	lang := getLangFromTraits(traits)
	tenant := getTenantFromTraits(traits, transientPayload)

	dir := fmt.Sprintf("%s/%s/%s/%s", tenant, lang, action, status)
	templatePath := fmt.Sprintf("%s/%s.gotmpl", dir, templateName)
	templateGlob := fmt.Sprintf("%s/%s.*", dir, templateName)

	return templatePath, templateGlob
}
