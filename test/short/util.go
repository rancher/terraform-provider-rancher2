package short

import (
	"os"
	"path/filepath"
	"testing"

	g "github.com/gruntwork-io/terratest/modules/git"
	"github.com/gruntwork-io/terratest/modules/random"
)

// Helpers.
func repoRoot(t *testing.T) (string, error) {
	gwd := g.GetRepoRoot(t)
	fwd, err := filepath.Abs(gwd)
	if err != nil {
		return "", err
	}
	return fwd, nil
}

func createTestDirectories(t *testing.T, id string) (testIdPath string, err error) {
	fwd, err := repoRoot(t)
	if err != nil {
		return "", err
	}
	paths := []string{
		filepath.Join(fwd, "test", "short", "data"),
		filepath.Join(fwd, "test", "short", "data", id),
		filepath.Join(fwd, "test", "short", "data", id, "data"),
	}
	for _, path := range paths {
		err = os.Mkdir(path, 0755)
		if err != nil && !os.IsExist(err) {
			return "", err
		}
	}
	return filepath.Join(fwd, "test", "short", "data", id), nil
}

func id() string {
	id := os.Getenv("IDENTIFIER")
	if id == "" {
		id = random.UniqueId()
	}
	id += "-" + random.UniqueId()
	return id
}
