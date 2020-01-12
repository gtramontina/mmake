package resolver_test

import (
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/pkg/errors"
	"github.com/tj/mmake/resolver"
)

type fakeResolver struct {
	contents map[string]string
}

func (r fakeResolver) Get(url string) (io.ReadCloser, error) {
	content, ok := r.contents[url]
	if !ok {
		return nil, errors.Errorf("Unexpected URL: %s", url)
	}
	return ioutil.NopCloser(strings.NewReader(content)), nil
}

func TestGithubResolver_Errors(t *testing.T) {
	res := resolver.NewGithubResolver(fakeResolver{})

	var cases = []struct {
		Path    string
		ErrorMsg string
	}{
		{"not url", "parsing include path"},
		{"bitbucket.org", resolver.ErrNotSupported.Error()},
		{"github.com", "user, repo required"},
		{"github.com/user", "user, repo required"},
	}

	for _, c := range cases {
		t.Run(c.Path, func(t *testing.T) {
			r, err := res.Get(c.Path)
			if r != nil {
				t.Errorf("expected nil, got %v", r)
			}

			if err == nil {
				t.Error("expected nil")
			}

			if !strings.HasPrefix(err.Error(), c.ErrorMsg) {
				t.Errorf("expected to start with '%v', got '%v'", c.ErrorMsg, err)
			}
		})
	}
}

func TestGithubResolver(t *testing.T) {
	res := resolver.NewGithubResolver(fakeResolver{
		contents: map[string]string {
			"https://raw.githubusercontent.com/user/repo/master/index.mk": "index.mk content",
			"https://raw.githubusercontent.com/user/repo/master/bar": "bar content",
			"https://raw.githubusercontent.com/user/repo/master/foo.mk": "foo.mk content",
			"https://raw.githubusercontent.com/user/repo/master/baz/stuff.mk": "baz/stuff.mk content",
		},
	})

	var cases = []struct {
		Path    string
		Content string
	}{
		{"github.com/user/repo", "index.mk content"},
		{"github.com/user/repo/bar", "bar content"},
		{"github.com/user/repo/foo.mk", "foo.mk content"},
		{"github.com/user/repo/baz/stuff.mk", "baz/stuff.mk content"},
	}

	for _, c := range cases {
		t.Run(c.Path, func(t *testing.T) {
			r, err := res.Get(c.Path)
			if err != nil {
				t.Fatal(err)
			}
			defer r.Close()

			b, err := ioutil.ReadAll(r)
			if err != nil {
				t.Fatal(err)
			}

			if string(b) != c.Content {
				t.Errorf("expected %q, got %q", c.Content, string(b))
			}
		})
	}
}
