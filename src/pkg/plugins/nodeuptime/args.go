package nodeuptime

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// Mode selects which end of the node age spectrum scores highest.
type Mode string

const (
	// ModeOld favors older (longer-running) nodes, consolidating pods
	// away from freshly scaled-up nodes so they can be scaled back down.
	ModeOld Mode = "Old"
	// ModeYoung favors younger (recently-created) nodes.
	ModeYoung Mode = "Young"
)

// Args configures the NodeUptime plugin.
type Args struct {
	metav1.TypeMeta `json:",inline"`

	// Mode is either "Old" or "Young". Defaults to "Old".
	Mode Mode `json:"mode,omitempty"`
}

// DeepCopyObject implements runtime.Object.
func (in *Args) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(Args)
	*out = *in
	return out
}
