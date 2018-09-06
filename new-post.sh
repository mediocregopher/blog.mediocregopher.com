#!/bin/sh

set -e

numargs=2
function usage {
    echo "Usage: $0 [options] <post title> <post description>
Options:
    -i                    Create image directory
    -d \"YYYY-MM-DD\"       Custom date to use instead of today
    -V                    Verbose
    -x                    Dry run, don't make any changes
"
    exit 1
}

td=$(date "+%Y-%m-%d")

while [ "$(echo $1 | head -c1)" = '-' -o "$#" -gt $numargs ]; do
    arg="$1"
    shift

    case "$arg" in
    "-i") IMG_DIR=1;;
    "-d") td=$1; shift;;
    "-V") VERBOSE=1;;
    "-x") DRY_RUN=1;;
    *)
        echo "Unknown option '$arg'"
        usage;;
    esac
done

if [ "$#" != $numargs ]; then usage; fi

if [ ! -z $VERBOSE ]; then set -x; fi

title="$1"
clean_title=$(echo "$title" |\
    tr '[:upper:]' '[:lower:]' |\
    sed 's/[^a-z0-9 ]//g' |\
    tr ' ' '-' \
    )

description="$2"
if $(echo "$description" | grep -q '[^.$!]$'); then
    echo 'Description needs to be a complete sentence, with ending punctuation.'
    exit 1
fi

postFileName=_posts/$td-$clean_title.md
echo "Creating $postFileName"
postContent=$(cat <<EOF
---
title: >-
    $title
description: >-
    $description
---

Write stuff here, title will automatically be added as an h1

## Secondary header

Title is already h1 so all sub-titles should be h2 or below.
EOF
)

if [ -z $DRY_RUN ]; then
    echo "$postContent" > "$postFileName"
fi

if [ ! -z $IMG_DIR ]; then
    imgDirName="img/$clean_title"
    echo "Creating directory $imgDirName"
    if [ -z $DRY_RUN ]; then
        mkdir -p "$imgDirName"
    fi
fi
