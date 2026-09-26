package articles

import "testing"

func TestResolveSlug(t *testing.T) {
	tests := []struct {
		name  string
		slug  string
		title string
		want  string
	}{
		{"falls back to title", "", "Hello  World!", "hello-world"},
		{"blank slug falls back", "   ", "Hello World", "hello-world"},
		{"explicit slug wins", "my-custom-slug", "Hello World", "my-custom-slug"},
		{"explicit slug is normalized", "My Custom Slug", "Hello World", "my-custom-slug"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := resolveSlug(tt.slug, tt.title); got != tt.want {
				t.Fatalf("resolveSlug(%q, %q) = %q, want %q", tt.slug, tt.title, got, tt.want)
			}
		})
	}
}
