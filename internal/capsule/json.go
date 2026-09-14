package capsule

import (
	"bytes"

	"arkmesh/internal/strictjson"
)

func rejectDuplicateJSONKeys(data []byte) error {
	return strictjson.RejectDuplicateKeys(bytes.NewReader(data))
}
