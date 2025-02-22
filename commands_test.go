package main

import (
	"net/url"
	"os"
	"path/filepath"
	"testing"
)

type _cloneArgs struct {
	remote  *url.URL
	local   string
	shallow bool
	branch  string
}

type _updateArgs struct {
	local string
}

func withFakeGitBackend(t *testing.T, block func(*testing.T, string, *_cloneArgs, *_updateArgs)) {
	tmpRoot := newTempDir(t)
	defer os.RemoveAll(tmpRoot)

	defer func(orig []string) { _localRepositoryRoots = orig }(_localRepositoryRoots)
	_localRepositoryRoots = []string{tmpRoot}

	var cloneArgs _cloneArgs
	var updateArgs _updateArgs

	var originalGitBackend = GitBackend
	tmpBackend := &VCSBackend{
		Clone: func(vg *vcsGetOption) error {
			cloneArgs = _cloneArgs{
				remote:  vg.url,
				local:   filepath.FromSlash(vg.dir),
				shallow: vg.shallow,
				branch:  vg.branch,
			}
			return nil
		},
		Update: func(vg *vcsGetOption) error {
			updateArgs = _updateArgs{
				local: vg.dir,
			}
			return nil
		},
	}
	GitBackend = tmpBackend
	vcsContentsMap[".git"] = tmpBackend
	defer func() { GitBackend = originalGitBackend; vcsContentsMap[".git"] = originalGitBackend }()

	block(t, tmpRoot, &cloneArgs, &updateArgs)
}
