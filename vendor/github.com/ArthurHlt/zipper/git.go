package zipper

import (
	"fmt"
	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing"
	"gopkg.in/src-d/go-git.v4/plumbing/transport"
	"gopkg.in/src-d/go-git.v4/plumbing/transport/client"
	githttp "gopkg.in/src-d/go-git.v4/plumbing/transport/http"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

type GitHandler struct {
	client *http.Client
}

func NewGitHandler() *GitHandler {
	customClient := &http.Client{
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
		},
	}
	client.InstallProtocol("https", githttp.NewClient(customClient))
	client.InstallProtocol("http", githttp.NewClient(customClient))
	return &GitHandler{customClient}
}
func (h GitHandler) Zip(src *Source) (ZipReadCloser, error) {
	h.setHttpClient(src)
	path := src.Path
	tmpDir, err := ioutil.TempDir("", "git-zipper")
	if err != nil {
		return nil, err
	}
	gitUtils := h.makeGitUtils(tmpDir, path)
	err = gitUtils.Clone()
	if err != nil {
		return nil, err
	}
	err = os.RemoveAll(filepath.Join(tmpDir, ".git"))
	newSrc := NewSource(tmpDir)
	newSrc.WithContext(src.Context())
	lh := &LocalHandler{}
	localFh, err := lh.Zip(newSrc)
	if err != nil {
		return nil, err
	}
	cleanFunc := func() error {
		err := localFh.Close()
		if err != nil {
			return err
		}
		return os.RemoveAll(tmpDir)
	}
	return NewZipFile(localFh, localFh.Size(), cleanFunc), nil
}
func (h GitHandler) makeGitUtils(tmpDir, path string) *GitUtils {
	u, err := url.Parse(path)
	if err != nil {
		u, _ = url.Parse("ssh://" + path)
	}
	refName := "master"
	if u.Fragment != "" {
		refName = u.Fragment
		u.Fragment = ""
	}
	var authMethod transport.AuthMethod
	if u.User != nil && IsWebURL(path) {
		password, _ := u.User.Password()
		authMethod = &githttp.BasicAuth{u.User.Username(), password}
		u.User = nil
	}
	gitUtils := &GitUtils{
		Url:        u.String(),
		Folder:     tmpDir,
		RefName:    refName,
		AuthMethod: authMethod,
	}
	return gitUtils
}
func (h GitHandler) Sha1(src *Source) (string, error) {
	h.setHttpClient(src)
	path := src.Path
	tmpDir, err := ioutil.TempDir("", "git-zipper")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmpDir)
	gitUtils := h.makeGitUtils(tmpDir, path)
	return gitUtils.CommitSha1()
}
func (h GitHandler) Detect(src *Source) bool {
	path := src.Path
	if !IsWebURL(path) {
		return false
	}
	u, err := url.Parse(path)
	if err != nil {
		return false
	}
	return HasExtFile(u.Path, ".git")
}
func (h *GitHandler) setHttpClient(src *Source) {
	*h.client = *CtxHttpClient(src)
}

func (h GitHandler) Name() string {
	return "git"
}

type GitUtils struct {
	Folder     string
	Url        string
	RefName    string
	AuthMethod transport.AuthMethod
}

var refTypes []string = []string{"heads", "tags"}

func (g GitUtils) Clone() error {
	_, err := g.findRepo(false)
	if err != nil {
		return err
	}
	return nil
}
func (g GitUtils) CommitSha1() (string, error) {
	if g.refNameIsHash() {
		return g.RefName, nil
	}
	repo, err := g.findRepo(true)
	if err != nil {
		return "", err
	}
	iter, err := repo.Log(&git.LogOptions{})
	if err != nil {
		return "", err
	}
	defer iter.Close()
	commit, err := iter.Next()
	if err != nil {
		return "", err
	}
	return commit.Hash.String(), nil
}
func (g GitUtils) refNameIsHash() bool {
	return len(g.RefName) == 40
}
func (g GitUtils) findRepoFromHash(isBare bool) (*git.Repository, error) {
	repo, err := git.PlainClone(g.Folder, isBare, &git.CloneOptions{
		URL:  g.Url,
		Auth: g.AuthMethod,
	})
	if err != nil {
		return nil, err
	}
	tree, err := repo.Worktree()
	if err != nil {
		return nil, err
	}
	err = tree.Checkout(&git.CheckoutOptions{
		Hash:  plumbing.NewHash(g.RefName),
		Force: true,
	})
	if err != nil {
		return nil, err
	}
	return repo, nil
}
func (g GitUtils) findRepo(isBare bool) (*git.Repository, error) {
	if g.refNameIsHash() {
		return g.findRepoFromHash(isBare)
	}
	var repo *git.Repository
	var err error
	for _, refType := range refTypes {
		repo, err = git.PlainClone(g.Folder, isBare, &git.CloneOptions{
			URL:          g.Url,
			SingleBranch: true,
			Auth:         g.AuthMethod,
			ReferenceName: plumbing.ReferenceName(fmt.Sprintf(
				"refs/%s/%s",
				refType,
				strings.ToLower(g.RefName),
			)),
			Depth: 1,
		})
		if err == nil {
			return repo, nil
		}
		if err.Error() == "reference not found" {
			os.RemoveAll(g.Folder)
			os.Mkdir(g.Folder, 0777)
			continue
		}
		return nil, err
	}
	return repo, err
}
