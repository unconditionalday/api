#!/bin/bash

# Check dependencies
if ! command -v git &> /dev/null; then
    echo "Error: git is not installed. Please install git before running the script."
    exit 1
fi

if ! command -v make &> /dev/null; then
    echo "Error: make is not installed. Please install make before running the script."
    exit 1
fi

# Temporary directory
tmp_dir="/tmp/pgvector"

# Clean up any previous residue
rm -rf "$tmp_dir"

# Clone the repository
echo "Cloning pgvector..."
git clone --branch v0.5.1 https://github.com/pgvector/pgvector.git "$tmp_dir" || exit 1

# Change directory
cd "$tmp_dir" || exit 1

# Compilation
echo "Compiling pgvector..."
make || exit 1

# Installation (may require sudo)
echo "Installing pgvector..."
make install || exit 1

echo "Installation completed successfully."
