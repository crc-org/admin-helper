//go:build darwin || linux

package hosts

import (
	"syscall"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDropPrivilegesFromSetuidRoot(t *testing.T) {
	var (
		groupsCleared bool
		setgidCalls   []int
		setuidCalls   []int
	)

	restoreSyscalls := stubPrivilegeSyscalls(
		func(groups []int) error {
			groupsCleared = len(groups) == 0
			return nil
		},
		func(gid int) error {
			setgidCalls = append(setgidCalls, gid)
			return nil
		},
		func(uid int) error {
			if uid == 0 {
				return syscall.EPERM
			}
			setuidCalls = append(setuidCalls, uid)
			return nil
		},
	)
	defer restoreSyscalls()

	restoreIdentity := stubProcessIdentity(501, 20, 0, 0)
	defer restoreIdentity()

	err := dropPrivileges()
	require.NoError(t, err)
	assert.True(t, groupsCleared)
	assert.Equal(t, []int{20}, setgidCalls)
	assert.Equal(t, []int{501}, setuidCalls)
}

func TestDropPrivilegesSkipsWhenAlreadyDropped(t *testing.T) {
	var setuidCalled bool

	restoreSyscalls := stubPrivilegeSyscalls(
		func([]int) error { return nil },
		func(int) error { return nil },
		func(int) error {
			setuidCalled = true
			return nil
		},
	)
	defer restoreSyscalls()

	restoreIdentity := stubProcessIdentity(501, 20, 501, 20)
	defer restoreIdentity()

	err := dropPrivileges()
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

func stubProcessIdentity(uid, gid, euid, egid int) func() {
	origGetuid, origGetgid, origGeteuid, origGetegid := sysGetuid, sysGetgid, sysGeteuid, sysGetegid
	sysGetuid = func() int { return uid }
	sysGetgid = func() int { return gid }
	sysGeteuid = func() int { return euid }
	sysGetegid = func() int { return egid }
	return func() {
		sysGetuid, sysGetgid, sysGeteuid, sysGetegid = origGetuid, origGetgid, origGeteuid, origGetegid
	}
}
