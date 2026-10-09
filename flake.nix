{
  description = "securitymd — SECURITY.md policy linter and generator";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
    flake-parts = {
      url = "github:hercules-ci/flake-parts";
      inputs.nixpkgs-lib.follows = "nixpkgs";
    };
    treefmt-nix = {
      url = "github:numtide/treefmt-nix";
      inputs.nixpkgs.follows = "nixpkgs";
    };
  };

  outputs =
    inputs@{
      self,
      flake-parts,
      treefmt-nix,
      ...
    }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      # Inline systems: nixpkgs 26.11 dropped x86_64-darwin, which github:nix-systems/default still lists.
      systems = [
        "x86_64-linux"
        "aarch64-linux"
        "aarch64-darwin"
      ];

      imports = [
        treefmt-nix.flakeModule
      ];

      perSystem =
        {
          config,
          pkgs,
          ...
        }:
        let
          docsGate = pkgs.writeShellScriptBin "docs-gate" ''
            set -euo pipefail

            repo="''${1:-$PWD}"
            status=0

            # Gate 1: every archived status/planning doc must carry inline
            # verdict annotations (the docs-health "~~" convention) — an
            # unannotated file means the annotation sweep missed it.
            for archived in "$repo"/docs/archive/pre-rebuild/*.md; do
              [ -e "$archived" ] || continue
              if ! grep -q '~~' "$archived"; then
                echo "docs-gate: unannotated archive file: $archived" >&2
                status=1
              fi
            done

            # Gate 2: backticked docs/ references in living docs must resolve
            # to an existing file or directory. Excluded: docs/SECURITY.md (a
            # policy candidate-location string, not a repo file) and
            # docs/reviews/ (convention decision still open, tracked in
            # TODO_LIST).
            for doc in "$repo"/README.md "$repo"/AGENTS.md "$repo"/TODO_LIST.md \
              "$repo"/ROADMAP.md "$repo"/FEATURES.md "$repo"/CHANGELOG.md; do
              [ -f "$doc" ] || continue
              # "|| true": grep exiting 1 on zero matches is not a failure.
              refs="$(grep -ohE '`docs/[A-Za-z0-9_./-]+`' "$doc" | tr -d '`' | sort -u || true)"
              for ref in $refs; do
                case "$ref" in
                  docs/SECURITY.md | docs/reviews/) continue ;;
                esac
                if [ ! -e "$repo/$ref" ]; then
                  echo "docs-gate: dangling reference in $(basename "$doc"): $ref" >&2
                  status=1
                fi
              done
            done

            exit "$status"
          '';
        in
        {
          treefmt = {
            projectRootFile = "go.mod";
            # goimports is deliberately absent: nixpkgs gotools bundles Go
            # 1.26 while go.mod requires 1.27, so the sandboxed check tries to
            # download a newer toolchain and fails (no network). Import
            # correctness is enforced by golangci-lint outside the sandbox.
            programs = {
              gofumpt.enable = true;
              nixfmt.enable = true;
            };
          };

          checks.format = config.treefmt.build.check self;

          checks.docs-gate =
            pkgs.runCommand "docs-gate"
              {
                src = self;
                nativeBuildInputs = [ pkgs.bash ];
              }
              ''
                ${pkgs.lib.getExe docsGate} "$src"
                touch $out
              '';

          apps.docs-gate = {
            type = "app";
            program = pkgs.lib.getExe docsGate;
          };

          devShells = {
            default = pkgs.mkShell {
              name = "securitymd-dev";

              packages = [
                pkgs.go_1_27
                pkgs.golangci-lint
                pkgs.gopls
                pkgs.delve
                pkgs.gotools
                pkgs.gofumpt
              ];

              GOWORK = "off";
              GOTOOLCHAIN = "local";
              GOPRIVATE = "github.com/LarsArtmann/*,github.com/larsartmann/*";
            };

            ci = pkgs.mkShellNoCC {
              packages = [
                pkgs.go_1_27
                pkgs.golangci-lint
              ];

              GOWORK = "off";
              GOTOOLCHAIN = "local";
              GOPRIVATE = "github.com/LarsArtmann/*,github.com/larsartmann/*";
            };
          };
        };
    };
}
