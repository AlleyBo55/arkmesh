package capsule

import (
	"encoding/json"
	"testing"
)

func FuzzCapsuleEnvelopeDecoders(f *testing.F) {
	seeds := [][]byte{
		[]byte(`{}`),
		[]byte(`null`),
		[]byte(`{"schema_version":"arkmesh.signature/v0alpha1"}`),
		[]byte(`{"schema_version":"arkmesh.authority/v0alpha1","action":"rotate"}`),
		[]byte(`{"schema_version":"arkmesh.checkpoint/v0alpha1"}`),
		[]byte(`{"schema_version":"arkmesh.recovery-policy/v0alpha1","threshold":2,"members":[]}`),
		[]byte("{\"x\":1}\n{\"y\":2}"),
	}
	for _, seed := range seeds {
		f.Add(seed)
	}
	f.Fuzz(func(t *testing.T, data []byte) {
		var signature SignatureEnvelope
		_ = decodeSignature(data, &signature)
		var authority AuthorityTransition
		_ = decodeAuthority(data, &authority)
		var checkpoint Checkpoint
		_ = decodeCheckpoint(data, &checkpoint)

		var manifest Manifest
		if json.Unmarshal(data, &manifest) == nil {
			_ = validateManifest(manifest)
			_ = capsuleID(manifest)
		}
		var policy RecoveryPolicy
		if json.Unmarshal(data, &policy) == nil {
			_ = policy.Validate()
			_ = recoveryPolicyID(policy)
		}
		var recovery RecoveryTransition
		if json.Unmarshal(data, &recovery) == nil {
			_ = recoveryPayload(recovery.PolicyID, recovery.Action, recovery.ParentCapsuleID, recovery.ChildCapsuleID, recovery.PreviousSignerID, recovery.NextSigner.KeyID)
			_, _ = recovery.NextSigner.Key()
			for _, approval := range recovery.Approvals {
				_, _ = decodeRecoverySignature(approval.Signature)
			}
		}
		var proof ChunkProof
		if json.Unmarshal(data, &proof) == nil {
			_ = VerifyChunkProof(manifest, proof)
			_ = verifyChunkPath(proof.ChunkRoot, proof.ChunkCount, proof.ChunkIndex, chunkLeafDigest(data), proof.Path)
		}
		var tree ChunkTreeFile
		if json.Unmarshal(data, &tree) == nil {
			_, _ = VerifyChunkTreeFile(manifest, tree)
		}
	})
}
