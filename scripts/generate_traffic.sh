#!/bin/sh

set -eu

BASE_URL="${BASE_URL:-https://movies-api-latest.onrender.com}"
ITERATIONS="${ITERATIONS:-200}"
SLEEP_SECONDS="${SLEEP_SECONDS:-1}"

echo "Generating traffic against ${BASE_URL}"

i=1
while [ "$i" -le "$ITERATIONS" ]; do
  echo "Iteration ${i}/${ITERATIONS}"

  # Successful requests that should produce 2xx metrics.
  curl --silent --output /dev/null "${BASE_URL}/health"
  curl --silent --output /dev/null "${BASE_URL}/movies/search?title=matrix"
  curl --silent --output /dev/null "${BASE_URL}/movies/search?title=batman"
  curl --silent --output /dev/null "${BASE_URL}/movies/550"

  # Intentional client-side errors so the dashboard shows 4xx responses.
  curl --silent --output /dev/null "${BASE_URL}/movies/search"
  curl --silent --output /dev/null "${BASE_URL}/movies/not-a-number"

  # Intentional server-side error so the dashboard shows 5xx responses.
  curl --silent --output /dev/null "${BASE_URL}/debug/error"

  sleep "${SLEEP_SECONDS}"
  i=$((i + 1))
done

echo "Done."
