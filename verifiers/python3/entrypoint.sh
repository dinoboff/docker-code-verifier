#!/bin/bash

if [ "$1" = 'bash' ] || [ "$1" = 'python3' ]; then
	exec "$@"
else
	exec python3 run-problem.py "$@"    
fi
