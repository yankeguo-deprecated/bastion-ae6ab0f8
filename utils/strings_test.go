package utils

import "testing"

func TestExtractReadableString(t *testing.T) {
	if ExtractReadableString([]byte("\x03\x1cHello 你好\x02\x9f")) != "Hello 你好" {
		t.Fatal("failed")
	}
}
