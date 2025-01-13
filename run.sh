#!/usr/bin/env bash

npx tailwindcss -i ./api/static/input.css -o ./api/static/styles.css --verbose

go run .


