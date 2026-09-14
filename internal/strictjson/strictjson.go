package strictjson

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type frame struct {
	object       bool
	expectingKey bool
	keys         map[string]struct{}
}

func RejectDuplicateKeys(reader io.Reader) error {
	decoder := json.NewDecoder(reader)
	decoder.UseNumber()
	var stack []frame
	for {
		token, err := decoder.Token()
		if errors.Is(err, io.EOF) {
			return nil
		}
		if err != nil {
			return err
		}
		if delimiter, ok := token.(json.Delim); ok {
			switch delimiter {
			case '{':
				stack = append(stack, frame{object: true, expectingKey: true, keys: map[string]struct{}{}})
			case '[':
				stack = append(stack, frame{})
			case '}', ']':
				if len(stack) == 0 {
					return errors.New("unexpected JSON delimiter")
				}
				stack = stack[:len(stack)-1]
				completeValue(stack)
			}
			continue
		}
		if len(stack) > 0 && stack[len(stack)-1].object && stack[len(stack)-1].expectingKey {
			key, ok := token.(string)
			if !ok {
				return errors.New("JSON object key is not a string")
			}
			current := &stack[len(stack)-1]
			if _, exists := current.keys[key]; exists {
				return fmt.Errorf("duplicate JSON field %q", key)
			}
			current.keys[key] = struct{}{}
			current.expectingKey = false
			continue
		}
		completeValue(stack)
	}
}

func completeValue(stack []frame) {
	if len(stack) > 0 && stack[len(stack)-1].object {
		stack[len(stack)-1].expectingKey = true
	}
}
