newline() {
  local target="${1:-.}"
  ensure_single_newline() {
    local file="$1"
    local tmp
    tmp=$(mktemp "${file}.XXXXXX") || { echo "mktemp failed"; return 1; }
    perl -0777 -pe 's/(\R)*\z/\n/' -- "$file" > "$tmp" && cat "$tmp" > "$file"
    rm -f "$tmp"
  }
  if [ -f "$target" ]; then
    ensure_single_newline "$target"
  elif [ -d "$target" ]; then
    find "$target" -mindepth 1 -type d -name ".*" -prune -o -type f -exec sh -c '
      file="$1"
      tmp=$(mktemp "${file}.XXXXXX") || { echo "mktemp failed"; exit 1; }
      perl -0777 -pe "s/(\\R)*\\z/\\n/" -- "$file" > "$tmp" && cat "$tmp" > "$file"
      rm -f "$tmp"
    ' _ {} \;
  else
    echo "Error: $target is not a file or directory" >&2
    return 1
  fi
}

test-newline() {
  local orig_dir=$PWD tmp_dir
  tmp_dir=$(mktemp -d) || { echo "mktemp failed"; return 1; }
  cd "$tmp_dir" || return 1

  local test_dirs=(
    testdir
    testdir/.hiddendir
    testdir/subdir
    testdir/.hiddendir/.deephidden
    testdir/subdir/.subhidden
  )
  for d in "${test_dirs[@]}"; do
    mkdir -p "$d"
    touch "$d/empty.txt"
    printf "\n"   > "$d/newline.txt"
    printf "\n\n" > "$d/extra-newline.txt"
    printf "content"    > "$d/content.txt"
    printf "content\n"  > "$d/content-newline.txt"
  done

  newline testdir
  newline testdir/.hiddendir

  echo "root of testdir:"; find testdir -maxdepth 1 -type f -exec wc -l {} +
  echo "root of testdir/subdir:"; find testdir/subdir -maxdepth 1 -type f -exec wc -l {} +
  echo "root of testdir/.hiddendir:"; find testdir/.hiddendir -maxdepth 1 -type f -exec wc -l {} +
  echo "testdir/.hiddendir/.deephidden (should NOT be processed):"; find testdir/.hiddendir/.deephidden -maxdepth 1 -type f -exec wc -l {} +
  echo "testdir/subdir/.subhidden (should NOT be processed):"; find testdir/subdir/.subhidden -maxdepth 1 -type f -exec wc -l {} +

  cd "$orig_dir"
}
