#!/usr/bin/env bash
set -euo pipefail

BASE_URL="http://localhost:8080/api/v1"
PROJECT_ID=""
EMAIL=""
PASSWORD=""

usage() {
    cat <<EOF
Usage:
  $(basename "$0") --project-id <id> --email <email> --password <password> [--base-url <url>]

Example:
  $(basename "$0") --project-id 1 --email demo@example.com --password 'your-password'
EOF
}

while [[ $# -gt 0 ]]; do
    case "$1" in
        --project-id)
            PROJECT_ID="${2:-}"
            shift 2
            ;;
        --email)
            EMAIL="${2:-}"
            shift 2
            ;;
        --password)
            PASSWORD="${2:-}"
            shift 2
            ;;
        --base-url)
            BASE_URL="${2:-}"
            shift 2
            ;;
        -h|--help)
            usage
            exit 0
            ;;
        *)
            echo "Unknown argument: $1" >&2
            usage
            exit 1
            ;;
    esac
done

if [[ -z "$PROJECT_ID" || -z "$EMAIL" || -z "$PASSWORD" ]]; then
    echo "Missing required arguments." >&2
    usage
    exit 1
fi

COOKIE_JAR="$(mktemp)"
cleanup() {
    rm -f "$COOKIE_JAR"
}
trap cleanup EXIT

echo "Logging in as $EMAIL ..."
LOGIN_STATUS="$(
    curl -sS -o /dev/null -w "%{http_code}" \
        -c "$COOKIE_JAR" \
        -H "Content-Type: application/json" \
        -X POST "${BASE_URL}/auth/login" \
        -d "{\"email\":\"${EMAIL}\",\"password\":\"${PASSWORD}\"}"
)"

if [[ "$LOGIN_STATUS" != "200" ]]; then
    echo "Login failed with status: $LOGIN_STATUS" >&2
    exit 1
fi

post_item() {
    local endpoint="$1"
    local label="$2"
    local color="$3"

    local status
    status="$(
        curl -sS -o /tmp/ticker-seed-response.json -w "%{http_code}" \
            -b "$COOKIE_JAR" \
            -H "Content-Type: application/json" \
            -X POST "${BASE_URL}${endpoint}" \
            -d "{\"label\":\"${label}\",\"color\":\"${color}\"}"
    )"

    if [[ "$status" == "201" ]]; then
        echo "Created: ${label}"
        return
    fi

    if [[ "$status" == "400" ]]; then
        echo "Skipped (already exists/invalid): ${label}"
        return
    fi

    echo "Failed (${status}): ${label}" >&2
    cat /tmp/ticker-seed-response.json >&2 || true
}

echo "Seeding statuses ..."
post_item "/projects/${PROJECT_ID}/statuses" "Backlog" "#64748B"
post_item "/projects/${PROJECT_ID}/statuses" "Todo" "#3B82F6"
post_item "/projects/${PROJECT_ID}/statuses" "In Progress" "#F59E0B"
post_item "/projects/${PROJECT_ID}/statuses" "Blocked" "#EF4444"
post_item "/projects/${PROJECT_ID}/statuses" "In Review" "#8B5CF6"
post_item "/projects/${PROJECT_ID}/statuses" "Done" "#22C55E"

echo "Seeding priorities ..."
post_item "/projects/${PROJECT_ID}/priorities" "Low" "#22C55E"
post_item "/projects/${PROJECT_ID}/priorities" "Medium" "#F59E0B"
post_item "/projects/${PROJECT_ID}/priorities" "High" "#EF4444"
post_item "/projects/${PROJECT_ID}/priorities" "Critical" "#DC2626"

echo "Seeding tags ..."
post_item "/projects/${PROJECT_ID}/tags" "frontend" "#38BDF8"
post_item "/projects/${PROJECT_ID}/tags" "backend" "#34D399"
post_item "/projects/${PROJECT_ID}/tags" "bug" "#F87171"
post_item "/projects/${PROJECT_ID}/tags" "feature" "#A78BFA"
post_item "/projects/${PROJECT_ID}/tags" "performance" "#FBBF24"
post_item "/projects/${PROJECT_ID}/tags" "refactor" "#94A3B8"
post_item "/projects/${PROJECT_ID}/tags" "docs" "#60A5FA"
post_item "/projects/${PROJECT_ID}/tags" "testing" "#F472B6"

echo "Done."
