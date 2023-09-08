#!/bin/bash
# This script serves a helper script to run the Tyk bundler tool to create a production-ready plugin bundle.
set -euo pipefail;
echo "Building plugin bundle...";

# Copy custom plugin to bundle directory
cp /opt/tyk-gateway/middleware/RateLimitingPlugin*.so /opt/tyk-gateway/bundle/RateLimitingPlugin.so;

# Run bundler tool in bundle directory
cd /opt/tyk-gateway/bundle && /opt/tyk-gateway/tyk bundle build -y;

# rename bundle zip to include gateway version targeted for and build version
mv bundle.zip RateLimitingPlugin_gw5.0.3_1.0.1.zip

# replace the path and file_list in manifest.json 
# sed ...

# also wondering if we should put this zip bundle artifact either back into github or aritfactory?
# talk of artifactory being replaced

# Cleanup
rm /opt/tyk-gateway/bundle/RateLimitingPlugin*.so;

# Exit
echo "Done building plugin bundle.";
exit 0;
