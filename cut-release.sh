#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

function main() {
  # Parse arguments
  if [[ "$#" -ne 1 && "$#" -ne 2 ]]; then
    usage
    exit 1
  fi
  local -r new_version=${1-}
  DRY_RUN=true
  if [[ "${2-}" == "--no-dry-run" ]]; then
    echo "!!! This NOT is a dry run."
    DRY_RUN=false
  else
    echo "This is a dry run."
  fi

  # Get and verify version info
  local -r version_regex="^v(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)\\.(0|[1-9][0-9]*)$"
  if [[ "${new_version}" =~ $version_regex ]]; then
    local -r version_major="${BASH_REMATCH[1]}"
    local -r version_minor="${BASH_REMATCH[2]}"
    local -r version_patch="${BASH_REMATCH[3]}"
  else
    usage
    echo
    echo "!!! You specified an invalid version '${new_version}'."
    exit 1
  fi

  git-delete-tag "${new_version}"
  git-tag "${new_version}"
  git-push-tags
}

function usage() {
  echo "Usage: ${0} <release_version> [--no-dry-run]"
  echo
  echo "<release_version> is the version you want to release,"
  echo "Please see docs/releasing.md for more info."
}

function current-git-commit() {
  git rev-parse --short HEAD
}

function git-delete-tag() {
  local -r new_version="${1}"
  if (git rev-parse ${new_version} >/dev/null 2>&1); then
    echo "Deleting ${new_version} tag."
    if $DRY_RUN; then
      echo "Dry run: would have done git tag -d ${new_version}"
    else
      git tag -d "${new_version}"
    fi
  fi
}

function git-tag() {
  local -r new_version="${1}"
  echo "Tagging ${new_version} at $(current-git-commit)."
  if $DRY_RUN; then
    echo "Dry run: would have done git tag -a -m \"API release ${new_version}\" ${new_version}"
  else
    git tag -a -m "API release ${new_version}" "${new_version}"
  fi
}

function git-push-tags() {
  echo "Pushing tags."
  if $DRY_RUN; then
    echo "Dry run: would have done git push -f --tags"
  else
    git push -f --tags
  fi
}

main "$@"
