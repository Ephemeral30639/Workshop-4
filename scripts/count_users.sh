#!/usr/bin/env bash
# count_users.sh - Count rows in the `users` table of an SQLite database.
# Usage: ./count_users.sh [path/to/users.db]
# Defaults to ./users.db when no argument supplied.
set -euo pipefail
DB_PATH=${1:-"$(pwd)/users.db"}
SQL='SELECT COUNT(*) AS count FROM users;'

# Check file exists
if [ ! -f "$DB_PATH" ]; then
  echo "Error: database file not found at '$DB_PATH'" >&2
  exit 2
fi

# Try sqlite3 -json first (if available), otherwise fall back to plain output
if sqlite3 --help 2>&1 | grep -q -- "-json"; then
  # Use json output and extract the count field (jq not required)
  sqlite3 -json "$DB_PATH" "$SQL" 2>/dev/null || {
    echo "Error: sqlite3 query failed" >&2
    exit 3
  }
else
  # Fallback: output a single number
  # Use -noheader and -csv to get clean output, then strip possible whitespace
  count=$(sqlite3 -noheader -csv "$DB_PATH" "$SQL" 2>/dev/null || { echo ""; })
  if [ -z "$count" ]; then
    echo "Error: sqlite3 query failed or returned no rows" >&2
    exit 3
  fi
  # count may be like "42" or "42\n"; ensure trimmed
  # Remove possible CSV quoting
  count=${count%\"}
  count=${count#\"}
  echo "$count"
fi
