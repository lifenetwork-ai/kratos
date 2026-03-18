package template

import (
	"fmt"
	"os"
	"strings"

	"github.com/ory/x/logrusx"
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
// It returns the primary path/glob for the requested language and a fallback path/glob using "en".
// If the requested language is already "en", fallback values are empty.
// It returns the primary path/glob for the requested language and a fallback path/glob using "en".
// If the requested language is already "en", fallback values are empty.
func GetTemplatePathAndGlob(
	traits map[string]interface{},
	transientPayload map[string]interface{},
	action string, // e.g. "registration_code"
	status string, // e.g. "valid" or "invalid"
	templateName string, // e.g. "email.subject"
) (templatePath, templateGlob, fallbackPath, fallbackGlob string) {
	return GetTemplatePathAndGlobWithLogger(nil, traits, transientPayload, action, status, templateName)
}

// GetTemplatePathAndGlobWithLogger constructs the template path and glob pattern with debug logging.
func GetTemplatePathAndGlobWithLogger(
	logger *logrusx.Logger,
	traits map[string]interface{},
	transientPayload map[string]interface{},
	action string, // e.g. "registration_code"
	status string, // e.g. "valid" or "invalid"
	templateName string, // e.g. "email.subject"
) (templatePath, templateGlob, fallbackPath, fallbackGlob string) {
	lang := getLangFromTraits(traits)
	tenant := getTenantFromTraits(traits, transientPayload)

	if logger != nil {
		logger.
			WithField("traits", traits).
			WithField("transient_payload", transientPayload).
			WithField("extracted_lang", lang).
			WithField("extracted_tenant", tenant).
			WithField("action", action).
			WithField("status", status).
			WithField("template_name", templateName).
			Debug("Resolving template path from traits")
	}

	dir := fmt.Sprintf("%s/%s/%s/%s", tenant, lang, action, status)
	templatePath = fmt.Sprintf("%s/%s.gotmpl", dir, templateName)
	templateGlob = fmt.Sprintf("%s/%s.*", dir, templateName)
	templatePath = fmt.Sprintf("%s/%s.gotmpl", dir, templateName)
	templateGlob = fmt.Sprintf("%s/%s.*", dir, templateName)

	if lang != "en" {
		fallbackDir := fmt.Sprintf("%s/%s/%s/%s", tenant, "en", action, status)
		fallbackPath = fmt.Sprintf("%s/%s.gotmpl", fallbackDir, templateName)
		fallbackGlob = fmt.Sprintf("%s/%s.*", fallbackDir, templateName)
	}

	if logger != nil {
		logger.
			WithField("template_path", templatePath).
			WithField("template_glob", templateGlob).
			WithField("fallback_path", fallbackPath).
			WithField("fallback_glob", fallbackGlob).
			Debug("Template path resolved")
	}

	return
}

// GetTraitsFromIdentity safely extracts the traits object from a Kratos identity-like map.
// Returns nil if not found or not a valid map.
func GetTraitsFromIdentity(identity map[string]interface{}) map[string]interface{} {
	if identity == nil {
		return nil
	}
	if traits, ok := identity["traits"].(map[string]interface{}); ok {
		return traits
	}
	return nil
}
