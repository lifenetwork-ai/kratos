#!/bin/bash
set -e

# Set default values if not provided
export SELFSERVICE_URL=${SELFSERVICE_URL:-http://localhost:3000}
export DOMAIN=${DOMAIN:-localhost}
export SMTP_CONNECTION_URI=${SMTP_CONNECTION_URI:-smtps://user:pass@smtp.gmail.com:465/?skip_ssl_verify=false&legacy_ssl=true}
export SMTP_FROM_ADDRESS=${SMTP_FROM_ADDRESS:-no-reply@lifenetwork.ai}
export SMS_WEBHOOK_URL=${SMS_WEBHOOK_URL:-https://webhook.example}
export ENABLE_COURIER=${ENABLE_COURIER:-false}

# Decide whether to enable courier worker
if [ "$ENABLE_COURIER" = "true" ]; then
  KRATOS_CMD="kratos serve --watch-courier -c /etc/config/kratos.yml"
else
  KRATOS_CMD="kratos serve -c /etc/config/kratos.yml"
fi

echo "Rendering kratos.yml from env vars..."
envsubst < /etc/config/kratos_template.yml > /etc/config/kratos.yml

echo "Starting ORY Kratos with command: $KRATOS_CMD"
exec $KRATOS_CMD
