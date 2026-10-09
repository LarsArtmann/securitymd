// Package policy validates and generates SECURITY.md files for repositories.
//
// It is the detector/generator core consumed by both the securitymd CLI and
// the BuildFlow toolsdk provider. All problems are reported as go-finding
// findings with stable rule IDs; generation writes the embedded canonical
// template atomically and never overwrites an existing policy.
package policy

import (
	"fmt"

	"github.com/larsartmann/go-finding"
)

// ToolName is the stable namespace for findings, BuildFlow tool selection
// (`buildflow -s securitymd`), and suppressions.
const ToolName finding.ToolName = "securitymd"

// CandidateLocations lists the SECURITY.md locations GitHub recognizes, in
// priority order. The first entry is the canonical write target when none
// exists.
var CandidateLocations = []string{ //nolint:gochecknoglobals // fixed discovery order, read-only after init
	"SECURITY.md",
	".github/SECURITY.md",
	"docs/SECURITY.md",
}

// OrderedCandidates returns the candidate locations with the preferred
// canonical location first — for repositories that treat .github/ or docs/
// as their policy home. Unknown preferences are rejected; the default order
// is returned unchanged for "" and LocationRoot.
func OrderedCandidates(preference string) ([]string, error) {
	switch preference {
	case "", LocationRoot:
		return CandidateLocations, nil
	case LocationGitHub, LocationDocs:
	default:
		return nil, fmt.Errorf("unknown policy location %q (want %s, %s, or %s)",
			preference, LocationRoot, LocationGitHub, LocationDocs)
	}

	preferred := preference + "/" + CandidateLocations[0]
	order := make([]string, 0, len(CandidateLocations))
	order = append(order, preferred)

	for _, candidate := range CandidateLocations {
		if candidate != preferred {
			order = append(order, candidate)
		}
	}

	return order, nil
}
