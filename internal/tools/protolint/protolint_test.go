package protolint_test

import (
	"slices"
	"testing"

	"buf.build/go/bufplugin/check"
	"buf.build/go/bufplugin/check/checktest"

	"github.com/fun-dotto/server/internal/tools/protolint"
)

func TestSpec(t *testing.T) {
	t.Parallel()
	checktest.SpecTest(t, protolint.Spec)
}

func TestGood(t *testing.T) {
	t.Parallel()
	checktest.CheckTest{
		Request: &checktest.RequestSpec{
			Files: &checktest.ProtoFileSpec{
				DirPaths:  []string{"testdata/good"},
				FilePaths: []string{"admin/v1/subject.proto"},
			},
		},
		Spec: protolint.Spec,
	}.Run(t)
}

func TestBad(t *testing.T) {
	t.Parallel()
	request, err := (&checktest.RequestSpec{
		Files: &checktest.ProtoFileSpec{
			DirPaths:  []string{"testdata/bad"},
			FilePaths: []string{"admin/v1/widget.proto", "admin/v1/mismatch.proto"},
		},
	}).ToRequest(t.Context())
	if err != nil {
		t.Fatal(err)
	}
	client, err := check.NewClientForSpec(protolint.Spec)
	if err != nil {
		t.Fatal(err)
	}
	response, err := client.Check(t.Context(), request)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, a := range response.Annotations() {
		got = append(got, a.RuleID())
		t.Logf("%s:%d: %s: %s", a.FileLocation().FileDescriptor().ProtoreflectFileDescriptor().Path(), a.FileLocation().StartLine()+1, a.RuleID(), a.Message())
	}
	// すべてのルールが少なくとも 1 回は違反を検出すること
	for _, rule := range protolint.Spec.Rules {
		if !slices.Contains(got, rule.ID) {
			t.Errorf("rule %s reported no annotation", rule.ID)
		}
	}
}
