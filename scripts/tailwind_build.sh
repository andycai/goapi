#!/bin/bash

NODE_ENV=production npx @tailwindcss/cli -i ./public/css/input.css -o ./public/css/tailwind.min.css -m -w