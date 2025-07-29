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

func (t *LoginCodeValid) EmailRecipient() (string, error) {
	return t.model.To, nil
}

func (t *LoginCodeValid) EmailSubject(ctx context.Context) (string, error) {
	data := struct {
		*LoginCodeValidModel
		Tenant string
	}{
		LoginCodeValidModel: t.model,
		Tenant:              template.GetTenantFromContext(t.model.Identity, t.model.TransientPayload),
	}

	subject, err := template.LoadText(
		ctx,
		t.deps,
		os.DirFS(t.deps.CourierConfig().CourierTemplatesRoot(ctx)),
		"login_code/valid/email.subject.gotmpl",
		"login_code/valid/email.subject*",
		data,
		t.deps.CourierConfig().CourierTemplatesLoginCodeValid(ctx).Subject,
	)

	return strings.TrimSpace(subject), err
}

func (t *LoginCodeValid) EmailBody(ctx context.Context) (string, error) {
	data := struct {
		*LoginCodeValidModel
		Tenant string
	}{
		LoginCodeValidModel: t.model,
		Tenant:              template.GetTenantFromContext(t.model.Identity, t.model.TransientPayload),
	}

	return template.LoadHTML(
		ctx,
		t.deps,
		os.DirFS(t.deps.CourierConfig().CourierTemplatesRoot(ctx)),
		"login_code/valid/email.body.gotmpl",
		"login_code/valid/email.body*",
		data,
		t.deps.CourierConfig().CourierTemplatesLoginCodeValid(ctx).Body.HTML,
	)
}

func (t *LoginCodeValid) EmailBodyPlaintext(ctx context.Context) (string, error) {
	data := struct {
		*LoginCodeValidModel
		Tenant string
	}{
		LoginCodeValidModel: t.model,
		Tenant:              template.GetTenantFromContext(t.model.Identity, t.model.TransientPayload),
	}

	return template.LoadText(
		ctx,
		t.deps,
		os.DirFS(t.deps.CourierConfig().CourierTemplatesRoot(ctx)),
		"login_code/valid/email.body.plaintext.gotmpl",
		"login_code/valid/email.body.plaintext*",
		data,
		t.deps.CourierConfig().CourierTemplatesLoginCodeValid(ctx).Body.PlainText,
	)
}

func (t *LoginCodeValid) MarshalJSON() ([]byte, error) {
	return json.Marshal(t.model)
}

func (t *LoginCodeValid) TemplateType() template.TemplateType {
	return template.TypeLoginCodeValid
}
