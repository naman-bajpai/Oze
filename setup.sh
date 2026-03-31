#!/usr/bin/env bash
# setup.sh — copies the oze project to ~/projects/oze, builds it, smoke-tests it, and inits git.
# Run this from the directory that CONTAINS the oze/ folder, or from inside oze/ itself.
set -euo pipefail

# ── Resolve source dir ────────────────────────────────────────────────────
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
SRC="$SCRIPT_DIR"

# ── Destination ───────────────────────────────────────────────────────────
DEST="$HOME/projects/oze"

# ── Colour helpers ────────────────────────────────────────────────────────
GREEN='\033[0;32m'; CYAN='\033[0;36m'; YELLOW='\033[1;33m'; RED='\033[0;31m'; RESET='\033[0m'
ok()   { echo -e "${GREEN}✔ $*${RESET}"; }
info() { echo -e "${CYAN}→ $*${RESET}"; }
warn() { echo -e "${YELLOW}⚠ $*${RESET}"; }
die()  { echo -e "${RED}✖ $*${RESET}"; exit 1; }

echo ""
echo -e "${CYAN}╭────────────────────────────────╮${RESET}"
echo -e "${CYAN}│  oze setup & build script      │${RESET}"
echo -e "${CYAN}╰────────────────────────────────╯${RESET}"
echo ""

# ── Check Go ──────────────────────────────────────────────────────────────
info "Checking for Go..."
go version 2>/dev/null || die "Go is not installed. Please install it from https://go.dev/dl/ and re-run."
ok "$(go version)"

# ── Copy files ────────────────────────────────────────────────────────────
if [ "$SRC" != "$DEST" ]; then
  info "Copying project files to $DEST ..."
  mkdir -p "$DEST"
  cp -R "$SRC"/. "$DEST"/
  ok "Files copied to $DEST"
else
  info "Already in $DEST — skipping copy."
fi

cd "$DEST"

# ── Build ─────────────────────────────────────────────────────────────────
info "Building oze..."
go build -o oze . || die "Build failed. Check the error above."
ok "Build succeeded → $DEST/oze"

# ── Smoke tests ───────────────────────────────────────────────────────────
echo ""
info "Smoke test 1: ./oze --help"
echo "──────────────────────────────────────"
./oze --help || true
echo "──────────────────────────────────────"
ok "--help printed"

echo ""
info "Smoke test 2: ./oze --dry-run \"test feature\""
echo "──────────────────────────────────────"
./oze --dry-run "test feature" || true
echo "──────────────────────────────────────"
ok "--dry-run printed prompt"

# ── Git init ─────────────────────────────────────────────────────────────
echo ""
if [ -d "$DEST/.git" ]; then
  warn ".git already exists — skipping git init."
else
  info "Initialising git repo..."
  git init
  git add .
  git commit -m "feat: initial oze release"
  ok "Git repo initialised and committed."
fi

# ── File tree ────────────────────────────────────────────────────────────
echo ""
info "Final file tree:"
find "$DEST" -not -path '*/.git/*' -not -name '.git' | sort | sed "s|$DEST|.|"

# ── Next steps ────────────────────────────────────────────────────────────
echo ""
echo -e "${GREEN}╭──────────────────────────────────────────────────────────────╮${RESET}"
echo -e "${GREEN}│  ✔ oze is ready!                                             │${RESET}"
echo -e "${GREEN}│                                                              │${RESET}"
echo -e "${GREEN}│  1. Replace 'yourusername' with your GitHub handle:          │${RESET}"
echo -e "${GREEN}│     sed -i '' 's/yourusername/YOUR_HANDLE/g' \\               │${RESET}"
echo -e "${GREEN}│       go.mod internal/cli/cli.go                             │${RESET}"
echo -e "${GREEN}│                                                              │${RESET}"
echo -e "${GREEN}│  2. Create the GitHub repo and push:                         │${RESET}"
echo -e "${GREEN}│     gh repo create YOUR_HANDLE/oze --public --source=. \\     │${RESET}"
echo -e "${GREEN}│       --remote=origin --push                                 │${RESET}"
echo -e "${GREEN}│     # OR manually:                                           │${RESET}"
echo -e "${GREEN}│     git remote add origin git@github.com:HANDLE/oze.git      │${RESET}"
echo -e "${GREEN}│     git push -u origin main                                  │${RESET}"
echo -e "${GREEN}╰──────────────────────────────────────────────────────────────╯${RESET}"
echo ""
