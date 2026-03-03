#!/bin/bash

# Configuration variables
LOCAL_BRANCH="rewrite"
ORIGIN_REMOTE="origin"
UPSTREAM_REMOTE="upstream"
UPSTREAM_BRANCH="main-rewrite"
UPSTREAM_URL="https://github.com/rancher/terraform-provider-rancher2.git"

echo "================================================="
echo " Starting upstream rebase process..."
echo " Local Branch: $LOCAL_BRANCH"
echo " Upstream Remote: $UPSTREAM_REMOTE/$UPSTREAM_BRANCH"
echo "================================================="

# 1. Verify and add upstream remote if it doesn't exist
echo "[1/4] Checking for upstream remote '$UPSTREAM_REMOTE'..."
if ! git remote get-url "$UPSTREAM_REMOTE" > /dev/null 2>&1; then
    echo "      Remote '$UPSTREAM_REMOTE' not found. Adding it now..."
    if ! git remote add "$UPSTREAM_REMOTE" "$UPSTREAM_URL"; then
        echo "[ERROR] Failed to add remote '$UPSTREAM_REMOTE'."
        exit 1
    fi
    echo "      Successfully added upstream remote: $UPSTREAM_URL"
fi

# 2. Fetch the latest changes from the upstream remote
echo "[2/4] Fetching latest changes from $UPSTREAM_REMOTE..."
if ! git fetch "$UPSTREAM_REMOTE"; then
    echo "[ERROR] Failed to fetch from remote '$UPSTREAM_REMOTE'."
    exit 1
fi

# 3. Checkout the local branch
echo "[3/4] Checking out local branch '$LOCAL_BRANCH'..."
if ! git checkout "$LOCAL_BRANCH"; then
    echo "[ERROR] Failed to checkout branch '$LOCAL_BRANCH'. Make sure it exists locally."
    exit 1
fi

# 4. Perform the rebase
echo "[4/4] Rebasing '$LOCAL_BRANCH' onto '$UPSTREAM_REMOTE/$UPSTREAM_BRANCH'..."
if git rebase "$UPSTREAM_REMOTE/$UPSTREAM_BRANCH"; then
    echo "[SUCCESS] Rebase completed successfully!"
else
    echo "[WARNING] Rebase halted due to conflicts."
    echo ""
    echo "To proceed:"
    echo "  1. Resolve the conflicts in your editor."
    echo "  2. Stage the resolved files: git add <file>"
    echo "  3. Continue the rebase:      git rebase --continue"
    echo ""
    echo "Or, to abort the rebase entirely: git rebase --abort"
    exit 1
fi

echo "================================================="
echo " Done."
echo " Note: Since you rewrote local history, you will need to force push"
echo " your updated branch back to your fork (origin):"
echo " "
echo "   git push --force-with-lease $ORIGIN_REMOTE $LOCAL_BRANCH"
echo "================================================="
