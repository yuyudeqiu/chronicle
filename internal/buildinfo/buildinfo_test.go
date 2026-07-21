package buildinfo

import "testing"

func TestResolvePseudoVersion(t *testing.T) {
	got := resolve(Info{}, "v0.0.0-20260721033902-a22e69fc0898", nil)

	if got.Version != "dev" {
		t.Fatalf("Version = %q, want dev", got.Version)
	}
	if got.Commit != "a22e69fc0898" {
		t.Fatalf("Commit = %q, want a22e69fc0898", got.Commit)
	}
	if got.CommitDate != "2026-07-21" {
		t.Fatalf("CommitDate = %q, want 2026-07-21", got.CommitDate)
	}
	if got.ModuleVersion != "v0.0.0-20260721033902-a22e69fc0898" {
		t.Fatalf("ModuleVersion = %q", got.ModuleVersion)
	}
}

func TestResolveDirtyPseudoVersion(t *testing.T) {
	got := resolve(Info{}, "v0.0.0-20260721033902-a22e69fc0898+dirty", nil)
	if got.Version != "dev" || got.Commit != "a22e69fc0898" || !got.Modified {
		t.Fatalf("resolve dirty pseudo-version = %#v", got)
	}
}

func TestResolveTaggedVersion(t *testing.T) {
	got := resolve(Info{}, "v0.1.0", nil)
	if got.Version != "v0.1.0" {
		t.Fatalf("Version = %q, want v0.1.0", got.Version)
	}
	if got.String() != "chronicle v0.1.0" {
		t.Fatalf("String() = %q", got.String())
	}
}

func TestInjectedVersionWinsAndIsNormalized(t *testing.T) {
	got := resolve(Info{Version: "0.2.0", Commit: "0123456789abcdef"}, "(devel)", map[string]string{
		"vcs.modified": "true",
	})
	if got.Version != "v0.2.0" {
		t.Fatalf("Version = %q, want v0.2.0", got.Version)
	}
	if got.String() != "chronicle v0.2.0 (commit 0123456789abcdef, modified)" {
		t.Fatalf("String() = %q", got.String())
	}
}

func TestCommitOnlyBuildIsDevelopmentVersion(t *testing.T) {
	got := resolve(Info{Version: "a22e69f-dirty", Commit: "a22e69f"}, "(devel)", nil)
	if got.Version != "dev" {
		t.Fatalf("Version = %q, want dev", got.Version)
	}
}
