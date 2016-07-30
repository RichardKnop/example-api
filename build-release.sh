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

  read -p "Container name (default: example-api): " container_name
  [ -z "$container_name" ] && container_name="example-api"

  local -r github=`git config --get remote.origin.url`
  declare -r temp_dir=$(mktemp -d "/tmp/${container_name}-${new_version}.XXXX")
  local -r tag="${container_name}:${new_version}"
  local -r registry_tag="registry.local/${container_name}:${new_version}"

  git-clone "${github}" "${temp_dir}"
  git-checkout "${new_version}" "${temp_dir}"
  docker-build "${tag}" "${temp_dir}"
  docker-tag "${tag}" "${registry_tag}"
  docker-push "${registry_tag}"
  docker-cleanup "${registry_tag}"
  rm -Rf "${temp_dir}"
}

function usage() {
  echo "Usage: ${0} <release_version> [--no-dry-run]"
  echo
  echo "<release_version> is the version you want to release,"
  echo "Please see docs/releasing.md for more info."
}

function git-clone() {
  local -r github="${1}"
  local -r dest="${2}"
  echo "Cloning from '${github}' to '${dest}'..."
  if $DRY_RUN; then
    echo "Dry run: would have done git clone ${github} ${dest}"
  else
    git clone "${github}" "${dest}"
  fi
}

function git-checkout() {
  local -r tag="${1}"
  local -r dir="${2}"
  echo "Checking out tag '${tag}'..."
  if $DRY_RUN; then
    echo "Dry run: would have done cd ${dir} ; git checkout -b ${tag} ${tag}"
  else
    (cd ${dir} ; git checkout -b "${tag}" "${tag}")
  fi
}

function docker-build() {
  local -r tag="${1}"
  local -r dir="${2}"
  echo "Building docker container '${tag}'..."
  if $DRY_RUN; then
    echo "Dry run: would have done docker build -t ${tag} ${dir}"
  else
    docker build -t "${tag}" "${dir}"
  fi
}

function docker-tag() {
  local -r tag="${1}"
  local -r registry_tag="${2}"
  echo "Tagging as '${tag}' as '${registry_tag}'..."
  if $DRY_RUN; then
    echo "Dry run: would have done docker tag ${tag} ${registry_tag}"
  else
    docker tag "${tag}" "${registry_tag}"
  fi
}

function docker-push() {
  local -r registry_tag="${1}"
  echo "Pushing '${registry_tag}' to registry..."
  if $DRY_RUN; then
    echo "Dry run: would have done docker push ${registry_tag}"
  else
    docker push "${registry_tag}"
  fi
}

function docker-cleanup() {
  local -r registry_tag="${1}"
  echo "Docker cleanup..."
  if $DRY_RUN; then
    echo "Dry run: would have done docker rmi ${registry_tag}"
  else
    docker rmi "${registry_tag}"
    docker rmi -f $(docker images -q -f "dangling=true") || true
  fi
}

main "$@"
