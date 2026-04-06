# Template Routing

This repo contains LIFE AI courier templates in:

```text
courier/template/courier/builtin/templates/life_ai/<lang>/<action>/<status>/<template>.gotmpl
```

For verification-code emails, the active templates are under:

```text
courier/template/courier/builtin/templates/life_ai/<lang>/verification_code/valid/
```

## Template Path Selection

Template path selection happens in `courier/template/helper.go`.

- `lang` comes from `identity.traits.lang`; if it is empty, the default is `en`.
- `tenant` comes from `identity.traits.tenant`, then `transient_payload.tenant`, then `TENANT_NAME`, then `Unknown`.
- Template path selection uses the raw tenant value as the directory name. For LIFE AI templates, that directory is `life_ai`.
- The primary template path is `<tenant>/<lang>/<action>/<status>/<template>.gotmpl`.
- If `lang` is not `en`, the loader also prepares an English fallback path.

The rendered `.Tenant` display value is normalized separately for known values such as `life_ai`, `lifeai`, and `life ai`.

## Verification `transient_payload`

The code verification strategy copies the request's `transient_payload` onto the verification flow. The sender later parses it as a JSON object and passes it to the email template model as `.TransientPayload`.

`transient_payload` is not validated against a schema in this path. A malformed JSON payload fails when it is parsed, but arbitrary object keys and values are otherwise passed through to templates.

## Verification `verify_type` Routing

For LIFE AI verification-code templates, `transient_payload.verify_type` is used only as a template-copy selector. It is not a validation enum and does not change whether the OTP is generated or sent.

Known verification types currently used by the LIFE AI verification-code templates:

- `commit_lifeai`
- `commit_unity`
- `withdraw_main_wallet`
- `withdraw_unity`
- `transfer_commit_unity`
- `pay_loan_transfer`
- `pay_loan`

If `verify_type` is unknown, for example:

```json
{
  "verify_type": "abcd"
}
```

the verification flow still proceeds. The templates do not match any special `verify_type` branch, so they fall through to generic content:

- Subject uses `transient_payload.subject` if present.
- If no subject override exists, subject falls back to `{{ .Tenant }} OTP - Expires in {{ .ExpiresInMinutes }} minutes`.
- Body uses the `action_description` branch if `transient_payload.action_description` is present.
- If no `action_description` exists, body falls back to the generic verification request copy.
- Greeting still uses `transient_payload.profile_name` when present; otherwise it uses the generic greeting.

Use `verify_type` only when the requested operation needs one of the predefined message variants. For custom one-off wording, prefer explicit template payload fields such as `subject` and `action_description`.
