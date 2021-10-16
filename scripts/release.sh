#!/usr/bin/env bash

prompt_for_tag() {
  read -p "GitHub tag must be created prior to release.
Enter tag (v0.0.0): v" input
  echo -e ${input}
}

main() {
  tag=$(prompt_for_tag)

  if [[ ! -z "${tag}" ]]; then
    echo "Creating GitHub tag: ${tag}"
    git tag ${tag}
    git push origin tag ${tag}
  else
    echo "Tag cannot be empty, aborting."
    exit 1
  fi

  if [[ ! -z "${GITHUB_TOKEN}" ]]; then
    echo "Releasing binaries... (tag: ${tag})"
    goreleaser release --rm-dist
  else
    echo "ERROR: GITHUB_TOKEN is not set"
  fi
}

main "$@"