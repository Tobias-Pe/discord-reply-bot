package logger

import "testing"

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name string
	}{
		{name: "init the logger"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if Logger != nil {
				t.Errorf("Logger shouldnt have been initialised yet")
			}
			InitLogger()
			if Logger == nil {
				t.Errorf("Logger should have been initialised")
			}
		})
	}
}
