#!/bin/bash

# Colors for pretty output
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}Building aishell...${NC}"

# Ensure we're in the project root
cd "$(dirname "$0")/.."

# Build the binary
go build -o aishell ./cmd/cli

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Build successful!${NC}"
    echo
    echo -e "${BLUE}Next steps:${NC}"
    echo "1. Create a .env file in the project root with your OpenRouter API key:"
    echo "   echo 'OPENROUTER_API_KEY=your-api-key-here' > .env"
    echo
    echo "2. Make the binary accessible from anywhere (optional):"
    echo "   sudo mv aishell /usr/local/bin/"
    echo
    echo -e "${BLUE}Example usage:${NC}"
    echo "# Interactive mode"
    echo "./aishell -i"
    echo
    echo "# Direct queries"
    echo "./aishell \"list all pdf files\""
    echo "./aishell \"find large files over 1GB\""
    echo
    echo "# Execute commands directly"
    echo "./aishell -e \"show system memory usage\""
    echo
    echo "# Verbose mode"
    echo "./aishell -v \"count files in current directory\""
    echo
    echo -e "${BLUE}For more information:${NC}"
    echo "./aishell --help"
else
    echo "❌ Build failed"
    exit 1
fi