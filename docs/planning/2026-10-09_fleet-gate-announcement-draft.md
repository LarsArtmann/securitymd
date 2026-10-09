# Fleet gate announcement — DRAFT (not sent)

_Status: prepared 2026-10-09, execution gated on Lars (TODO_LIST "Fleet gate
announcement"). Fill the two `TODO` placeholders before sending._

---

## What is coming

`securitymd` is becoming a default BuildFlow provider: every manifest-carrying
repository gets its `SECURITY.md` validated on each run. A repository without
a policy — or with a policy missing the sections GitHub expects — fails the
findings gate with clear, rule-tagged messages (`missing-file`,
`missing-contact`, …).

## What each repository needs

Nothing, in the common case. Running

```bash
buildflow -s securitymd --fix
```

generates a compliant `SECURITY.md` from the embedded template when none
exists (it never overwrites a policy you wrote by hand).

## If a finding is a false positive

Two per-repo escape hatches, both visible in reports instead of silent:

- **Suppress** the rule inside the policy, with a reason (reason-less is inert):

  ```markdown
  <!-- securitymd:ignore(no-version-info) bootstrap stage, filled next sprint -->
  ```

- **Downgrade** the rule for the repo in `.buildflow.yml`:

  ```yaml
  tool_options:
    securitymd:
      severity-overrides: missing-file=warning
  ```

`missing-file` cannot be suppressed (there is no file to carry the comment);
downgrade it or generate a policy.

## The sweep

Before the gate lands in CI defaults, one sweep runs `buildflow -s securitymd
--fix` across the fleet and reports which repositories were repaired. The
script is ready: `scripts/fleet-securitymd-sweep.sh` in the securitymd
repository.

- Sweep date: TODO (Lars)
- Announcement channel: TODO (Lars)

## Rollout order

1. Announce (this document, filled in).
2. Run the sweep; share the repair summary.
3. Flip the provider default in BuildFlow; repositories that still fail have
   the escape hatches above.
