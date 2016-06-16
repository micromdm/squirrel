package datastore

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/groob/ape/models"
)

// GitRepo lets you push back changes to a repository
type GitRepo struct {
	URL     string // Repository URL
	Path    string // Directory to pull to
	Host    string // Git domain host e.g. github.com
	Branch  string // Git branch
	KeyPath string // Path to private ssh key
	Datastore
}

// NewManifest creates a new manifest and commits it to git.
func (r *GitRepo) NewManifest(name string) (*models.Manifest, error) {
	manifest, err := r.Datastore.NewManifest(name)
	if err != nil {
		return nil, err
	}
	manifestPath := fmt.Sprintf("manifests/%v", name)
	commitMsg := fmt.Sprintf("Created manifest %v", name)
	err = r.addCommitPush(manifestPath, commitMsg)
	return manifest, err
}

func (r *GitRepo) addCommitPush(path string, msg ...string) error {
	// add
	cmd := exec.Command("git", "add", path)
	cmd.Dir = r.Path
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running command %q: %v: %s", strings.Join(cmd.Args, " "), err, string(output))
	}

	// commit
	// combine commit messages each on a separate line
	commitMsg := strings.Join(msg, "\n")
	cmd = exec.Command("git", "commit", "-m", commitMsg)
	cmd.Dir = r.Path
	output, err = cmd.CombinedOutput()
	if err != nil {
		log.Printf("error running command %q: %v: %s", strings.Join(cmd.Args, " "), err, string(output))
	}
	// push
	cmd = exec.Command("git", "push")
	cmd.Dir = r.Path
	output, err = cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("error running command %q: %v: %s", strings.Join(cmd.Args, " "), err, string(output))
	}
	return nil
}
