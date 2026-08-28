package nodeuptime

import (
	"context"
	"testing"

	"k8s.io/kubernetes/pkg/scheduler/framework"
)

func TestNormalizeScoreOldFavorsOlderNodes(t *testing.T) {
	pl := &NodeUptime{mode: ModeOld}
	scores := framework.NodeScoreList{
		{Name: "old-node", Score: 1000},   // created earliest
		{Name: "young-node", Score: 2000}, // created latest
	}

	if status := pl.NormalizeScore(context.Background(), nil, nil, scores); !status.IsSuccess() {
		t.Fatalf("unexpected status: %v", status)
	}

	old := scoreFor(t, scores, "old-node")
	young := scoreFor(t, scores, "young-node")
	if old <= young {
		t.Fatalf("expected old-node to score higher than young-node, got old=%d young=%d", old, young)
	}
	if old != framework.MaxNodeScore {
		t.Fatalf("expected old-node to reach MaxNodeScore, got %d", old)
	}
	if young != 0 {
		t.Fatalf("expected young-node to score 0, got %d", young)
	}
}

func TestNormalizeScoreYoungFavorsYoungerNodes(t *testing.T) {
	pl := &NodeUptime{mode: ModeYoung}
	scores := framework.NodeScoreList{
		{Name: "old-node", Score: 1000},
		{Name: "young-node", Score: 2000},
	}

	if status := pl.NormalizeScore(context.Background(), nil, nil, scores); !status.IsSuccess() {
		t.Fatalf("unexpected status: %v", status)
	}

	old := scoreFor(t, scores, "old-node")
	young := scoreFor(t, scores, "young-node")
	if young <= old {
		t.Fatalf("expected young-node to score higher than old-node, got old=%d young=%d", old, young)
	}
}

func TestNormalizeScoreAllNodesEqualAge(t *testing.T) {
	pl := &NodeUptime{mode: ModeOld}
	scores := framework.NodeScoreList{
		{Name: "a", Score: 500},
		{Name: "b", Score: 500},
	}

	if status := pl.NormalizeScore(context.Background(), nil, nil, scores); !status.IsSuccess() {
		t.Fatalf("unexpected status: %v", status)
	}

	for _, s := range scores {
		if s.Score != framework.MaxNodeScore {
			t.Fatalf("expected %s to score MaxNodeScore when all nodes tie, got %d", s.Name, s.Score)
		}
	}
}

func scoreFor(t *testing.T, scores framework.NodeScoreList, name string) int64 {
	t.Helper()
	for _, s := range scores {
		if s.Name == name {
			return s.Score
		}
	}
	t.Fatalf("node %q not found in scores", name)
	return 0
}
