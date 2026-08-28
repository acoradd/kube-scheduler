// Package nodeuptime implements a Score plugin that ranks nodes by age,
// so pods can be consolidated onto (or away from) long-running nodes as
// part of a weighted bin-packing pipeline.
package nodeuptime

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/kubernetes/pkg/scheduler/framework"
	frameworkruntime "k8s.io/kubernetes/pkg/scheduler/framework/runtime"
)

// Name is the name registered in the scheduler configuration's plugin lists.
const Name = "NodeUptime"

// NodeUptime scores nodes by creation timestamp.
type NodeUptime struct {
	handle framework.Handle
	mode   Mode
}

var _ framework.ScorePlugin = &NodeUptime{}

func (pl *NodeUptime) Name() string { return Name }

// New builds a NodeUptime plugin instance from the raw plugin config.
func New(_ context.Context, obj runtime.Object, h framework.Handle) (framework.Plugin, error) {
	args := &Args{Mode: ModeOld}
	if obj != nil {
		if err := frameworkruntime.DecodeInto(obj, args); err != nil {
			return nil, fmt.Errorf("decoding NodeUptime args: %w", err)
		}
	}
	if args.Mode == "" {
		args.Mode = ModeYoung
	}
	if args.Mode != ModeOld && args.Mode != ModeYoung {
		return nil, fmt.Errorf("invalid NodeUptime mode %q: must be %q or %q", args.Mode, ModeOld, ModeYoung)
	}
	return &NodeUptime{handle: h, mode: args.Mode}, nil
}

// Score returns the node's creation time as a unix timestamp; normalized
// in NormalizeScore relative to the other candidate nodes.
func (pl *NodeUptime) Score(_ context.Context, _ *framework.CycleState, _ *v1.Pod, nodeName string) (int64, *framework.Status) {
	nodeInfo, err := pl.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		return 0, framework.AsStatus(fmt.Errorf("getting node %q from snapshot: %w", nodeName, err))
	}
	return nodeInfo.Node().CreationTimestamp.Unix(), nil
}

func (pl *NodeUptime) ScoreExtensions() framework.ScoreExtensions { return pl }

// NormalizeScore rescales raw unix timestamps to the [0, MaxNodeScore] range
// expected by the framework, inverting the ordering when mode is ModeOld
// since a smaller timestamp means an older node.
func (pl *NodeUptime) NormalizeScore(_ context.Context, _ *framework.CycleState, _ *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	if len(scores) == 0 {
		return nil
	}

	min, max := scores[0].Score, scores[0].Score
	for _, s := range scores {
		if s.Score < min {
			min = s.Score
		}
		if s.Score > max {
			max = s.Score
		}
	}

	span := max - min
	for i, s := range scores {
		if span == 0 {
			scores[i].Score = framework.MaxNodeScore
			continue
		}
		if pl.mode == ModeOld {
			scores[i].Score = framework.MaxNodeScore * (max - s.Score) / span
		} else {
			scores[i].Score = framework.MaxNodeScore * (s.Score - min) / span
		}
	}
	return nil
}
