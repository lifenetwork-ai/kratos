#!/bin/bash
set -e

# Set default values if not provided
export SELFSERVICE_URL=${SELFSERVICE_URL:-http://localhost:3000}
export DOMAIN=${DOMAIN:-localhost}
export SMTP_CONNECTION_URI=${SMTP_CONNECTION_URI:-smtps://user:pass@smtp.gmail.com:465/?skip_ssl_verify=false&legacy_ssl=true}
export SMTP_FROM_ADDRESS=${SMTP_FROM_ADDRESS:-no-reply@ory.kratos.sh}
export TWILIO_URL=${TWILIO_URL:-https://api.twilio.com/2010-04-01/Accounts/XXX/Messages.json}
export TWILIO_USER=${TWILIO_USER:-default_user}
export TWILIO_PASS=${TWILIO_PASS:-default_pass}
export TWILIO_BODY=${TWILIO_BODY:-base64://ZnVuY3Rpb24oY3R4KSB7CiAgVG86IGN0eC5yZWNpcGllbnQsCiAgQm9keTogY3R4LmJvZHksCn0=}


echo "Rendering kratos.yml from env vars..."
envsubst < /etc/config/kratos_template.yml > /etc/config/kratos.yml

echo "Starting ORY Kratos with rendered config..."
exec /usr/bin/kratos serve -c /etc/config/kratos.yml