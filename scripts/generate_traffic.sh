#!/bin/sh

set -eu

BASE_URL="${BASE_URL:-http://localhost:8080}"
ITERATIONS="${ITERATIONS:-5}"
SLEEP_SECONDS="${SLEEP_SECONDS:-1}"

echo "Generating traffic against ${BASE_URL}"

i=1
while [ "$i" -le "$ITERATIONS" ]; do
  echo "Iteration ${i}/${ITERATIONS}"

  curl --silent --output /dev/null "${BASE_URL}/health"
  curl --silent --output /dev/null "${BASE_URL}/movies/search?title=matrix"
  curl --silent --output /dev/null "${BASE_URL}/movies/search?title=batman"
  curl --silent --output /dev/null "${BASE_URL}/movies/550"

  # Intentional client-side errors so the dashboard shows non-200 responses.
  curl --silent --output /dev/null "${BASE_URL}/movies/search"
  curl --silent --output /dev/null "${BASE_URL}/movies/not-a-number"

  sleep "${SLEEP_SECONDS}"
  i=$((i + 1))
done

echo "Done."
