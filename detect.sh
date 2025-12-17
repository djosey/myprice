#!/usr/bin/env bash
set -euo pipefail

IMG="/Users/donaldjosey/developer/myprice/PXL_20251206_222017258.jpg"
OUT="/Users/donaldjosey/developer/myprice/textract_output.json"

echo "Running Textract on: $IMG"
echo "Output will be saved to: $OUT"

aws textract detect-document-text \
  --region us-east-1 \
  --document "Bytes=$(base64 < "$IMG" | tr -d '\n')" \
  > "$OUT"

echo "Done. Written to $OUT"
zsh: event not found: /usr/bin/env
