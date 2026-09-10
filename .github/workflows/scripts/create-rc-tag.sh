#!/bin/bash
BASE_VERSION=$(echo "$GITHUB_REF" | sed 's/refs\/heads\/release\///')
echo "Base version is: $BASE_VERSION"
git fetch --tags
LATEST_MINOR_NUM=$(git tag | grep "^${BASE_VERSION}\." | awk -F. '{print $2}' | sort -n | tail -1)
echo "Latest minor num: $LATEST_MINOR_NUM"
LATEST_PATCH_NUM=$(git tag | grep "^${BASE_VERSION}\.${LATEST_MINOR_NUM}\." | awk -F- '{print $1}' | awk -F. '{print $3}' | sort -n | tail -1)
echo "Latest patch num: $LATEST_PATCH_NUM"
LATEST_RC_NUM=$(git tag | grep "^${BASE_VERSION}\.${LATEST_MINOR_NUM}\.${LATEST_PATCH_NUM}-rc\." | awk -F. '{print $4}' | grep -v '^$' | sort -n | tail -1)
echo "Latest rc num: $LATEST_RC_NUM"
LATEST_FULL_RELEASE=$(git tag | grep "^${BASE_VERSION}\.${LATEST_MINOR_NUM}\.${LATEST_PATCH_NUM}$")
echo "Detected a full release for latest rc: $LATEST_FULL_RELEASE"
if [ -n "$LATEST_FULL_RELEASE" ]; then
  if [ -z "$LATEST_MINOR_NUM" ]; then
    LATEST_MINOR_NUM=0
  else
    LATEST_MINOR_NUM=$((LATEST_MINOR_NUM + 1))
  fi
  LATEST_PATCH_NUM=0
  NEXT_RC_NUM=0
else
  LATEST_MINOR_NUM=${LATEST_MINOR_NUM:-0}
  LATEST_PATCH_NUM=${LATEST_PATCH_NUM:-0}
  if [ -z "$LATEST_RC_NUM" ]; then
    NEXT_RC_NUM=0
  else
    NEXT_RC_NUM=$((LATEST_RC_NUM + 1))
  fi
fi
NEXT_RC_TAG="${BASE_VERSION}.${LATEST_MINOR_NUM}.${LATEST_PATCH_NUM}-rc.${NEXT_RC_NUM}"
echo "Calculated next RC tag: $NEXT_RC_TAG"
git config user.name "${GITHUB_ACTOR}"
git config user.email "${GITHUB_ACTOR}@users.noreply.github.com"
git tag "$NEXT_RC_TAG" -m "Release Candidate $NEXT_RC_TAG"
git push origin "$NEXT_RC_TAG"
echo "RC_TAG=$NEXT_RC_TAG" >> "$GITHUB_ENV"
echo "RC_BRANCH=release/$BASE_VERSION" >> "$GITHUB_ENV"
