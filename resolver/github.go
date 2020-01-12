package resolver

import (
	"fmt"
	"io"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

// github implementation.
type github struct{
	resolver Interface
}

// NewGithubResolver returns a github resolver.
func NewGithubResolver(resolver Interface) Interface {
	return &github{resolver:resolver}
}

// Get implementation.
func (r *github) Get(s string) (io.ReadCloser, error) {
	u, err := url.Parse(fmt.Sprintf("https://%s", s))
	if err != nil {
		return nil, errors.Wrap(err, "parsing include path")
	}

	if u.Host != "github.com" {
		return nil, ErrNotSupported
	}

	parts := strings.SplitN(u.Path, "/", 4)
	if len(parts) < 3 {
		return nil, errors.New("user, repo required in include url")
	}

	if len(parts) < 4 {
		parts = append(parts, "index.mk")
	}

	user := parts[1]
	repo := parts[2]
	file := parts[3]
	raw := fmt.Sprintf("https://raw.githubusercontent.com/%s/%s/master/%s", user, repo, file)

	return r.resolver.Get(raw)
}
