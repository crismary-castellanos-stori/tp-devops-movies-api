#!/bin/sh

set -eu

BASE_URL="${BASE_URL:-https://movies-api-latest.onrender.com}"
ITERATIONS="${ITERATIONS:-120}"
SLEEP_SECONDS="${SLEEP_SECONDS:-1}"
SUCCESS_BURST="${SUCCESS_BURST:-12}"
LATENCY_SPIKE_EVERY="${LATENCY_SPIKE_EVERY:-10}"
LATENCY_MS="${LATENCY_MS:-1200}"
INCLUDE_ERRORS="${INCLUDE_ERRORS:-1}"

echo "Generating traffic against ${BASE_URL}"

i=1
while [ "$i" -le "$ITERATIONS" ]; do
  echo "Iteration ${i}/${ITERATIONS}"

  burst=1
  while [ "$burst" -le "$SUCCESS_BURST" ]; do
    curl --silent --output /dev/null "${BASE_URL}/health"
    curl --silent --output /dev/null "${BASE_URL}/movies/search?title=matrix"
    curl --silent --output /dev/null "${BASE_URL}/movies/search?title=batman"
    curl --silent --output /dev/null "${BASE_URL}/movies/550"
    burst=$((burst + 1))
  done

  # Add a slower request every few iterations to create visible latency spikes.
  if [ $((i % LATENCY_SPIKE_EVERY)) -eq 0 ]; then
    curl --silent --output /dev/null "${BASE_URL}/debug/slow?ms=${LATENCY_MS}"
  fi

  # Keep 2xx dominant, but inject a small amount of 4xx and 5xx for dashboards.
  if [ "${INCLUDE_ERRORS}" = "1" ]; then
    if [ $((i % 6)) -eq 0 ]; then
      curl --silent --output /dev/null "${BASE_URL}/movies/search"
      curl --silent --output /dev/null "${BASE_URL}/movies/not-a-number"
    fi

    if [ $((i % 15)) -eq 0 ]; then
      curl --silent --output /dev/null "${BASE_URL}/debug/error"
    fi
  fi

  sleep "${SLEEP_SECONDS}"
  i=$((i + 1))
done

echo "Done."
