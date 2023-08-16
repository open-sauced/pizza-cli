#!/bin/sh

# This is a convenience script that can be downloaded from GitHub and
# piped into "sh" for conveniently downloading the latest GitHub release
# of the pizza CLI:
#
# curl -fsSL https://raw.githubusercontent.com/open-sauced/pizza-cli/main/install.sh | sh
#
# Warning: It may not be advisable to pipe scripts from GitHub directly into
# a command line interpreter! If you do not fully trust the source, first
# download the script, inspect it manually to ensure its integrity, and then
# run it:
#
# curl -fsSL https://raw.githubusercontent.com/open-sauced/pizza-cli/main/install.sh > install.sh
# vim install.sh
# ./install.sh

PIZZA_REPO="open-sauced/pizza-cli"
ARCH=""

# Detect architecture
case "$(uname -m)" in
    x86_64) ARCH="x86_64" ;;
    arm64)  ARCH="arm64"  ;;
    *)      echo "Unsupported architecture"; exit 1 ;;
esac

# Detect OS system type. Windows not supported.
case "$(uname -s)" in
    Darwin) OSTYPE="darwin" ;;
    *)      OSTYPE="linux"  ;;
esac

# Fetch download URL for the architecture from the GitHub API
ASSET_URL=$(curl -s https://api.github.com/repos/$PIZZA_REPO/releases/latest | \
           grep -o "https:\/\/github\.com\/open-sauced\/pizza-cli\/releases\/download\/.*${OSTYPE}-${ARCH}.*")

if [ -z "$ASSET_URL" ]; then
    echo "Could not find a binary for latest version of Pizza CLI release and architecture ${ARCH} on OS type ${OSTYPE}"
    exit 1
fi

# Download and install
curl -L "${ASSET_URL}" -o ./pizza
chmod +x ./pizza

echo
echo "Download complete. Stay saucy 🍕"
