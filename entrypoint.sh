#!/bin/bash
set -e

# Set default values if not provided
export SELFSERVICE_URL=${SELFSERVICE_URL:-http://localhost:3000}
export DOMAIN=${DOMAIN:-localhost}
export SMTP_CONNECTION_URI=${SMTP_CONNECTION_URI:-smtps://user:pass@smtp.gmail.com:465/?skip_ssl_verify=false&legacy_ssl=true}
export SMTP_FROM_ADDRESS=${SMTP_FROM_ADDRESS:-no-reply@lifenetwork.ai}
export SMTP_FROM_NAME=${SMTP_FROM_NAME:-LIFE AI}
export SMS_WEBHOOK_URL=${SMS_WEBHOOK_URL:-https://webhook.example}
export ENABLE_COURIER=${ENABLE_COURIER:-false}
export KRATOS_CONFIG_TEMPLATE=${KRATOS_CONFIG_TEMPLATE:-/etc/config/kratos_template.yml}
export KRATOS_RENDERED_CONFIG=${KRATOS_RENDERED_CONFIG:-/etc/config/kratos.yml}

CONFIG_ARGS=()
if [ -n "$KRATOS_CONFIG_FILES" ]; then
  IFS=',' read -ra CONFIG_FILES <<< "$KRATOS_CONFIG_FILES"
  for file in "${CONFIG_FILES[@]}"; do
    CONFIG_ARGS+=("-c" "$file")
  done
else
  echo "Rendering ${KRATOS_CONFIG_TEMPLATE} from env vars..."
  envsubst < "$KRATOS_CONFIG_TEMPLATE" > "$KRATOS_RENDERED_CONFIG"
  CONFIG_ARGS=("-c" "$KRATOS_RENDERED_CONFIG")
fi

# Decide whether to enable courier worker.
if [ "$ENABLE_COURIER" = "true" ]; then
  KRATOS_CMD=(kratos serve --watch-courier "${CONFIG_ARGS[@]}")
else
  KRATOS_CMD=(kratos serve "${CONFIG_ARGS[@]}")
fi

echo "Starting ORY Kratos with command: ${KRATOS_CMD[*]}"
exec "${KRATOS_CMD[@]}"
