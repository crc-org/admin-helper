# Testing the Sandboxed admin-helper

## Quick Test

You need to run this with sudo since it modifies /etc/hosts:

```bash
# Run the automated test script
./test_sandbox.sh
```

This will test add/remove/clean operations with the sandbox enabled.

## Verifying the Sandbox Works

### Check the binary is CGO-enabled
```bash
file ./crc-admin-helper
# Output: crc-admin-helper: Mach-O 64-bit executable arm64

otool -L ./crc-admin-helper | grep Security
# Output: /System/Library/Frameworks/Security.framework/...

nm ./crc-admin-helper | grep sandbox
# Output shows: _sandbox_init, _sandbox_free_error
```

### Test sandbox restrictions
```bash
CGO_ENABLED=1 go build -o verify_sandbox verify_sandbox.go
./verify_sandbox
```

Expected output:
- ✓ Can read `/etc/hosts` and `/private/etc/hosts`
- ✗ Cannot read `/etc/ssh/ssh_config`, `~/.ssh/id_rsa`, `/Users`

### Manual functional test
```bash
# Add entry
sudo ./crc-admin-helper add 192.168.130.11 api.crc.testing console-openshift-console.apps-crc.testing

# Verify
grep "Added by CRC" -A10 /etc/hosts

# Remove
sudo ./crc-admin-helper remove api.crc.testing

# Clean up
sudo ./crc-admin-helper clean
```

## Security Benefits

The sandbox restricts admin-helper to only:
1. Read/write `/etc/hosts`
2. Read system libraries (for Go runtime)
3. Access temp directories (for atomic writes)

It denies:
- Network access
- Executing other programs
- Reading/writing other files (SSH keys, documents, etc.)

This provides defense-in-depth: even if admin-helper has a vulnerability and is running as root, the attacker cannot:
- Exfiltrate data over the network
- Read SSH keys or other sensitive files
- Execute backdoors or other binaries

## Troubleshooting

### Build fails with "library 'crt0.o' not found"
- Make sure you removed `-static` from LDFLAGS for macOS builds
- CGO doesn't support static linking on macOS

### Permission errors accessing /etc/hosts
- Make sure sandbox profile allows `/etc/hosts` and `/private/etc/hosts`
- Check Console.app for "Sandbox: deny file-write" messages
