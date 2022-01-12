#!/bin/sh

echo "terminate"
ps aux | grep gim | grep -v 'grep' | awk '{print $2}' | xargs kill -TERM