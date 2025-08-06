#!/bin/bash

# Example usage script for detect-server

echo "=== Server Detection Tool Examples ==="

# Quick detection with stdout output (default)
echo "1. Quick detection with stdout output:"
./detect-server

# More comprehensive detection with HTML output
echo -e "\n2. More comprehensive detection with HTML output:"
./detect-server -c more -s html

# Quick detection with PDF output
echo -e "\n3. Quick detection with PDF output:"
./detect-server -s pdf

# More detection with stdout output
echo -e "\n4. More detection with stdout output:"
./detect-server -c more

echo -e "\n=== End of Examples ==="