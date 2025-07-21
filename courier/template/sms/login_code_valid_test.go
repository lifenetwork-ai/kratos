// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package sms_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/kratos/courier/template/sms"
	"github.com/ory/kratos/internal"
)

func TestNewLoginCodeValid(t *testing.T) {
	_, reg := internal.NewFastRegistryWithMocks(t)

	const (
		expectedPhone  = "+12345678901"
		otp            = "012345"
		expectedTenant = "tenant"
	)

	tpl := sms.NewLoginCodeValid(reg, &sms.LoginCodeValidModel{
		To:               expectedPhone,
		LoginCode:        otp,
		ExpiresInMinutes: 0,
		Identity: map[string]interface{}{
			"traits": map[string]interface{}{
				"tenant": expectedTenant,
			},
		},
	})

	// Update expected body to match the new template
	expectedBody := fmt.Sprintf("[%s] Your login code is: %s\nIt expires in 0 minutes.\n", expectedTenant, otp)

	actualBody, err := tpl.SMSBody(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expectedBody, actualBody)

	actualPhone, err := tpl.PhoneNumber()
	require.NoError(t, err)
	assert.Equal(t, expectedPhone, actualPhone)
}
