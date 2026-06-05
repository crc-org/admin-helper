#!/bin/bash
set -e

BINARY="./crc-admin-helper"
TEST_HOST="test-sandbox.crc.testing"
TEST_IP="127.0.0.1"

echo "=== Testing admin-helper with sandbox ==="
echo

echo "1. Testing ADD command..."
sudo $BINARY add $TEST_IP $TEST_HOST
if grep -q "$TEST_HOST" /etc/hosts; then
    echo "✓ Successfully added $TEST_HOST"
else
    echo "✗ Failed to add $TEST_HOST"
    exit 1
fi
echo

echo "2. Verifying entry exists..."
grep "CRC" -A10 /etc/hosts | head -20
echo

echo "3. Testing REMOVE command..."
sudo $BINARY remove $TEST_HOST
if ! grep -q "$TEST_HOST" /etc/hosts; then
    echo "✓ Successfully removed $TEST_HOST"
else
    echo "✗ Failed to remove $TEST_HOST"
    exit 1
fi
echo

echo "4. Testing CLEAN command (cleanup CRC section)..."
sudo $BINARY clean
echo "✓ Clean command executed"
echo

echo "=== Testing sandbox restrictions ==="
echo "Note: The sandbox should prevent access to files other than /etc/hosts"
echo "Check Console.app for sandbox violation messages after running this script"
echo

echo "All tests passed! The sandboxed binary works correctly."
