#!/bin/sh

set -e

if [ "$1" = 'git-robot' -a "$(id -u)" = '0' ]; then
    exec su-exec zhimiao "$0" "$@"
fi

exec $@