package v1

import (
	"strings"
	"testing"
)

func TestMdToHTMLSanitizesUnsafeHTML(t *testing.T) {
	for _, tc := range []struct {
		name  string
		input string
		deny  []string
	}{
		{name: "script", input: "<script>alert(1)</script>", deny: []string{"<script", "</script>"}},
		{name: "event handler", input: `<img src="x" onerror="alert(1)">`, deny: []string{"onerror"}},
		{name: "javascript link", input: `[click](javascript:alert(1))`, deny: []string{"javascript:"}},
		{name: "iframe", input: `<iframe src="//evil.example"></iframe>`, deny: []string{"<iframe", "</iframe>"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := string(mdToHTML([]byte(tc.input)))
			for _, forbidden := range tc.deny {
				if strings.Contains(strings.ToLower(got), strings.ToLower(forbidden)) {
					t.Fatalf("sanitized HTML contains %q: %s", forbidden, got)
				}
			}
		})
	}
}

func TestMdToHTMLPreservesCommonMarkdown(t *testing.T) {
	got := string(mdToHTML([]byte("# Heading\n\n**bold** and `code`")))
	for _, want := range []string{"<h1", "Heading", "<strong>bold</strong>", "<code>code</code>"} {
		if !strings.Contains(got, want) {
			t.Fatalf("sanitized HTML missing %q: %s", want, got)
		}
	}
}
