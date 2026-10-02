#!/usr/bin/env bash
# Retargets src/go.mod onto a given Kubernetes release so the same code base
# can be built against several minors (see the k8s matrix in CI).
#
# Usage: hack/set-k8s-version.sh 1.35.9
set -euo pipefail

version="${1:?usage: $0 <kubernetes version, e.g. 1.35.9>}"
version="${version#v}"
staging="v0.${version#1.}" # k8s.io/kubernetes v1.X.Y <=> staging modules v0.X.Y

cd "$(dirname "$0")/../src"

# Indirect requirements are dropped so `go mod tidy` recomputes them from the
# target release's dependency graph instead of keeping newer versions around.
sed -E -i '/\/\/ indirect\r?$/d' go.mod

edits=(-require="k8s.io/kubernetes@v${version}")
for mod in $(awk '/=> k8s\.io\// && $1 != "k8s.io/legacy-cloud-providers" { print $1 }' go.mod); do
	edits+=(-replace="${mod}=${mod}@${staging}")
done
for mod in $(awk '/^[[:space:]]+k8s\.io\/[^ ]+ v0\./ && !/=>/ { print $1 }' go.mod); do
	edits+=(-require="${mod}@${staging}")
done
go mod edit "${edits[@]}"
go mod tidy
