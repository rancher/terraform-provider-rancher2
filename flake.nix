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
            echo "$checksum  $src" | sha256sum -c -
            install -d $out/bin
            unzip -o $src -d $out/bin
            chmod +x $out/bin/terraform
          '';
        };

        gemini-cli = pkgs.gemini-cli-bin.overrideAttrs (oldAttrs: rec {
          version = "0.62.0";
          src = pkgs.fetchzip {
            url = "https://github.com/google-gemini/gemini-cli/releases/download/v${version}/gemini-cli-bundle.zip";
            hash = "sha256-OJx8+Lu/QPszslPzot3d3LgefU+XVR0uI7AWXIsLLhY=";
            stripRoot = false;
          };
        });

        macVscode = pkgs.writeShellScriptBin "code" ''
          exec /usr/local/bin/code "$@"
        '';

        swVers = pkgs.writeShellScriptBin "sw_vers" ''
          exec /usr/bin/sw_vers "$@"
        '';

        initNodeEnv = pkgs.writeShellScriptBin "init-node-env" ''
          echo "Setting up Node environment..."
          if [ ! -f package.json ]; then
            npm init -y
          fi
          if ls eslint.config.* 1> /dev/null 2>&1; then
            echo "ESLint config already exists, skipping."
          else
            echo "Generating default eslint.config.mjs..."
            cat <<'EOF' > eslint.config.mjs
import eslintPluginJs from "@eslint/js";

export default [
  {
    ...eslintPluginJs.configs.recommended,
    rules: {
      semi: ["error", "always"],
      indent: ["error", 2],
      "no-unused-vars": "warn",
    },
  },
];
EOF
          fi
        '';

        unfreePkgs = import nixpkgs {
          inherit system;
          config = {
            allowUnfree = true;
          };
        };

        claude-code = unfreePkgs.claude-code;
        github-copilot-cli = unfreePkgs.github-copilot-cli;

        devPackages = [
          # downloaded packages here
          leftovers
          terraform
          initNodeEnv
          gemini-cli
        ] ++ ([
          # unfree packages from the nix repository
          claude-code
          github-copilot-cli
        ]) ++ (with pkgs; [
          # free packages from the nix repository
          actionlint
          age
          awscli2
          bashInteractive
          cspell
          curl
          dig
          docker-client
          docker-compose
          eslint
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
          markdownlint-cli
          nodejs_24
          openssh
          openssl
          prettier
          ripgrep
          shellcheck
          shfmt
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
        ]) ++ ( if pkgs.stdenv.hostPlatform.isDarwin then [
          # mac only packages
          macVscode
          swVers
          pkgs.pinentry_mac
          pkgs.colima
          pkgs.age-plugin-se
        ] else []);

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
            '';
          };
        }
      );
}
