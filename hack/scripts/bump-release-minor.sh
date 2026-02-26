#!/usr/bin/env bash

# Copyright AppsCode Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

if [[ $# -ne 1 ]]; then
    echo "Usage: $0 <product>" >&2
    echo "Supported products: voyager kubevault kubedb ace stash kubestash" >&2
    exit 1
fi

product="${1,,}"

case "${product}" in
    voyager | voyager_create_release.go)
        target_basename="voyager_create_release.go"
        ;;
    kubevault | kubevault_create_release.go)
        target_basename="kubevault_create_release.go"
        ;;
    kubedb | kubedb_create_release.go)
        target_basename="kubedb_create_release.go"
        ;;
    ace | ace_create_release.go)
        target_basename="ace_create_release.go"
        ;;
    stash | stash_create_release.go)
        target_basename="stash_create_release.go"
        ;;
    kubestash | kubestash_create_release.go)
        target_basename="kubestash_create_release.go"
        ;;
    *)
        echo "Unsupported product: ${1}" >&2
        echo "Supported products: voyager kubevault kubedb ace stash kubestash" >&2
        exit 1
        ;;
esac

TARGET_FILE="${REPO_ROOT}/cmds/${target_basename}"

if [[ ! -f "${TARGET_FILE}" ]]; then
    echo "Target file not found: ${TARGET_FILE}" >&2
    exit 1
fi

perl -i.bak -pe '
  s{(?<![0-9A-Za-z])(v?)(\d+)\.(\d+)\.(\d+)((?:-[0-9A-Za-z.-]+)?(?:\+[0-9A-Za-z.-]+)?)?(?![0-9A-Za-z])}{
    my ($prefix, $major, $minor, $patch, $suffix) = ($1, $2, $3, $4, $5 // "");
    $prefix . $major . "." . ($minor + 1) . ".0" . $suffix;
  }gex;
' "${TARGET_FILE}"

rm -f "${TARGET_FILE}.bak"

echo "Updated semver minor versions in ${TARGET_FILE}"
