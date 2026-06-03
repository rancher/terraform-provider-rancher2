{
  description = "A reliable testing environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachSystem [ "aarch64-darwin" "x86_64-linux" "aarch64-linux" ]
      (system:
        let
          pkgs = nixpkgs.legacyPackages.${system};

          leftovers-version = {
            "selected" = "v0.70.0";
          };
          leftovers-prep = {
            "aarch64-darwin" = {
              "url" = "https://github.com/genevieve/leftovers/releases/download/${leftovers-version.selected}/leftovers-${leftovers-version.selected}-darwin-arm64";
              "sha" = "sha256-Tw7G538RYZrwIauN7kI68u6aKS4d/0Efh+dirL/kzoM=";
            };
            "x86_64-linux" = {
              "url" = "https://github.com/genevieve/leftovers/releases/download/${leftovers-version.selected}/leftovers-${leftovers-version.selected}-linux-amd64";
              "sha" = "sha256-D2OPjLlV5xR3f+dVHu0ld6bQajD5Rv9GLCMCk9hXlu8=";
            };
            # linux container running on darwin, actual arm linux isn't in the artifacts
            "aarch64-linux" = {
              "url" = "https://github.com/genevieve/leftovers/releases/download/${leftovers-version.selected}/leftovers-${leftovers-version.selected}-darwin-arm64";
              "sha" = "sha256-Tw7G538RYZrwIauN7kI68u6aKS4d/0Efh+dirL/kzoM=";
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

          terraform-version = {
            "selected" = "1.5.7";
          };
          terraform-prep = {
            "aarch64-darwin" = {
              "url" = "https://releases.hashicorp.com/terraform/${terraform-version.selected}/terraform_${terraform-version.selected}_darwin_arm64.zip";
              "sha" = "sha256-23wz6xpEa3OkQ+LFW1MoRfe3DNVhAL7EyW8Vz6tfUMs=";
              "checksum" = "db7c33eb1a446b73a443e2c55b532845f7b70cd56100bec4c96f15cfab5f50cb";
            };
            "x86_64-linux" = {
              "url" = "https://releases.hashicorp.com/terraform/${terraform-version.selected}/terraform_${terraform-version.selected}_linux_amd64.zip";
              "sha" = "sha256-wO17wy7lKuJVr5mCyMiKekxhBIXPHVX+6wN+q3X6CCw=";
              "checksum" = "c0ed7bc32ee52ae255af9982c8c88a7a4c610485cf1d55feeb037eab75fa082c";
            };
            # linux container running on darwin or arm linux
            "aarch64-linux" = {
              "url" = "https://releases.hashicorp.com/terraform/${terraform-version.selected}/terraform_${terraform-version.selected}_linux_arm64.zip";
              "sha" = "sha256-9LStfGtgiJYKZn40SVyuSQ+wcpR6n/Jmv1kp9TM1ZeQ=";
              "checksum" = "f4b4ad7c6b6088960a667e34495cae490fb072947a9ff266bf5929f5333565e4";
            };
          };
          terraform = pkgs.stdenv.mkDerivation {
            name = "terraform-${terraform-version.selected}";
            src = pkgs.fetchurl {
              url = terraform-prep."${system}".url;
              sha256 = terraform-prep."${system}".sha;
            };
            checksum = terraform-prep."${system}".checksum;
            nativeBuildInputs = [ pkgs.unzip ];
            phases = [ "installPhase" ];
            installPhase = ''
              echo "Verifying checksum..."
              SUM="$(sha256sum $src | awk '{print $1}')"
              if [ "$SUM" = "$checksum" ]; then 
                echo "Valid!";
              else 
                echo "Invalid!";
                echo "expected: $checksum, got: $SUM"
                echo "check your flake.nix file"
                exit 1;
              fi

              install -d $out/bin
              unzip -o $src -d $out/bin
              chmod +x $out/bin/terraform
            '';
          };
          goreleaser-version = {
            "selected" = "2.15.4";
          };
          goreleaser-prep = {
            "aarch64-darwin" = { # local workstation
              "url" = "https://github.com/goreleaser/goreleaser/releases/download/v${goreleaser-version.selected}/goreleaser_Darwin_arm64.tar.gz";
              "sha" = "sha256-2c0JeNZGhutU5T6fCLQA7cix+QNPmPykPsej9kEycIE=";
              "checksum" = "d9cd0978d64686eb54e53e9f08b400edc8b1f9034f98fca43ec7a3f641327081";
            };
            "x86_64-linux" = { # github ubuntu latest runner
              "url" = "https://github.com/goreleaser/goreleaser/releases/download/v${goreleaser-version.selected}/goreleaser_Linux_x86_64.tar.gz";
              "sha" = "sha256-quAMcaSm1V4IzOknOhUWvc4zweB8/7flAvpv7EN33t4=";
              "checksum" = "aae00c71a4a6d55e08cce9273a1516bdce33c1e07cffb7e502fa6fec4377dede";
            };
            "aarch64-linux" = { # linux container running on darwin
              "url" = "https://github.com/goreleaser/goreleaser/releases/download/v${goreleaser-version.selected}/goreleaser_Linux_arm64.tar.gz";
              "sha" = "sha256-3gHKFJdXHps0hBPNLn90vkm41XaWrjhvfu3QYXZUSog=";
              "checksum" = "de01ca1497571e9b348413cd2e7f74be49b8d57696ae386f7eedd06176544a88";
            };
          };
          goreleaser = pkgs.stdenv.mkDerivation {
            name = "goreleaser-${goreleaser-version.selected}";
            src = pkgs.fetchurl {
              url = goreleaser-prep."${system}".url;
              sha256 = goreleaser-prep."${system}".sha;
            };
            checksum = goreleaser-prep."${system}".checksum;
            phases = [ "installPhase" ];
            installPhase = ''
              echo "Verifying checksum..."
              SUM="$(sha256sum $src | awk '{print $1}')"
              if [ "$SUM" = "$checksum" ]; then 
                echo "Valid!";
              else 
                echo "Invalid!";
                echo "expected: $checksum, got: $SUM"
                echo "check your flake.nix file"
                exit 1;
              fi
              install -d $out/bin
              tar -xf $src goreleaser
              install -m 755 goreleaser $out/bin/goreleaser
            '';
          };
          aspellWithDicts = pkgs.aspellWithDicts (d: [d.en d.en-computers]);

          devShellPackage = pkgs.symlinkJoin {
            name = "dev-shell-package";
            paths = [
              # place our downloaded packages here
              aspellWithDicts
              goreleaser
              leftovers
              terraform
            ] ++ (with pkgs; [
              # here are the packages from the nix repository
              actionlint
              awscli2
              bashInteractive
              cmctl
              curl
              eslint
              gh
              git
              gitleaks
              gnupg
              go
              golangci-lint
              gotestfmt
              gotestsum
              jq
              kubectl
              kubernetes-helm
              less
              nodejs_24
              openssh
              openssl
              shellcheck
              tflint
              vim
              which
              yq-go
            ]);
          };
        in
        {
          packages.default = devShellPackage;

          devShells.default = pkgs.mkShell {
            buildInputs = [ devShellPackage ];
            shellHook = ''
              while read word; do echo -e "*$word\n#" | aspell -a --dont-validate-words >/dev/null 2>&1; done < aspell_custom.txt
              export GOROOT="${pkgs.go}/share/go"
              export PATH="${pkgs.go}/share/go/bin:$PATH"
              export PS1="nix:# ";
            '';
          };
        }
      );
}
