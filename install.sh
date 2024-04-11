#!/bin/bash

REPO_URL="https://github.com/vh8t/SuperUtils.git"
TOOLS=("ccat.py" "sls.py" "sutil.py")

INSTALL_DIR="$HOME/superutils"

git clone "$REPO_URL" "$INSTALL_DIR"

cd "$INSTALL_DIR"

for tool in "${TOOLS[@]}"; do
  tool_name="$(basename "$tool" .py)"

  chmod +x "$tool"
  sudo ln -sf "$INSTALL_DIR/$tool" "/usr/local/bin/$tool_name"

  MAN_PAGE="$tool_name.1"
  if [ -f "$MAN_PAGE" ]; then
    groff -man -Tascii "$MAN_PAGE" | gzip > "$MAN_PAGE".gz
    sudo cp "$MAN_PAGE".gz /usr/share/man/man1/
  fi
done

sudo mandb

echo "SuperUtils installed successfully! Run 'sutil --setup' to get started"
