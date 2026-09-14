package strictjson

import (
	"strings"
	"testing"
)

func TestRejectDuplicateKeys(t *testing.T) {
	t.Parallel()
	for _, input := range []string{
		`{"a":1,"a":2}`,
		`{"outer":{"x":1,"x":2}}`,
		`[{"x":1,"x":2}]`,
	} {
		if err := RejectDuplicateKeys(strings.NewReader(input)); err == nil || !strings.Contains(err.Error(), "duplicate JSON field") {
			t.Fatalf("RejectDuplicateKeys(%s) error = %v, want duplicate rejection", input, err)
		}
	}
}

func TestRejectDuplicateKeysAllowsSameNameInSeparateObjects(t *testing.T) {
	t.Parallel()
	input := `[{"x":1},{"x":2}]`
	if err := RejectDuplicateKeys(strings.NewReader(input)); err != nil {
		t.Fatalf("RejectDuplicateKeys(%s) error = %v", input, err)
	}
}

func FuzzRejectDuplicateKeys(f *testing.F) {
	f.Add([]byte(`{"a":1}`))
	f.Add([]byte(`{"a":1,"a":2}`))
	f.Add([]byte(`[{"a":{"b":1}}]`))
	f.Fuzz(func(t *testing.T, data []byte) {
		_ = RejectDuplicateKeys(strings.NewReader(string(data)))
	})
}
