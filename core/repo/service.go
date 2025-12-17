package repo

import (
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
)

var (
	MkdirPattern = "shipit-repo-"
)

type CloneService struct {
	RepoUrl *url.URL
	Host    string
	User    string
	TempDir string
	AbsPath string

	Out    io.Writer
	ErrOut io.Writer
}

type CopyContext struct {
	TempDir    string
	AbsPath    string
	RemoteDir  string
	TargetHost string
	ScpTarget  string
}

func NewCloneService(host, user string, out, errOut io.Writer) *CloneService {
	return &CloneService{Host: host, User: user, Out: out, ErrOut: errOut}
}

func (s *CloneService) Run() error {

	if err := s.validate(); err != nil {
		return err
	}

	ctx, err := s.buildCopyContext()
	if err != nil {
		return err
	}

	if err := s.cloneRepo(ctx); err != nil {
		return err
	}

	if err := s.copyRepo(ctx); err != nil {
		return err
	}
	return nil

}

func (s *CloneService) validate() error {

	if err := isValidHost(s.Host); err != nil {
		return fmt.Errorf("repo: %w: %w", ErrisValidHost, err)
	}

	if err := isValidUser(s.User); err != nil {
		return fmt.Errorf("repo: %w: %w", ErrisValidUser, err)
	}

	if err := isValidUrl(s.RepoUrl); err != nil {
		return fmt.Errorf("repo: %w: %w", ErrisValidUrl, err)
	}

	return nil
}

func (s *CloneService) buildCopyContext() (*CopyContext, error) {

	dir, err := os.MkdirTemp("", MkdirPattern)
	if err != nil {
		return nil, fmt.Errorf("repo: %w: %w", ErrMakeTempDir, err)
	}

	s.TempDir = dir

	path, err := filepath.Abs(s.TempDir)
	if err != nil {
		return nil, fmt.Errorf("repo: %w: %w", ErrMakeTempDir, err)
	}

	s.AbsPath = path

	remoteDir := fmt.Sprintf("~/%s", filepath.Base(s.AbsPath))
	targetHost := fmt.Sprintf("%s@%s", s.User, s.Host)

	scpTarget := fmt.Sprintf("%s:~", targetHost)

	return &CopyContext{
		TempDir:    dir,
		AbsPath:    path,
		RemoteDir:  remoteDir,
		TargetHost: targetHost,
		ScpTarget:  scpTarget,
	}, nil
}

func (s *CloneService) ParseRepoUrl(repoCloneURL string) error {

	url, err := url.Parse(repoCloneURL)
	if err != nil {
		return fmt.Errorf("repo: %w: %w", ErrParseUrl, err)
	}
	s.RepoUrl = url
	return nil

}

func (s *CloneService) cloneRepo(ctx *CopyContext) error {

	gitCmd := exec.Command("git", "clone", s.RepoUrl.String(), ctx.TempDir)
	gitCmd.Stdout = s.Out
	gitCmd.Stderr = s.ErrOut

	if err := os.Remove(s.TempDir); err != nil {
		return fmt.Errorf("repo: %w: %w", ErrCloneRepo, err)
	}

	if err := gitCmd.Run(); err != nil {
		return fmt.Errorf("repo %w: %w", ErrCloneRepo, err)
	}

	return nil
}

func (s *CloneService) copyRepo(ctx *CopyContext) error {

	scpCmd := exec.Command("scp", "-r", ctx.AbsPath, ctx.ScpTarget)
	scpCmd.Stdout = s.Out
	scpCmd.Stderr = s.ErrOut

	fmt.Fprintf(s.Out, "\nCopying repository to %s:%s\n", ctx.TargetHost, ctx.RemoteDir)

	if err := scpCmd.Run(); err != nil {
		return fmt.Errorf("repo: %w: %w", ErrCopyRepo, err)
	}

	fmt.Fprintf(s.Out, "\nRepository copied successfully\n")
	fmt.Fprintf(s.Out, "\nRun:\n  cd %s\n", s.AbsPath)

	return nil
}

func isValidHost(Host string) error {
	if Host == "" {
		return fmt.Errorf("empty string")
	}
	return nil
}

func isValidUser(User string) error {
	if User == "" {
		return fmt.Errorf("empty string")
	}
	return nil
}

func isValidUrl(Url *url.URL) error {
	if Url.String() == "" {
		return fmt.Errorf("empty string")
	}
	return nil
}
