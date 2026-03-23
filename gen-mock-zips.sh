#!/usr/bin/env bash
set -euo pipefail

# Run this script from the repo root.

need_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "error: required command not found: $1" >&2
    exit 1
  fi
}

need_cmd zip
need_cmd unzip
need_cmd dd
need_cmd mktemp

ROOT="$(pwd)"

slugs=(
  example-game
  downloading-game
  extracting-game
  not-installed-game
  error-game
  not-installed-game-2
  installed-game
  z-alternating
)

make_tiny_zip() {
  local slug="$1"
  local out="${ROOT}/${slug}-mock.zip"
  local tmp
  tmp="$(mktemp -d)"
  trap 'rm -rf "$tmp"' RETURN

  : > "${tmp}/${slug}.exe"
  (cd "$tmp" && zip -q -9 "$out" "${slug}.exe")

  trap - RETURN
  rm -rf "$tmp"
}

make_downloading_zip() {
  local slug="downloading-game"
  local out="${ROOT}/${slug}-mock.zip"
  local tmp
  tmp="$(mktemp -d)"
  trap 'rm -rf "$tmp"' RETURN

  : > "${tmp}/${slug}.exe"

  # ~2GB payload so download progress is visible. Using /dev/zero keeps generation fast.
  dd if=/dev/zero of="${tmp}/payload.bin" bs=1M count=2048 status=none

  # Store (no compression) so archive size ~= payload size.
  (cd "$tmp" && zip -q -0 "$out" "${slug}.exe" payload.bin)

  trap - RETURN
  rm -rf "$tmp"
}

make_extracting_zip() {
  local slug="extracting-game"
  local out="${ROOT}/${slug}-mock.zip"
  local tmp
  tmp="$(mktemp -d)"
  trap 'rm -rf "$tmp"' RETURN

  : > "${tmp}/${slug}.exe"

  mkdir -p "${tmp}/files"
  for i in $(seq -w 1 5000); do
    : > "${tmp}/files/${i}.dat"
  done

  # ~300MB incompressible file so download progress is visible.
  dd if=/dev/urandom of="${tmp}/payload.bin" bs=1M count=300 status=none

  # Store (no compression) so archive size ~= payload size.
  (cd "$tmp" && zip -q -0 "$out" "${slug}.exe" payload.bin files/*.dat)

  trap - RETURN
  rm -rf "$tmp"
}

echo "Generating mock zips in: ${ROOT}"
echo

echo "Rebuilding: downloading-game-mock.zip (~2GB)"
make_downloading_zip

make_extracting_zip

for slug in "${slugs[@]}"; do
  if [[ "$slug" == "downloading-game" || "$slug" == "extracting-game" ]]; then
    continue
  fi
  make_tiny_zip "$slug"
done

echo
echo "Created/updated:"
for slug in "${slugs[@]}"; do
  echo "  - ${slug}-mock.zip"
done

echo
echo "Zip contents summary:"
echo "  - For every slug, the zip must contain: /<slug>.exe (at zip root)."
echo "  - downloading-game-mock.zip additionally contains:"
echo "      - /payload.bin (~2GB, stored)"
echo "  - extracting-game-mock.zip additionally contains:"
echo "      - /payload.bin (~300MB, stored)"
echo "      - /files/0001.dat ... /files/5000.dat (0 bytes each)"
echo
echo "To verify a zip layout:"
echo "  unzip -l ./extracting-game-mock.zip | head"
echo
echo "Done."
