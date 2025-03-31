# envault

A CLI tool for securely managing environment variables.

## Features

- Encrypt `.env` files into `.env.vaulted` files
- Export environment variables from encrypted `.env.vaulted` files
- Unset environment variables that were exported from `.env.vaulted` files

## Usage

```
# Encrypt a .env file
envault .env

# Export environment variables from a .env.vaulted file
envault export

# Unset environment variables from a .env.vaulted file
envault unset
```

## Requirements

- Linux or macOS with bash
