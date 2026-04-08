// Copyright © 2023 Ory Corp
// SPDX-License-Identifier: Apache-2.0

package email_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ory/kratos/courier/template"
	"github.com/ory/kratos/courier/template/email"
	"github.com/ory/kratos/courier/template/testhelpers"
	"github.com/ory/kratos/internal"
)

func TestVerifyCodeValid(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	t.Run("test=with courier templates directory", func(t *testing.T) {
		_, reg := internal.NewFastRegistryWithMocks(t)
		tpl := email.NewVerificationCodeValid(reg, &email.VerificationCodeValidModel{})

		testhelpers.TestRendered(t, ctx, tpl)
	})

	t.Run("test=with remote resources", func(t *testing.T) {
		testhelpers.TestRemoteTemplates(t, "../courier/builtin/templates/verification_code/valid", template.TypeVerificationCodeValid)
	})
}

func TestVerifyCodeValidLifeAICommitUnityUsesGenericActionCopy(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	t.Cleanup(cancel)

	for _, tc := range []struct {
		name             string
		lang             string
		profileName      string
		expectedGreeting string
		expectedAction   string
		expectedValidity string
		expectedWarning  string
		expectedTeam     string
	}{
		{
			name:             "en",
			lang:             "en",
			profileName:      "Alex",
			expectedGreeting: "Dear Mr./Ms. Alex,",
			expectedAction:   "We have received a request to commit Life Points from Main wallet.",
			expectedValidity: "This OTP will be valid for 2 minutes.",
			expectedWarning:  "If you did NOT initiate this request",
			expectedTeam:     "LIFE AI Team",
		},
		{
			name:             "vi",
			lang:             "vi",
			expectedGreeting: "Kính gửi Quý khách,",
			expectedAction:   "Chúng tôi vừa nhận được yêu cầu commit Life Points từ Ví chính.",
			expectedValidity: "Mã OTP chỉ có hiệu lực trong 2 phút tới.",
			expectedWarning:  "Nếu bạn KHÔNG phải là người thực hiện yêu cầu này",
			expectedTeam:     "Đội ngũ LIFE AI",
		},
		{
			name:             "other_locale_uses_english",
			lang:             "id",
			profileName:      "Alex",
			expectedGreeting: "Dear Mr./Ms. Alex,",
			expectedAction:   "We have received a request to commit Life Points from Main wallet.",
			expectedValidity: "This OTP will be valid for 2 minutes.",
			expectedWarning:  "If you did NOT initiate this request",
			expectedTeam:     "LIFE AI Team",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, reg := internal.NewVeryFastRegistryWithoutDB(t)

			transientPayload := map[string]interface{}{
				"verify_type": "commit_unity",
			}
			if tc.profileName != "" {
				transientPayload["profile_name"] = tc.profileName
			}

			tpl := email.NewVerificationCodeValid(reg, &email.VerificationCodeValidModel{
				VerificationCode: "123456",
				ExpiresInMinutes: 2,
				Identity: map[string]interface{}{
					"traits": map[string]interface{}{
						"tenant": "life_ai",
						"lang":   tc.lang,
					},
				},
				TransientPayload: transientPayload,
			})

			plain, err := tpl.EmailBodyPlaintext(ctx)
			require.NoError(t, err)
			assert.Contains(t, plain, tc.expectedGreeting)
			assert.Contains(t, plain, tc.expectedAction)
			assert.Contains(t, plain, "123456")
			assert.Contains(t, plain, tc.expectedValidity)
			assert.Contains(t, plain, tc.expectedWarning)
			assert.Contains(t, plain, tc.expectedTeam)
			assert.NotContains(t, plain, "Transaction details")
			assert.NotContains(t, plain, "Commitment amount")
			assert.NotContains(t, plain, "Do not share this OTP")

			html, err := tpl.EmailBody(ctx)
			require.NoError(t, err)
			assert.Contains(t, html, tc.expectedAction)
			assert.Contains(t, html, "123456")
			assert.NotContains(t, html, "Transaction details")
			assert.NotContains(t, html, "Commitment amount")
		})
	}
}
