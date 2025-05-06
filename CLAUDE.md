# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Envault is a CLI tool designed to securely manage environment variables by encrypting `.env` files. It creates encrypted `.env.vaulted` files and provides commands to export/unset environment variables when needed, enhancing security for sensitive configuration data.

## Key Commands

### Building

```bash
# Install dependencies
go mod tidy

# Build the application
go build -o envault cmd/envault/main.go

# Run tests
go test ./...
```

### Encryption

```bash
# Encrypt .env file (creates .env.vaulted)
./envault encrypt .env

# Encrypt with custom output path
./envault encrypt .env -f /path/to/output.vaulted
```

### Exporting Environment Variables

```bash
# Script evaluation method (variables set in current shell)
eval $(./envault export -o)
source <(./envault export --output-script-only)

# New shell method (variables set in new bash session)
./envault export -n
./envault export --new-shell

# Command execution method (variables set for specific command)
./envault export -- node app.js
./envault export -- docker-compose up
```

### Unsetting Environment Variables

```bash
# Unset variables in current shell
eval $(./envault unset -o)
source <(./envault unset --output-script-only)
```

### Viewing Encrypted Content

```bash
# View content of encrypted file
./envault dump

# Save decrypted content to file
./envault dump > decrypted.env
```

### Interactive Selection

```bash
# Select specific environment variables to export
./envault export -s
./envault export select

# Select variables and run in new shell
./envault export -s -n

# Select variables and run command
./envault export -s -- npm start
```

## Architecture

Envault is structured into several key components:

1. **CLI Interface** (`internal/cli/cli.go`)
   - Manages command structure: encrypt, export, unset, dump
   - Handles command-line argument parsing and workflow management

2. **Cryptography** (`internal/crypto/crypto.go`)
   - Implements AES-256-GCM encryption for confidentiality and integrity
   - Uses Argon2id for secure password derivation
   - Handles the secure file format with magic bytes, salt, and nonce

3. **Environment Variable Management** (`internal/env/`)
   - Parses and processes environment variables
   - Generates export/unset scripts

4. **File Operations** (`internal/file/file.go`)
   - Manages reading/writing encrypted files
   - Handles file format validation

5. **TUI (Terminal User Interface)** (`internal/tui/`)
   - Provides interactive selection of environment variables
   - Based on BubbleTea library

6. **Utilities** (`pkg/utils/utils.go`)
   - Handles password input (interactive and stdin)
   - Contains script execution utilities

## Technical Constraints

Envault works around OS security model limitations where:
- Child processes cannot modify the environment of their parent process
- To set environment variables in the current shell, Envault must generate shell scripts to be evaluated with `eval` or `source`

The project implements three methods to handle this constraint:
1. Script evaluation method with `eval` or `source`
2. New shell method that launches a new bash session with variables set
3. Command execution method that runs a specified command with variables set

## Security Features

- AES-256-GCM authenticated encryption
- Argon2id for key derivation with secure parameters
- Random salt and nonce for each encryption
- No plaintext stored in encrypted files
- Password handling with memory zeroing after use