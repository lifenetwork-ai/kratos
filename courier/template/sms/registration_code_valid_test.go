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

func TestNewRegistrationCodeValid(t *testing.T) {
	_, reg := internal.NewFastRegistryWithMocks(t)

	const (
		expectedPhone  = "+12345678901"
		otp            = "012345"
		expectedTenant = "tenant"
	)

	tpl := sms.NewRegistrationCodeValid(reg, &sms.RegistrationCodeValidModel{
		To:               expectedPhone,
		RegistrationCode: otp,
		ExpiresInMinutes: 0,
		Identity: map[string]interface{}{
			"traits": map[string]interface{}{
				"tenant": expectedTenant,
			},
		},
	})

	expectedBody := fmt.Sprintf("[%s] Your registration code is: %s\nIt expires in 0 minutes.\n", expectedTenant, otp)

	actualBody, err := tpl.SMSBody(context.Background())
	require.NoError(t, err)
	assert.Equal(t, expectedBody, actualBody)

	actualPhone, err := tpl.PhoneNumber()
	require.NoError(t, err)
	assert.Equal(t, expectedPhone, actualPhone)
}
