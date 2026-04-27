package gitx

import "testing"

func TestDeriveCloneDir(t *testing.T) {
	cases := []struct {
		in, want string
	}{
		{"https://github.com/foo/bar.git", "bar"},
		{"https://github.com/foo/bar", "bar"},
		{"https://github.com/foo/bar/", "bar"},
		{"git@github.com:foo/bar.git", "bar"},
		{"git@github.com:foo/bar", "bar"},
		{"/local/path/to/repo.git/", "repo"},
		{"/local/path/to/repo", "repo"},
		{"repo", "repo"},
		{"", ""},
	}
	for _, tc := range cases {
		if got := DeriveCloneDir(tc.in); got != tc.want {
			t.Errorf("DeriveCloneDir(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}