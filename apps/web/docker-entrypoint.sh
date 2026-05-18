#!/bin/sh
# apps/web/docker-entrypoint.sh

set -e

echo "Starting Nginx..."

# Replace API URL if environment variable is set
if [ ! -z "$VITE_API_URL" ]; then
    echo "Setting API URL to: $VITE_API_URL"
    
    # Find and replace API URL in JS files
    find /usr/share/nginx/html -type f -name "*.js" -exec sed -i "s|VITE_API_URL_PLACEHOLDER|$VITE_API_URL|g" {} \;
fi

# Execute the main command
exec "$@"