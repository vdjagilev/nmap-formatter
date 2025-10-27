package formatter

import (
	"bytes"
	"encoding/json"
	"io"
	"unicode"
)

// snakeCaseEncoder wraps a standard JSON encoder and converts keys to snake_case
type snakeCaseEncoder struct {
	writer      io.Writer
	indent      string
	prettyPrint bool
}

// newSnakeCaseEncoder creates a new encoder that converts JSON keys to snake_case
func newSnakeCaseEncoder(w io.Writer, prettyPrint bool) *snakeCaseEncoder {
	indent := ""
	if prettyPrint {
		indent = "  "
	}
	return &snakeCaseEncoder{
		writer:      w,
		indent:      indent,
		prettyPrint: prettyPrint,
	}
}

// Encode encodes the value to JSON with snake_case keys
func (e *snakeCaseEncoder) Encode(v interface{}) error {
	// First, encode normally to a buffer
	buf := new(bytes.Buffer)
	encoder := json.NewEncoder(buf)
	if e.prettyPrint {
		encoder.SetIndent("", e.indent)
	}
	if err := encoder.Encode(v); err != nil {
		return err
	}

	// Convert keys to snake_case
	converted := convertKeysToSnakeCase(buf.Bytes())

	// Write to the actual writer
	_, err := e.writer.Write(converted)
	return err
}

// toSnakeCase converts a CamelCase string to snake_case
// This is optimized for performance - it processes the string in a single pass
func toSnakeCase(s string) string {
	if s == "" {
		return s
	}

	// Pre-allocate buffer with estimated size (original length + 30% for underscores)
	var buf bytes.Buffer
	buf.Grow(len(s) + len(s)/3)

	var prevLower bool
	var prevUnderscore bool

	for i, r := range s {
		if unicode.IsUpper(r) {
			// Add underscore before uppercase if:
			// 1. Not the first character
			// 2. Previous character was lowercase
			// 3. Previous character wasn't already an underscore
			if i > 0 && prevLower && !prevUnderscore {
				buf.WriteByte('_')
			}
			buf.WriteRune(unicode.ToLower(r))
			prevLower = false
			prevUnderscore = false
		} else if r == '_' {
			buf.WriteRune(r)
			prevLower = false
			prevUnderscore = true
		} else {
			buf.WriteRune(r)
			prevLower = true
			prevUnderscore = false
		}
	}

	return buf.String()
}

// convertKeysToSnakeCase converts all JSON keys in the byte slice to snake_case
// This processes the JSON in a streaming fashion for better performance
func convertKeysToSnakeCase(data []byte) []byte {
	result := make([]byte, 0, len(data))
	inString := false
	escaped := false
	keyStart := 0

	for i := 0; i < len(data); i++ {
		b := data[i]

		if escaped {
			result = append(result, b)
			escaped = false
			continue
		}

		switch b {
		case '\\':
			if inString {
				escaped = true
			}
			result = append(result, b)

		case '"':
			if !inString {
				// Starting a potential key
				inString = true
				result = append(result, b)
				keyStart = len(result)
			} else {
				// Ending a string
				inString = false

				// Check if this was a key (followed by colon after whitespace)
				isKey := false
				j := i + 1
				for j < len(data) && (data[j] == ' ' || data[j] == '\t' || data[j] == '\n' || data[j] == '\r') {
					j++
				}
				if j < len(data) && data[j] == ':' {
					isKey = true
				}

				if isKey {
					// Extract the key and convert it
					key := result[keyStart:]
					snakeKey := toSnakeCase(string(key))
					result = result[:keyStart]
					result = append(result, []byte(snakeKey)...)
				}

				result = append(result, b)
			}

		default:
			result = append(result, b)
		}
	}

	return result
}
