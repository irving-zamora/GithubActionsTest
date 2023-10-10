#!/bin/bash
# This script serves a helper script to run the Tyk bundler tool to create a production-ready plugin bundle.
set -euo pipefail;
echo "Building plugin bundle...";

# Run bundler tool in bundle directory
cd /opt/tyk-gateway/bundle && /opt/tyk-gateway/tyk bundle build -y;

# TODO: also wondering if we should put this zip bundle artifact either back into github or aritfactory?
# talk of artifactory being replaced

# Cleanup
rm /opt/tyk-gateway/bundle/RateLimitingPlugin*.so;

# Exit
echo "Done building plugin bundle.";
exit 0;