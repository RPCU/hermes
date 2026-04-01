let
  sources = import ./npins;
in
{
  pkgs,
  lib,
  ...
}:
{
  imports = [ "${sources.nixbook}/devenvModules/devenv.nix" ];

  packages = with pkgs; [
    go
    delve
    air
    git
    jq
    curl
    npins
  ];

  scripts = {
    build.exec = "go build -o hermes .";
    test.exec = "go test -v ./...";
    lint.exec = "golangci-lint run";
    release-snapshot.exec = "goreleaser build --snapshot --clean";
    "run-dry".exec = "go run . --dry-run";
    clean.exec = "rm -rf dist/ hermes";
  };

  enterShell = ''
    echo ""
    echo -e "\\033[1;35m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\\033[0m"
    echo -e "\\033[1;36m   Hermes Development Environment\\033[0m"
    echo -e "\\033[1;35m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\\033[0m"
    echo ""

    echo -e "\\033[1;33m Available tools:\\033[0m"
    echo -e "  \\033[32m✓\\033[0m go                ${pkgs.go.version}"
    echo -e "  \\033[32m✓\\033[0m goreleaser        ${pkgs.goreleaser.version}"
    echo -e "  \\033[32m✓\\033[0m golangci-lint     ${pkgs.golangci-lint.version}"
    echo -e "  \\033[32m✓\\033[0m delve             ${pkgs.delve.version}"
    echo -e "  \\033[32m✓\\033[0m air               ${pkgs.air.version}"
    echo -e "  \\033[32m✓\\033[0m git               ${pkgs.git.version}"
    echo -e "  \\033[32m✓\\033[0m jq                ${pkgs.jq.version}"
    echo -e "  \\033[32m✓\\033[0m curl              ${pkgs.curl.version}"

    echo ""
    echo -e "\\033[1;33m🔧 Quick Commands:\\033[0m"
    echo -e "  \\033[32m•\\033[0m \\033[1mdevenv build\\033[0m                   Build the project"
    echo -e "  \\033[32m•\\033[0m \\033[1mdevenv test\\033[0m                    Run all tests"
    echo -e "  \\033[32m•\\033[0m \\033[1mdevenv run-dry\\033[0m                 Test without API calls"
    echo -e "  \\033[32m•\\033[0m \\033[1mdevenv release-snapshot\\033[0m       Build for all platforms"
    echo -e "  \\033[32m•\\033[0m \\033[1mdevenv lint\\033[0m                   Lint the code"

    echo ""
    echo -e "\\033[1;33m Project Info:\\033[0m"
    if [ -f go.mod ]; then
      MODULE=$(grep '^module' go.mod | awk '{print $2}')
      GO_VERSION=$(grep '^go' go.mod | awk '{print $2}')
      echo -e "  \\033[32m•\\033[0m Module: \\033[1m$MODULE\\033[0m"
      echo -e "  \\033[32m•\\033[0m Go Version: \\033[1m$GO_VERSION\\033[0m"
    fi
    if [ -d .git ]; then
      BRANCH=$(git branch --show-current 2>/dev/null || echo 'unknown')
      echo -e "  \\033[32m•\\033[0m Git Branch: \\033[1m$BRANCH\\033[0m"
    fi

    echo ""
    echo -e "\\033[1;35m━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\\033[0m"
    echo ""
  '';
}
