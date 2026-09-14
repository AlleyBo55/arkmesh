package identity

import "testing"

func FuzzIdentityDecoders(f *testing.F) {
	seeds := [][]byte{
		[]byte(`{}`),
		[]byte(`null`),
		[]byte(`{"schema_version":"arkmesh.identity/v0alpha1","algorithm":"ed25519"}`),
		[]byte(`{"schema_version":"arkmesh.private-identity/v0alpha1","algorithm":"ed25519"}`),
		[]byte("{\"x\":1}\n{\"y\":2}"),
	}
	for _, seed := range seeds {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		var public PublicIdentity
		if decodeStrict(data, &public) == nil {
			_, _ = public.Key()
		}
		var private PrivateIdentity
		if decodeStrict(data, &private) == nil {
			_, _, _ = private.Keys()
		}
	})
}
