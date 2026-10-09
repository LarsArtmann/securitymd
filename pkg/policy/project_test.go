package policy

import (
	"path/filepath"
	"testing"

	finding "github.com/larsartmann/go-finding"
	"github.com/onsi/gomega"
)

func TestParseGitRemote(t *testing.T) {
	t.Parallel()

	tests := []struct {
		remote string
		want   RepoIdentity
	}{
		{"https://github.com/AcmeCorp/widget.git", RepoIdentity{"AcmeCorp", "widget"}},
		{"https://github.com/AcmeCorp/widget", RepoIdentity{"AcmeCorp", "widget"}},
		{"git@github.com:AcmeCorp/widget.git", RepoIdentity{"AcmeCorp", "widget"}},
		{"ssh://git@github.com/AcmeCorp/widget.git", RepoIdentity{"AcmeCorp", "widget"}},
		{"ssh://git@github.com:22/AcmeCorp/widget.git", RepoIdentity{"AcmeCorp", "widget"}},
		{"https://gitlab.com/group/subproject/repo.git", RepoIdentity{"subproject", "repo"}},
		{"https://github.com/AcmeCorp/widget/", RepoIdentity{"AcmeCorp", "widget"}},
		{"git@github.com:AcmeCorp/widget.git\r\n", RepoIdentity{"AcmeCorp", "widget"}},
		{"git@github.com:AcmeCorp", RepoIdentity{}},
		{"git://github.com/AcmeCorp/widget.git", RepoIdentity{}},
		{"file:///srv/git/repo", RepoIdentity{}},
		{"/srv/git/repo", RepoIdentity{}},
		{"not-a-remote", RepoIdentity{}},
		{"", RepoIdentity{}},
	}

	for _, test := range tests {
		t.Run(test.remote, func(t *testing.T) {
			t.Parallel()

			gomega.NewWithT(t).Expect(parseGitRemote(test.remote)).To(gomega.Equal(test.want))
		})
	}
}

func TestRepoIdentity_IsComplete(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	g.Expect(RepoIdentity{"a", "b"}.IsComplete()).To(gomega.BeTrue())
	g.Expect(RepoIdentity{"a", ""}.IsComplete()).To(gomega.BeFalse())
	g.Expect(RepoIdentity{}.IsComplete()).To(gomega.BeFalse())
}

func TestDetect_reports_missing_file_as_direct_fix(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	dir := t.TempDir()
	ctx := withWorkingDir(t.Context(), dir)

	findings, err := Detect(ctx)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(findings).To(gomega.HaveLen(1))

	f := findings[0]
	g.Expect(f.Rule).To(gomega.Equal(finding.RuleName("missing-file")))
	g.Expect(f.Severity).To(gomega.Equal(finding.SeverityError))
	g.Expect(f.FixStrategy).To(gomega.Equal(finding.FixStrategySuggest),
		"no git remote in a temp dir: the fix cannot carry content, so it stays suggest-only")
	g.Expect(f.Suggestion).NotTo(gomega.BeEmpty())
}

func TestDetect_validates_existing_policy(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	dir := t.TempDir()
	ctx := withWorkingDir(t.Context(), dir)

	_, err := Generate(ctx, GenerateOptions{
		Organization: "AcmeCorp",
		Repository:   "widget",
	})
	g.Expect(err).NotTo(gomega.HaveOccurred())

	findings, err := Detect(ctx)
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(findings).To(gomega.BeEmpty())
}

func TestDetect_finds_policy_in_github_dir(t *testing.T) {
	t.Parallel()
	g := gomega.NewWithT(t)

	dir := t.TempDir()
	writePolicyFixture(t, filepath.Join(dir, ".github", "SECURITY.md"))

	findings, err := Detect(withWorkingDir(t.Context(), dir))
	g.Expect(err).NotTo(gomega.HaveOccurred())
	g.Expect(findings).To(gomega.BeEmpty(),
		"a compliant .github/SECURITY.md must be discovered and pass")
}
