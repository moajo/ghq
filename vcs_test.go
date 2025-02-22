package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/motemen/ghq/cmdutil"
)

var remoteDummyURL = mustParseURL("https://example.com/git/repo")

func TestVCSBackend(t *testing.T) {
	tempDir := newTempDir(t)
	defer os.RemoveAll(tempDir)
	localDir := filepath.Join(tempDir, "repo")
	_commands := []*exec.Cmd{}
	lastCommand := func() *exec.Cmd { return _commands[len(_commands)-1] }
	defer func(orig func(cmd *exec.Cmd) error) {
		cmdutil.CommandRunner = orig
	}(cmdutil.CommandRunner)
	cmdutil.CommandRunner = func(cmd *exec.Cmd) error {
		_commands = append(_commands, cmd)
		return nil
	}

	testCases := []struct {
		name   string
		f      func() error
		expect []string
		dir    string
	}{{
		name: "[git] clone",
		f: func() error {
			return GitBackend.Clone(&vcsGetOption{
				url: remoteDummyURL,
				dir: localDir,
			})
		},
		expect: []string{"git", "clone", remoteDummyURL.String(), localDir},
	}, {
		name: "[git] shallow clone",
		f: func() error {
			return GitBackend.Clone(&vcsGetOption{
				url:     remoteDummyURL,
				dir:     localDir,
				shallow: true,
				silent:  true,
			})
		},
		expect: []string{"git", "clone", "--depth", "1", remoteDummyURL.String(), localDir},
	}, {
		name: "[git] clone specific branch",
		f: func() error {
			return GitBackend.Clone(&vcsGetOption{
				url:    remoteDummyURL,
				dir:    localDir,
				branch: "hello",
			})
		},
		expect: []string{"git", "clone", "--branch", "hello", "--single-branch", remoteDummyURL.String(), localDir},
	}, {
		name: "[git] update",
		f: func() error {
			return GitBackend.Update(&vcsGetOption{
				dir: localDir,
			})
		},
		expect: []string{"git", "pull", "--ff-only"},
		dir:    localDir,
	}, {
		name: "[svn] checkout",
		f: func() error {
			return SubversionBackend.Clone(&vcsGetOption{
				url: remoteDummyURL,
				dir: localDir,
			})
		},
		expect: []string{"svn", "checkout", remoteDummyURL.String(), localDir},
	}, {
		name: "[svn] checkout shallow",
		f: func() error {
			return SubversionBackend.Clone(&vcsGetOption{
				url:     remoteDummyURL,
				dir:     localDir,
				shallow: true,
			})
		},
		expect: []string{"svn", "checkout", "--depth", "1", remoteDummyURL.String(), localDir},
	}, {
		name: "[svn] checkout specific branch",
		f: func() error {
			return SubversionBackend.Clone(&vcsGetOption{
				url:    remoteDummyURL,
				dir:    localDir,
				branch: "hello",
			})
		},
		expect: []string{"svn", "checkout", remoteDummyURL.String() + "/branches/hello", localDir},
	}, {
		name: "[svn] update",
		f: func() error {
			return SubversionBackend.Update(&vcsGetOption{
				dir:    localDir,
				silent: true,
			})
		},
		expect: []string{"svn", "update"},
		dir:    localDir,
	}, {
		name: "[git-svn] clone",
		f: func() error {
			return GitsvnBackend.Clone(&vcsGetOption{
				url: remoteDummyURL,
				dir: localDir,
			})
		},
		expect: []string{"git", "svn", "clone", remoteDummyURL.String(), localDir},
	}, {
		name: "[git-svn] update",
		f: func() error {
			return GitsvnBackend.Update(&vcsGetOption{
				dir: localDir,
			})
		},
		expect: []string{"git", "svn", "rebase"},
		dir:    localDir,
	}, {
		name: "[git-svn] clone shallow",
		f: func() error {
			return GitsvnBackend.Clone(&vcsGetOption{
				url:     remoteDummyURL,
				dir:     localDir,
				shallow: true,
			})
		},
		expect: []string{"git", "svn", "clone", remoteDummyURL.String(), localDir},
	}, {
		name: "[git-svn] clone specific branch",
		f: func() error {
			return GitsvnBackend.Clone(&vcsGetOption{
				url:    remoteDummyURL,
				dir:    localDir,
				branch: "hello",
			})
		},
		expect: []string{"git", "svn", "clone", remoteDummyURL.String() + "/branches/hello", localDir},
	}, {
		name: "[hg] clone",
		f: func() error {
			return MercurialBackend.Clone(&vcsGetOption{
				url: remoteDummyURL,
				dir: localDir,
			})
		},
		expect: []string{"hg", "clone", remoteDummyURL.String(), localDir},
	}, {
		name: "[hg] update",
		f: func() error {
			return MercurialBackend.Update(&vcsGetOption{
				dir: localDir,
			})
		},
		expect: []string{"hg", "pull", "--update"},
		dir:    localDir,
	}, {
		name: "[hg] clone shallow",
		f: func() error {
			return MercurialBackend.Clone(&vcsGetOption{
				url:     remoteDummyURL,
				dir:     localDir,
				shallow: true,
			})
		},
		expect: []string{"hg", "clone", remoteDummyURL.String(), localDir},
	}, {
		name: "[hg] clone specific branch",
		f: func() error {
			return MercurialBackend.Clone(&vcsGetOption{
				url:    remoteDummyURL,
				dir:    localDir,
				branch: "hello",
			})
		},
		expect: []string{"hg", "clone", "--branch", "hello", remoteDummyURL.String(), localDir},
	}, {
		name: "[darcs] clone",
		f: func() error {
			return DarcsBackend.Clone(&vcsGetOption{
				url: remoteDummyURL,
				dir: localDir,
			})
		},
		expect: []string{"darcs", "get", remoteDummyURL.String(), localDir},
	}, {
		name: "[darcs] clone shallow",
		f: func() error {
			return DarcsBackend.Clone(&vcsGetOption{
				url:     remoteDummyURL,
				dir:     localDir,
				shallow: true,
			})
		},
		expect: []string{"darcs", "get", "--lazy", remoteDummyURL.String(), localDir},
	}, {
		name: "[darcs] update",
		f: func() error {
			return DarcsBackend.Update(&vcsGetOption{
				dir: localDir,
			})
		},
		expect: []string{"darcs", "pull"},
		dir:    localDir,
	}, {
		name: "[bzr] clone",
		f: func() error {
			return BazaarBackend.Clone(&vcsGetOption{
				url: remoteDummyURL,
				dir: localDir,
			})
		},
		expect: []string{"bzr", "branch", remoteDummyURL.String(), localDir},
	}, {
		name: "[bzr] update",
		f: func() error {
			return BazaarBackend.Update(&vcsGetOption{
				dir: localDir,
			})
		},
		expect: []string{"bzr", "pull", "--overwrite"},
		dir:    localDir,
	}, {
		name: "[bzr] clone shallow",
		f: func() error {
			return BazaarBackend.Clone(&vcsGetOption{
				url:     remoteDummyURL,
				dir:     localDir,
				shallow: true,
			})
		},
		expect: []string{"bzr", "branch", remoteDummyURL.String(), localDir},
	}, {
		name: "[fossil] clone",
		f: func() error {
			return FossilBackend.Clone(&vcsGetOption{
				url: remoteDummyURL,
				dir: localDir,
			})
		},
		expect: []string{"fossil", "open", fossilRepoName},
		dir:    localDir,
	}, {
		name: "[fossil] update",
		f: func() error {
			return FossilBackend.Update(&vcsGetOption{
				dir: localDir,
			})
		},
		expect: []string{"fossil", "update"},
		dir:    localDir,
	}}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.f(); err != nil {
				t.Errorf("error should be nil, but: %s", err)
			}
			c := lastCommand()
			if !reflect.DeepEqual(c.Args, tc.expect) {
				t.Errorf("\ngot:  %+v\nexpect: %+v", c.Args, tc.expect)
			}
			if c.Dir != tc.dir {
				t.Errorf("got: %s, expect: %s", c.Dir, tc.dir)
			}
		})
	}
}

