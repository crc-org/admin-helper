//go:build darwin || linux
// +build darwin linux

package hosts

import (
	"errors"
	"os"
	"path/filepath"
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNeedsPrivilegeDrop(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		targetUID int
		targetGID int
		ok        bool
		euid      int
		egid      int
		uid       int
		wantDrop  bool
	}{
		{
			name: "real root process",
			ok:   false,
			euid: 0,
			uid:  0,
		},
		{
			name: "unprivileged process",
			ok:   true,
			euid: 501,
			egid: 20,
			uid:  501,
		},
		{
			name:      "already dropped to target identity",
			ok:        true,
			targetUID: 501,
			targetGID: 20,
			euid:      501,
			egid:      20,
			uid:       501,
		},
		{
			name:      "setuid root helper",
			ok:        true,
			targetUID: 501,
			targetGID: 20,
			euid:      0,
			egid:      0,
			uid:       501,
			wantDrop:  true,
		},
		{
			name:      "setuid root with mismatched effective gid",
			ok:        true,
			targetUID: 501,
			targetGID: 20,
			euid:      0,
			egid:      0,
			uid:       501,
			wantDrop:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.wantDrop, needsPrivilegeDrop(tt.targetUID, tt.targetGID, tt.ok, tt.euid, tt.egid, tt.uid))
		})
	}
}

func TestDropPrivilegesFromSetuidRoot(t *testing.T) {
	var (
		groupsCleared bool
		setgidCalls   []int
		setuidCalls   []int
	)

	restore := stubPrivilegeSyscalls(
		func(groups []int) error {
			groupsCleared = len(groups) == 0
			return nil
		},
		func(gid int) error {
			setgidCalls = append(setgidCalls, gid)
			return nil
		},
		func(uid int) error {
			setuidCalls = append(setuidCalls, uid)
			if uid == 0 {
				return syscall.EPERM
			}
			return nil
		},
	)
	defer restore()

	err := dropPrivilegesWithIdentity(501, 20, true, 0, 0, 501)
	require.NoError(t, err)
	assert.True(t, groupsCleared)
	assert.Equal(t, []int{20}, setgidCalls)
	assert.Equal(t, []int{501, 0}, setuidCalls)
}

func TestDropPrivilegesSkipsWhenUnprivileged(t *testing.T) {
	var setuidCalled bool

	restore := stubPrivilegeSyscalls(
		func([]int) error { return nil },
		func(int) error { return nil },
		func(int) error {
			setuidCalled = true
			return nil
		},
	)
	defer restore()

	err := dropPrivilegesWithIdentity(501, 20, true, 501, 20, 501)
	require.NoError(t, err)
	assert.False(t, setuidCalled)
}

func TestApplyPrivilegeDropErrors(t *testing.T) {
	tests := []struct {
		name      string
		setgroups func([]int) error
		setgid    func(int) error
		setuid    func(int) error
		wantErr   string
	}{
		{
			name: "setgroups failure",
			setgroups: func([]int) error {
				return syscall.EPERM
			},
			setgid:  func(int) error { return nil },
			setuid:  func(int) error { return nil },
			wantErr: "failed to drop supplementary groups",
		},
		{
			name:    "setgid failure",
			setgid:  func(int) error { return syscall.EPERM },
			setuid:  func(int) error { return nil },
			wantErr: "failed to setgid(20)",
		},
		{
			name:    "setuid failure",
			setgid:  func(int) error { return nil },
			setuid:  func(int) error { return syscall.EPERM },
			wantErr: "failed to setuid(501)",
		},
		{
			name:    "can regain root",
			setgid:  func(int) error { return nil },
			setuid:  func(int) error { return nil },
			wantErr: "able to regain root after dropping privileges",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setgroups := tt.setgroups
			if setgroups == nil {
				setgroups = func([]int) error { return nil }
			}

			restore := stubPrivilegeSyscalls(setgroups, tt.setgid, tt.setuid)
			defer restore()

			err := applyPrivilegeDrop(501, 20)
			require.Error(t, err)
			assert.Contains(t, err.Error(), tt.wantErr)
		})
	}
}

func TestOpenHostsFileAndDropPrivilegesClosesFileOnDropFailure(t *testing.T) {
	hostsFile := filepath.Join(t.TempDir(), "hosts")
	require.NoError(t, os.WriteFile(hostsFile, []byte("keep"), 0600))

	orig := dropPrivilegesFn
	dropPrivilegesFn = func() error { return errors.New("drop failed") }
	defer func() { dropPrivilegesFn = orig }()

	file, err := OpenHostsFileAndDropPrivileges(hostsFile)
	require.Error(t, err)
	assert.Nil(t, file)
	assert.Contains(t, err.Error(), "drop failed")
	contents, readErr := os.ReadFile(hostsFile)
	require.NoError(t, readErr)
	assert.Equal(t, "keep", string(contents))
}

func TestOpenHostsFileAndDropPrivilegesInvalidPath(t *testing.T) {
	_, err := OpenHostsFileAndDropPrivileges("/nonexistent/path/to/hosts")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "failed to open hosts file")
}

func TestPrivilegeDropTarget(t *testing.T) {
	uid, gid, ok := privilegeDropTarget()

	if syscall.Getuid() == 0 {
		assert.False(t, ok)
		return
	}

	assert.True(t, ok)
	assert.Equal(t, syscall.Getuid(), uid)
	assert.Equal(t, syscall.Getgid(), gid)
}

func stubPrivilegeSyscalls(
	setgroups func([]int) error,
	setgid func(int) error,
	setuid func(int) error,
) func() {
	origSetgroups, origSetgid, origSetuid := sysSetgroups, sysSetgid, sysSetuid
	sysSetgroups, sysSetgid, sysSetuid = setgroups, setgid, setuid
	return func() {
		sysSetgroups, sysSetgid, sysSetuid = origSetgroups, origSetgid, origSetuid
	}
}
