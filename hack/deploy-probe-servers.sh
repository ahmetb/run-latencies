#!/usr/bin/env bash

set -euo pipefail
[[ -n "${DEBUG:-}" ]] && set -x

SCRIPTDIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "${SCRIPTDIR}/.."

regions=(\
    arn1 \
    bom1 \
    bru1 \
    cdg1 \
    cle1 \
    dub1 \
    gru1 \
    hnd1 \
    iad1 \
    icn1 \
    lhr1 \
    pdx1 \
    sfo1 \
    sin1 \
    syd1 \
)

clr=$(tput sgr0)
y=$(tput setaf 3)

for r in "${regions[@]}"; do
  echo "$y-> deploying to $r...$clr"
  now --name prober-"$r" --regions="$r"
done
