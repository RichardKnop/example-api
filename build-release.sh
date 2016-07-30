#!/bin/bash

set -o errexit
set -o nounset
set -o pipefail

function main() {
  # Parse arguments
  if [[ "$#" -ne 1 && "$#" -ne 2 && "$#" -ne 3 ]]; then
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

  INTERACTIVE=true
  if [[ "${3-}" == "-y" ]]; then
    INTERACTIVE=false
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

  # Interactive questions
  local container_name="example-api"
  if [ $INTERACTIVE == true ]; then
    read -p "Container name (default: example-api): " container_name
    [ -z "$container_name" ] && container_name="example-api"
  fi

  local -r github=`git config --get remote.origin.url`
  declare -r temp_dir=$(mktemp -d "/tmp/${container_name}-${new_version}.XXXX")
  local -r version_tag="registry.local/${container_name}:${new_version}"
  local -r latest_tag="registry.local/${container_name}:latest"

  git-clone "${github}" "${temp_dir}"
  git-checkout "${new_version}" "${temp_dir}"
  docker-build "${version_tag}" "${temp_dir}"
  docker-tag "${version_tag}" "${latest_tag}"
  docker-push "${version_tag}"
  docker-cleanup "${version_tag}" "${registry_tag}"
  rm -Rf "${temp_dir}"
}

function usage() {
  echo "Usage: ${0} <release_version> [--no-dry-run] [-y]"
  echo
  echo "<release_version> is the version you want to release,"
  echo "--no-dry-run flag turns own the dry run mode and executes real commands"
  echo "-y flag turns own interactive mode (no questions asked, uses default values)"
  echo
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
  local -r version_tag="${1}"
  local -r dir="${2}"
  echo "Building docker container '${version_tag}'..."
  if $DRY_RUN; then
    echo "Dry run: would have done docker build -t ${version_tag} ${dir}"
  else
    docker build -t "${version_tag}" "${dir}"
  fi
}

function docker-tag() {
  local -r version_tag="${1}"
  local -r registry_tag="${2}"
  echo "Tagging as '${version_tag}' as '${registry_tag}'..."
  if $DRY_RUN; then
    echo "Dry run: would have done docker tag ${version_tag} ${registry_tag}"
  else
    docker tag "${version_tag}" "${registry_tag}"
  fi
}

function docker-push() {
  local -r version_tag="${1}"
  echo "Pushing '${version_tag}' to registry..."
  if $DRY_RUN; then
    echo "Dry run: would have done docker push ${version_tag}"
  else
    docker push "${version_tag}"
  fi
}

function docker-cleanup() {
  local -r version_tag="${1}"
  local -r registry_tag="${2}"
  echo "Docker cleanup..."
  if $DRY_RUN; then
    echo "Dry run: would have done "
    echo "docker rmi ${version_tag}"
    echo "docker rmi ${registry_tag}"
  else
    docker rmi "${version_tag}"
    docker rmi "${registry_tag}"
    docker rmi -f $(docker images -q -f "dangling=true") || true
  fi
}

main "$@"
