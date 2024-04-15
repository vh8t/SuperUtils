#!/bin/bash

if ! command -v git &> /dev/null; then
  echo "Error: Git is not installed. Please install Git before proceeding."
  exit 1
fi

if ! command -v go &> /dev/null; then
  echo "Error: Go (golang) is not installed. Please install Go before proceeding."
  exit 1
fi

INSTALL_DIR="$HOME/superutils"
BIN_DIR="$INSTALL_DIR/bin"
REPO_URL="https://github.com/vh8t/SuperUtils.git"
TOOLS=("sls.go" "ccat.go")

if [ -z "$(ls -A $INSTALL_DIR)" ]; then
  echo "Installation directory '$INSTALL_DIR' is empty."
  echo "Cloning repository into '$INSTALL_DIR'..."

  git clone $REPO_URL $INSTALL_DIR

  mkdir -p $BIN_DIR
else
  echo "Installation cancelled. To reinstall, remove the ~/superutils/ directory first"
  exit 1
fi

echo "Initializing project"
go mod init superutils

echo "Installing dependencies (chroma, pflag)"
go get github.com/alecthomas/chroma/v2
go get github.com/spf13/pflag

echo "Compiling tools..."
for tool in "${TOOLS[@]}"; do
  tool_name=$(basename -s .go "$tool")
  echo "Compiling $tool_name..."
  go build -o "$BIN_DIR/$tool_name" "$INSTALL_DIR/$tool"
done

update_path_in_shell_profile() {
  local shell_profile="$1"
  echo "Updating PATH in $shell_profile..."
  echo "export PATH=\"$BIN_DIR:\$PATH\"" >> "$shell_profile"
  source "$shell_profile"
}

shell_profiles=()
if [ -f "$HOME/.bashrc" ]; then
  available_profiles+=("$HOME/.bashrc")
fi
if [ -f "$HOME/.bash_profile" ]; then
  available_profiles+=("$HOME/.bash_profile")
fi
if [ -f "$HOME/.zshrc" ]; then
  available_profiles+=("$HOME/.zshrc")
fi

if [ ${#available_profiles[@]} -eq 0 ]; then
  echo "Unable to locate shell profile (.bashrc, .bash_profile, or .zshrc)."
  echo "Please manually add the following line to your shell configuration file:"
  echo "export PATH=\"$BIN_DIR:\$PATH\""
  exit 1
fi

for profile in "${shell_profiles[@]}"; do
  if [ -f "$profile" ]; then
    update_path_in_shell_profile "$profile"
  else
    echo "Warning: Shell profile '$profile' not found. Skipping."
  fi
done

echo "Installation complete. Tools are now accessible from the command line."
