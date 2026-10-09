package policy

import (
	"strings"
	"testing"

	"github.com/onsi/gomega"
)

// FuzzParseGitRemote proves the parser's core invariant for arbitrary input:
// a complete identity consists of clean coordinates (no separators, no
// whitespace, non-empty), so a rendered advisory link can never be garbage.
func FuzzParseGitRemote(f *testing.F) {
	for _, seed := range []string{
		"https://github.com/AcmeCorp/widget.git",
		"https://github.com/AcmeCorp/widget",
		"git@github.com:AcmeCorp/widget.git",
		"ssh://git@github.com/AcmeCorp/widget.git",
		"ssh://git@github.com:22/AcmeCorp/widget.git",
		"https://gitlab.com/group/subproject/repo.git",
		"https://github.com/AcmeCorp/widget/",
		"git@github.com:AcmeCorp/widget.git\r\n",
		"git@github.com:AcmeCorp",
		"git://github.com/AcmeCorp/widget.git",
		"file:///srv/git/repo",
		"/srv/git/repo",
		"https://github.com//widget.git",
		"not-a-remote",
		"",
	} {
		f.Add(seed)
	}

	f.Fuzz(func(t *testing.T, remote string) {
		identity := parseGitRemote(remote)
		if !identity.IsComplete() {
			return
		}

		g := gomega.NewWithT(t)

		clean := func(part string) bool {
			return part != "" && !strings.ContainsAny(part, ":/ \t\r\n")
		}

		g.Expect(clean(identity.Organization)).
			To(gomega.BeTrue(), "organization %q contains separators or whitespace", identity.Organization)
		g.Expect(clean(identity.Repository)).
			To(gomega.BeTrue(), "repository %q contains separators or whitespace", identity.Repository)
	})
}
