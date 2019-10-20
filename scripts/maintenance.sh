#!/usr/bin/env bash
set -euo pipefail

# Update root (Golang) .gitignore
curl https://www.gitignore.io/api/go > .gitignore

# Update web (React) .gitignore
echo "build/" > web/.gitignore
curl https://www.gitignore.io/api/react >> web/.gitignore

# Build .dockerignore from **/.gitignore
function prepend() { while read line; do echo "${1}${line}"; done; }
if [ -f .dockerignore ]; then
  rm .dockerignore
fi
for FILE in $(find . -name .gitignore); do
  DIRNAME=$(dirname "${FILE}" | sed 's/^[\.\/]*//')
  echo $DIRNAME
  cat "${FILE}" | grep -v ^# | awk NF | prepend "${DIRNAME:+${DIRNAME}/}" >> .dockerignore
done
