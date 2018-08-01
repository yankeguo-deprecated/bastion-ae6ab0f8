package utils

import "testing"

var matchAsteriskCases = []struct {
	P string
	V string
	R bool
}{
	{P: "hello.*", V: "hello.world", R: true},
	{P: "你好*", V: "你好世界", R: true},
	{P: "你好*", V: "好世界", R: false},
	{P: "你*世界", V: "你好世界", R: true},
}

func TestMatchAsterisk(t *testing.T) {
	for _, c := range matchAsteriskCases {
		if MatchAsterisk(c.P, c.V) != c.R {
			t.Errorf("failed %s vs %s should be %v", c.P, c.V, c.R)
		}
	}
}
