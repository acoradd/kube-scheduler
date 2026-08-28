// Command scheduler is an out-of-tree kube-scheduler binary that bundles
// the default upstream plugins plus NodeUptime, so a single weighted
// scoring pipeline (InterPodAffinity, NodeResourcesFit, NodeUptime, ...)
// can be configured via a KubeSchedulerConfiguration.
package main

import (
	"os"

	"k8s.io/component-base/cli"
	"k8s.io/kubernetes/cmd/kube-scheduler/app"

	"acoradd/kube-scheduler/pkg/plugins/nodeuptime"
)

func main() {
	command := app.NewSchedulerCommand(
		app.WithPlugin(nodeuptime.Name, nodeuptime.New),
	)

	code := cli.Run(command)
	os.Exit(code)
}
