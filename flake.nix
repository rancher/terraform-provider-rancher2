{
  description = "A reliable testing environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachSystem [ "x86_64-darwin" "aarch64-darwin" "x86_64-linux" ]
      (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};

          leftovers-version = {
            "selected" = "v0.70.0";
          };
          leftovers-prep = {
            "x86_64-darwin" = {
              "url" = "https://github.com/genevieve/leftovers/releases/download/${leftovers-version.selected}/leftovers-${leftovers-version.selected}-darwin-amd64";
              "sha" = "sha256-HV12kHqB14lGDm1rh9nD1n7Jvw0rCnxmjC9gusw7jfo=";
            };
            "aarch64-darwin" = {
              "url" = "https://github.com/genevieve/leftovers/releases/download/${leftovers-version.selected}/leftovers-${leftovers-version.selected}-darwin-arm64";
              "sha" = "sha256-Tw7G538RYZrwIauN7kI68u6aKS4d/0Efh+dirL/kzoM=";
            };
            "x86_64-linux" = {
              "url" = "https://github.com/genevieve/leftovers/releases/download/${leftovers-version.selected}/leftovers-${leftovers-version.selected}-linux-amd64";
              "sha" = "sha256-D2OPjLlV5xR3f+dVHu0ld6bQajD5Rv9GLCMCk9hXlu8=";
            };
          };
          leftovers = pkgs.stdenv.mkDerivation {
            name = "leftovers-${leftovers-version.selected}";
            src = pkgs.fetchurl {
              url = leftovers-prep."${system}".url;
              sha256 = leftovers-prep."${system}".sha;
            };
            phases = [ "installPhase" ];
            installPhase = ''
              mkdir -p $out/bin
              cp $src $out/bin/leftovers
              chmod +x $out/bin/leftovers
            '';
          };
          aspellWithDicts = pkgs.aspellWithDicts (d: [d.en d.en-computers]);

          devShellPackage = pkgs.symlinkJoin {
            name = "dev-shell-package";
            paths = with pkgs; [
              actionlint
              age
              aspellWithDicts
              awscli2
              bashInteractive
              curl
              dig
              gh
              git
              gitleaks
              gnupg
              go
              gotestfmt
              gotestsum
              jq
              kubectl
              kubernetes-helm
              leftovers
              less
              openssh
              openssl
              shellcheck
              tflint
              tfsec
              tfswitch
              trivy
              updatecli
              vim
              which
              yq
            ];
          };
        in
        {
          packages.default = devShellPackage;

          devShells.default = pkgs.mkShell {
            buildInputs = [ devShellPackage ];
            shellHook = ''
              homebin=$HOME/bin;
              install -d $homebin;
              tfswitch -b $homebin/terraform 1.5.7 &>/dev/null;
              export PATH="$homebin:$PATH";
              export PS1="nix:# ";
            '';
          };
        }
      );
}
