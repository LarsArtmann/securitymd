package main

import (
	"bytes"
	"strings"
	"testing"

	. "github.com/onsi/gomega"
)

func TestPromptIdentity_reads_both_coordinates(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	var out bytes.Buffer

	org, repo, err := promptIdentity(strings.NewReader("AcmeCorp\nwidget\n"), &out)
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(org).To(Equal("AcmeCorp"))
	g.Expect(repo).To(Equal("widget"))
	g.Expect(out.String()).To(ContainSubstring("Organization"))
	g.Expect(out.String()).To(ContainSubstring("Repository"))
}

func TestPromptIdentity_trims_whitespace(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	org, repo, err := promptIdentity(strings.NewReader("  AcmeCorp  \n\twidget\n"), &bytes.Buffer{})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(org).To(Equal("AcmeCorp"))
	g.Expect(repo).To(Equal("widget"))
}

func TestPromptIdentity_accepts_final_line_without_newline(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	org, repo, err := promptIdentity(strings.NewReader("AcmeCorp\nwidget"), &bytes.Buffer{})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(org).To(Equal("AcmeCorp"))
	g.Expect(repo).To(Equal("widget"))
}

func TestPromptIdentity_rejects_empty_answers_with_flag_guidance(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	_, _, err := promptIdentity(strings.NewReader("\nwidget\n"), &bytes.Buffer{})
	g.Expect(err).To(MatchError(errEmptyPromptAnswer))
	g.Expect(err.Error()).To(ContainSubstring("--organization"),
		"the error must point non-interactive users at the flags")

	_, _, err = promptIdentity(strings.NewReader("AcmeCorp\n   \n"), &bytes.Buffer{})
	g.Expect(err).To(MatchError(errEmptyPromptAnswer))
	g.Expect(err.Error()).To(ContainSubstring("--repository"))
}

func TestPromptIdentity_rejects_immediate_eof(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	_, _, err := promptIdentity(strings.NewReader(""), &bytes.Buffer{})
	g.Expect(err).To(HaveOccurred())
}

// In tests and CI, stdin is never a char device: identity resolution must
// stay silent and hand back the flags untouched — the honest skip (never a
// blocking prompt) is what keeps pipelines moving.
func TestResolveIdentity_is_inert_without_tty(t *testing.T) {
	t.Parallel()
	g := NewWithT(t)

	org, repo, err := resolveIdentity(t.Context(), setupFlags{})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(org).To(BeEmpty(), "no prompt may run without an interactive terminal")
	g.Expect(repo).To(BeEmpty())

	org, repo, err = resolveIdentity(t.Context(), setupFlags{Organization: "flag-org", Repository: "flag-repo"})
	g.Expect(err).NotTo(HaveOccurred())
	g.Expect(org).To(Equal("flag-org"))
	g.Expect(repo).To(Equal("flag-repo"))
}
