#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

echo "=========================================="
echo "Building Frontend..."
echo "=========================================="

cd frontend
# Install dependencies
echo "Installing frontend dependencies..."
npm install

# Build frontend
echo "Building frontend application..."
npm run build

if [ $? -ne 0 ]; then
    echo "Frontend build failed!"
    exit 1
fi

cd ..

echo "=========================================="
echo "Updating Backend Static Files..."
echo "=========================================="

# Remove old static files
echo "Cleaning up old static files..."
rm -rf internal/static/dist

# Copy new static files
echo "Copying new static files..."
cp -r frontend/dist internal/static/

if [ $? -ne 0 ]; then
    echo "Failed to copy static files!"
    exit 1
fi

echo "=========================================="
echo "Building Backend..."
echo "=========================================="

echo "Tidying go modules..."
go mod tidy

echo "Building backend binary..."
go build -o ai_novel .

if [ $? -ne 0 ]; then
    echo "Backend build failed!"
    exit 1
fi

echo "=========================================="
echo "Starting Server..."
echo "=========================================="

# Make sure the binary is executable
chmod +x ai_novel

./ai_novel