func TestCvsDummyBackend(t *testing.T) {
	tempDir := newTempDir(t)
	defer os.RemoveAll(tempDir)
	localDir := filepath.Join(tempDir, "repo")

	if err := cvsDummyBackend.Clone(&vcsGetOption{
		url: remoteDummyURL,
		dir: localDir,
	}); err == nil {
		t.Error("error should be occurred, but nil")
	}

	if err := cvsDummyBackend.Clone(&vcsGetOption{
		url:     remoteDummyURL,
		dir:     localDir,
		shallow: true,
	}); err == nil {
		t.Error("error should be occurred, but nil")
	}

	if err := cvsDummyBackend.Update(&vcsGetOption{
		dir: localDir,
	}); err == nil {
		t.Error("error should be occurred, but nil")
	}
}

func TestBranchOptionIgnoredErrors(t *testing.T) {
	tempDir := newTempDir(t)
	defer os.RemoveAll(tempDir)
	localDir := filepath.Join(tempDir, "repo")

	if err := DarcsBackend.Clone(&vcsGetOption{
		url:    remoteDummyURL,
		dir:    localDir,
		branch: "hello",
	}); err == nil {
		t.Error("error should be occurred, but nil")
	}

	if err := FossilBackend.Clone(&vcsGetOption{
		url:    remoteDummyURL,
		dir:    localDir,
		branch: "hello",
	}); err == nil {
		t.Error("error should be occurred, but nil")
	}

	if err := BazaarBackend.Clone(&vcsGetOption{
		url:    remoteDummyURL,
		dir:    localDir,
		branch: "hello",
	}); err == nil {
		t.Error("error should be occurred, but nil")
	}
}
