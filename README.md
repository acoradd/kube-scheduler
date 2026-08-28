# kube-scheduler

Out-of-tree `kube-scheduler` binary bundling the upstream default plugins plus
a custom `NodeUptime` Score plugin, so a second scheduler profile can run a
fully configurable, weighted bin-packing pipeline:

1. **InterPodAffinity** (built-in) - respects pod (anti-)affinity constraints.
2. **NodeResourcesFit** (built-in) - `MostAllocated` / `LeastAllocated` / etc.
   via `scoringStrategy`.
3. **NodeUptime** (this repo) - favors older or younger nodes (`Old` /
   `Young`) so pods consolidate away from freshly scaled-up nodes, letting
   the cluster-autoscaler reclaim idle ones sooner.

This repo only builds and publishes the scheduler image to
`ghcr.io/acoradd/kube-scheduler`. Deployment (Helm chart, RBAC, the
`KubeSchedulerConfiguration`) lives in the consuming cluster's own
infrastructure repo, which pulls the published image.

Every plugin in the pipeline is enabled/disabled and weighted through the
standard `KubeSchedulerConfiguration` - see
[`config/scheduler-config.example.yaml`](config/scheduler-config.example.yaml)
for a reference profile to copy into the consuming repo.

## Important: kube-scheduler is weighted-sum, not priority order

kube-scheduler runs every enabled Score plugin and combines them as
`sum(weight_i * normalizedScore_i)`. There is no first-match / short-circuit
priority chain. Tune relative influence between plugins via each plugin's
`weight` in `plugins.score.enabled`.

## Layout

All Go sources live under [`src/`](src) (`go.mod`, `cmd/`, `pkg/`), kept
separate from the repo's CI/CD, Docker, and release tooling at the root.

## Build

```bash
cd src
go mod tidy   # resolves go.sum against the pinned k8s.io/* replace block
go build ./cmd/scheduler
go test ./...
```

## Docker image (multi-arch)

```bash
docker buildx build --platform linux/amd64,linux/arm64 -t kube-scheduler:dev .
```

CI builds and tests on every push/PR (including a non-pushing multi-arch
Docker build check). Tagging a `vX.Y.Z` release publishes a multi-arch
(`linux/amd64`, `linux/arm64`) image to `ghcr.io/acoradd/kube-scheduler`.

## NodeUptime plugin args

```yaml
pluginConfig:
  - name: NodeUptime
    args:
      mode: Old   # or "Young"
```
