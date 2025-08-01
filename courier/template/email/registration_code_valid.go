// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package email

import (
	"context"
	"encoding/json"
	"os"
	"strings"

	"github.com/ory/kratos/courier/template"
)

type (
	RegistrationCodeValid struct {
		deps  template.Dependencies
		model *RegistrationCodeValidModel
	}
	RegistrationCodeValidModel struct {
		To               string                 `json:"to"`
		Traits           map[string]interface{} `json:"traits"`
		RegistrationCode string                 `json:"registration_code"`
		RequestURL       string                 `json:"request_url"`
		TransientPayload map[string]interface{} `json:"transient_payload"`
		ExpiresInMinutes int                    `json:"expires_in_minutes"`
		Tenant           string                 `json:"tenant"`
	}
)

func NewRegistrationCodeValid(d template.Dependencies, m *RegistrationCodeValidModel) *RegistrationCodeValid {
	m.Tenant = template.GetNormalizedTenantFromTraits(m.Traits, m.TransientPayload)

	return &RegistrationCodeValid{deps: d, model: m}
}

func (t *RegistrationCodeValid) EmailRecipient() (string, error) {
	return t.model.To, nil
}

func (t *RegistrationCodeValid) EmailSubject(ctx context.Context) (string, error) {
	templatePath, templateGlob := template.GetTemplatePathAndGlob(
		t.model.Traits,
		t.model.TransientPayload,
		"registration_code",
		"valid",
		"email.subject",
	)

	subject, err := template.LoadText(
		ctx,
		t.deps,
		os.DirFS(t.deps.CourierConfig().CourierTemplatesRoot(ctx)),
		templatePath,
		templateGlob,
		t.model,
		t.deps.CourierConfig().CourierTemplatesRegistrationCodeValid(ctx).Subject,
	)

	return strings.TrimSpace(subject), err
}

func (t *RegistrationCodeValid) EmailBody(ctx context.Context) (string, error) {
	traits := template.GetTraitsFromIdentity(t.model.Traits)
	templatePath, templateGlob := template.GetTemplatePathAndGlob(
		traits,
		t.model.TransientPayload,
		"registration_code",
		"valid",
		"email.body",
	)

	return template.LoadHTML(
		ctx,
		t.deps,
		os.DirFS(t.deps.CourierConfig().CourierTemplatesRoot(ctx)),
		templatePath,
		templateGlob,
		t.model,
		t.deps.CourierConfig().CourierTemplatesRegistrationCodeValid(ctx).Body.HTML,
	)
}

func (t *RegistrationCodeValid) EmailBodyPlaintext(ctx context.Context) (string, error) {
	traits := template.GetTraitsFromIdentity(t.model.Traits)
	templatePath, templateGlob := template.GetTemplatePathAndGlob(
		traits,
		t.model.TransientPayload,
		"registration_code",
		"valid",
		"email.body.plaintext",
	)

	return template.LoadText(
		ctx,
		t.deps,
		os.DirFS(t.deps.CourierConfig().CourierTemplatesRoot(ctx)),
		templatePath,
		templateGlob,
		t.model,
		t.deps.CourierConfig().CourierTemplatesRegistrationCodeValid(ctx).Body.PlainText,
	)
}

func (t *RegistrationCodeValid) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.model)
}

func (t *RegistrationCodeValid) TemplateType() template.TemplateType {
	return template.TypeRegistrationCodeValid
}
