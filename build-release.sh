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

  check-prereqs

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

  read -p "Container name (default: recall): " container_name
  [ -z "$container_name" ] && container_name="recall"

  read -p "S3 bucket (default: recall.builds): " s3_bucket
  [ -z "$s3_bucket" ] && s3_bucket="recall.builds"

  local -r github=`git config --get remote.origin.url`
  declare -r temp_dir=$(mktemp -d "/tmp/${container_name}-${new_version}.XXXX")
  local -r tag="${container_name}:${new_version}"
  local -r tarball="/tmp/${tag}.tar.gz"
  local -r s3_path="s3://${s3_bucket}/${container_name}/${new_version}.tar.gz"

  git-clone "${github}" "${temp_dir}"
  git-checkout "${new_version}" "${temp_dir}"
  docker-build "${tag}" "${temp_dir}"
  docker-save "${tag}" "${tarball}"
  s3-copy "${tarball}" "${s3_path}"
  docker-cleanup "${tag}"
  delete-tarball "${tarball}"
  rm -Rf "${temp_dir}"
}

function usage() {
  echo "Usage: ${0} <release_version> [--no-dry-run]"
  echo
  echo "<release_version> is the version you want to release,"
  echo "Please see docs/releasing.md for more info."
}

function check-prereqs() {
  if ! (aws --version 2>&1 | grep -q aws-cli); then
    echo "!!! AWS SDK is required. Use 'go get github.com/aws/aws-sdk-go'."
    exit 1
  fi
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

function docker-save() {
  local -r image="${1}"
  local -r dest="${2}"
  echo "Saving '${image}' to '${dest}'..."
  if $DRY_RUN; then
    echo "Dry run: would have done docker save ${image} | gzip > ${dest}"
  else
    docker save "${image}" | gzip > "${dest}"
  fi
}

function docker-cleanup() {
  local -r tag="${1}"
  echo "Docker cleanup..."
  if $DRY_RUN; then
    echo "Dry run: would have done docker rmi ${tag}"
  else
    docker rmi "${tag}"
    docker rmi $(docker images -q -f "dangling=true") || true
  fi
}

function delete-tarball() {
  local -r tarball="${1}"
  echo "Deleting tarball '${tarball}'..."
  if $DRY_RUN; then
    echo "Dry run: would have done rm ${tarball}"
  else
    rm  "${tarball}"
  fi
}

function s3-copy() {
  local -r src="${1}"
  local -r dest="${2}"
  echo "Pushing '${src}' to '${dest}'..."
  if $DRY_RUN; then
    echo "Dry run: would have done aws s3 cp ${src} ${dest} --acl=private --content-encoding=gzip"
  else
    aws s3 cp "${src}" "${dest}" --acl=private --content-encoding=gzip
  fi
}

main "$@"
