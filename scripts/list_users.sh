#!/usr/bin/env bash
# List users from the SQLite `users.db` and print first_name, member_id, point_balance
# Usage: ./scripts/list_users.sh [path/to/users.db]

set -euo pipefail

DB_PATH="${1:-./users.db}"

if [ ! -f "$DB_PATH" ]; then
  echo "Error: database file not found at '$DB_PATH'" >&2
  exit 2
fi

# Check sqlite3 is available
if ! command -v sqlite3 >/dev/null 2>&1; then
  echo "Error: sqlite3 is not installed or not in PATH" >&2
  exit 3
fi

# Query and print header
printf "%-20s %-20s %12s\n" "first_name" "member_id" "point_balance"
printf "%-20s %-20s %12s\n" "--------------------" "--------------------" "------------"

sqlite3 -readonly "$DB_PATH" <<SQL
.mode list
.separator ' | '
SELECT IFNULL(first_name, '') AS first_name,
       IFNULL(member_id, '') AS member_id,
       IFNULL(point_balance, 0) AS point_balance
FROM users
ORDER BY first_name COLLATE NOCASE ASC;
SQL

# Reformat the output from pipe-separated to aligned columns
# We read the output and align columns using awk
sqlite3 -readonly "$DB_PATH" "SELECT IFNULL(first_name, ''), IFNULL(member_id, ''), IFNULL(point_balance, 0) FROM users ORDER BY first_name COLLATE NOCASE ASC;" \
  | awk -F '\t' '{ printf "%-20s %-20s %12s\n", $1, $2, $3 }'

exit 0
