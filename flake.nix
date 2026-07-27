{
  description = "A reliable testing environment";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachSystem [ "aarch64-darwin" "x86_64-linux" "aarch64-linux" ]
      (system: let
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
          # linux container running on darwin, actual arm linux isnt in the artifacts
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
            echo "$checksum  $src" | sha256sum -c -
            install -d $out/bin
            unzip -o $src -d $out/bin
            chmod +x $out/bin/terraform
          '';
        };

        macVscode = pkgs.writeShellScriptBin "code" ''
          exec /usr/local/bin/code "$@"
        '';

        swVers = pkgs.writeShellScriptBin "sw_vers" ''
          exec /usr/bin/sw_vers "$@"
        '';

        claude-code = (import nixpkgs {
          inherit system;
          config = {
            allowUnfree = true;
          };
        }).claude-code;

        devPackages = [
          # place our downloaded packages here
          leftovers
          terraform
          macVscode
          swVers
          claude-code
        ] ++ (with pkgs; [
          # here are the packages from the nix repository
          actionlint
          age
          awscli2
          bashInteractive
          colima
          cspell
          curl
          dig
          docker-client
          docker-compose
          eslint
          gemini-cli
          gh
          git
          gitleaks
          gnupg
          go
          golangci-lint
          google-cloud-sdk
          goreleaser
          gotestfmt
          gotestsum
          jq
          kubectl
          kubernetes-helm
          less
          nodejs_26
          openssh
          openssl
          ripgrep
          shellcheck
          tflint
          tfsec
          time
          tree
          trivy
          updatecli
          vim
          which
          xz
          yq-go
        ]) ++ (if pkgs.stdenv.isDarwin then [ pkgs.colima ] else []);

        devShellPackage = pkgs.symlinkJoin {
          name = "dev-shell-package";

          # buildEnv properly handles combining all outputs (like bin, out, etc.)
          paths = [
            (pkgs.buildEnv {
              name = "dev-shell-env";
              paths = devPackages;
            })
          ];
        };
        in
        {
          packages.default = devShellPackage;

          devShells.default = pkgs.mkShell {
            buildInputs = [ devShellPackage ];
            shellHook = ''
              export PS1="nix:# ";
              install -d ~/.docker/cli-plugins/ || true;
              ln -sfn $(which docker-compose) ~/.docker/cli-plugins/docker-compose || true;
            '';
          };
        }
      );
}
