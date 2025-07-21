// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sms

import (
	"context"
	"encoding/json"
	"os"

	"github.com/ory/kratos/courier/template"
)

type (
	LoginCodeValid struct {
		deps  template.Dependencies
		model *LoginCodeValidModel
	}
	LoginCodeValidModel struct {
		To               string                 `json:"to"`
		LoginCode        string                 `json:"login_code"`
		Identity         map[string]interface{} `json:"identity"`
		RequestURL       string                 `json:"request_url"`
		TransientPayload map[string]interface{} `json:"transient_payload"`
		ExpiresInMinutes int                    `json:"expires_in_minutes"`
	}
)

func NewLoginCodeValid(d template.Dependencies, m *LoginCodeValidModel) *LoginCodeValid {
	return &LoginCodeValid{deps: d, model: m}
}

func (t *LoginCodeValid) PhoneNumber() (string, error) {
	return t.model.To, nil
}

// getTenant extracts the tenant information from the identity traits.
func (t *LoginCodeValid) getTenant() string {
	// Prefer to get tenant from identity traits
	if traits, ok := t.model.Identity["traits"].(map[string]interface{}); ok {
		if tenant, ok := traits["tenant"].(string); ok {
			return tenant
		}
	}
	// Fallback from transient payload if set in the flow
	if tenant, ok := t.model.TransientPayload["tenant"].(string); ok {
		return tenant
	}
	// Fallback to environment variable TENANT_NAME if not set in traits or transient payload
	if fallback := os.Getenv("TENANT_NAME"); fallback != "" {
		return fallback
	}
	return "Unknown"
}

func (t *LoginCodeValid) SMSBody(ctx context.Context) (string, error) {
	data := struct {
		*LoginCodeValidModel
		Tenant string
	}{
		LoginCodeValidModel: t.model,
		Tenant:              t.getTenant(),
	}

	return template.LoadText(
		ctx,
		t.deps,
		os.DirFS(t.deps.CourierConfig().CourierTemplatesRoot(ctx)),
		"login_code/valid/sms.body.gotmpl",
		"login_code/valid/sms.body*",
		data,
		t.deps.CourierConfig().CourierSMSTemplatesLoginCodeValid(ctx).Body.PlainText,
	)
}

func (t *LoginCodeValid) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.model)
}

func (t *LoginCodeValid) TemplateType() template.TemplateType {
	return template.TypeLoginCodeValid
}
