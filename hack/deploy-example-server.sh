#!/usr/bin/env bash

set -euo pipefail
[[ -n "${DEBUG:-}" ]] && set -x

SCRIPTDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPTDIR}/.."

tag="$(git describe --always --dirty)"
project="${PROJECT:-ahmetb-samples-playground}"
app_name="${APP_NAME:-example-server}"

function regions() {
  local access_token
  access_token="$(gcloud auth print-access-token -q)"

  curl -fsS -H "Authorization: Bearer ${access_token}" \
    https://run.googleapis.com/v1alpha1/projects/"${project}"/locations |
    jq -r '.locations[].locationId'
}

function image_name() {
  echo "gcr.io/${project}/${app_name}:${tag}"
}

function build_image() {
  local img push_log
  img="$(image_name)"
  push_log="$(mktemp)"
  trap 'rm -- "${push_log}"' RETURN

  docker build --quiet --tag "${img}" ./cmd/example_server 1>/dev/null
  set +e
  docker push "${img}" &> "${push_log}"
  ec=$?
  set -e
  if [ $ec -ne 0 ]; then
    echo -n >&2 "$(tput setaf 1)Error pushing image: $(tput sgr0)"
    cat "${push_log}"
    exit $ec
  fi
  set -e
}

function service_name() {
  local region
  region="$1"
  echo -n "${app_name}-${region}"
}

function deploy_to_cloud_run() {
  local image region deploy_log
  image="$1"
  region="$2"
  deploy_log="$(mktemp)"
  trap 'rm -- "${deploy_log}"' RETURN

  set +e
  CLOUDSDK_CORE_DISABLE_PROMPTS=1 gcloud beta run deploy -q \
    --platform=managed \
    --region="${region}" \
    --allow-unauthenticated \
    --image="${image}" \
    "$(service_name "${region}")" 2>"${deploy_log}"
  ec=$?
  set -e
  if [ $ec -ne 0 ]; then
    echo -n >&2 "$(tput setaf 1)Error deploying to Cloud Run: $(tput sgr0)"
    cat "${deploy_log}"
    exit $ec
  fi
}

function cloud_run_url() {
  local region
  region="${1}"
  CLOUDSDK_CORE_DISABLE_PROMPTS=1 gcloud beta run services describe -q \
    --format=get\(status.url\) \
    --platform=managed \
    --region="${region}" \
    "$(service_name "${region}")"
}

build_image "$(image_name)"
for region in $(regions); do
  echo >&2 "Deploying to $(tput setaf 2)$region$(tput sgr0)..."
  deploy_to_cloud_run "$(image_name)" "${region}"
  cloud_run_url "${region}"
done
