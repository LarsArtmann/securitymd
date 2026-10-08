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
