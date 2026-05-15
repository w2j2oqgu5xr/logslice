package parser

import "testing"

func TestParseLevel_KnownLevels(t *testing.T) {
	cases := []struct {
		input    string
		want     Level
	}{
		{"debug", LevelDebug},
		{"DEBUG", LevelDebug},
		{"info", LevelInfo},
		{"INFO", LevelInfo},
		{"warn", LevelWarn},
		{"WARN", LevelWarn},
		{"warning", LevelWarn},
		{"WARNING", LevelWarn},
		{"error", LevelError},
		{"ERROR", LevelError},
		{"fatal", LevelFatal},
		{"FATAL", LevelFatal},
	}

	for _, tc := range cases {
		t.Run(tc.input, func(t *testing.T) {
			got := ParseLevel(tc.input)
			if got != tc.want {
				t.Errorf("ParseLevel(%q) = %v, want %v", tc.input, got, tc.want)
			}
		})
	}
}

func TestParseLevel_Unknown(t *testing.T) {
	cases := []string{"", "trace", "verbose", "notice", "bogus"}
	for _, input := range cases {
		t.Run(input, func(t *testing.T) {
			got := ParseLevel(input)
			if got != LevelUnknown {
				t.Errorf("ParseLevel(%q) = %v, want LevelUnknown", input, got)
			}
		})
	}
}

func TestLevel_String(t *testing.T) {
	cases := []struct {
		level Level
		want  string
	}{
		{LevelDebug, "DEBUG"},
		{LevelInfo, "INFO"},
		{LevelWarn, "WARN"},
		{LevelError, "ERROR"},
		{LevelFatal, "FATAL"},
		{LevelUnknown, "UNKNOWN"},
	}

	for _, tc := range cases {
		t.Run(tc.want, func(t *testing.T) {
			got := tc.level.String()
			if got != tc.want {
				t.Errorf("Level.String() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestLevel_IsValid(t *testing.T) {
	valid := []Level{LevelDebug, LevelInfo, LevelWarn, LevelError, LevelFatal}
	for _, l := range valid {
		if !l.IsValid() {
			t.Errorf("expected %v to be valid", l)
		}
	}

	if LevelUnknown.IsValid() {
		t.Error("expected LevelUnknown to be invalid")
	}
}
