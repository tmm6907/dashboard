#!/bin/bash
# Start postfix service (run in background)
service postfix start

# Start your Go mailserver (keep in foreground)
exec /mail/mailserver
