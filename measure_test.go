package gs1

import "testing"

func TestMeasure(t *testing.T) {
	tests := []struct {
		name     string
		prefix   string
		ai       string
		wantUnit string
	}{
		{name: "net kilograms", prefix: "310", ai: "3102", wantUnit: "kg"},
		{name: "net pounds", prefix: "320", ai: "3203", wantUnit: "lb"},
		{name: "gross kilograms", prefix: "330", ai: "3304", wantUnit: "kg"},
		{name: "gross pounds", prefix: "340", ai: "3405", wantUnit: "lb"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b, err := Parse(tt.ai + "001250")
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			got, ok := b.Measure(tt.prefix)
			if !ok {
				t.Fatalf("Measure(%q) = not found", tt.prefix)
			}
			if got.Raw != "001250" || got.Scaled != 1250 {
				t.Errorf("raw/scaled = %q/%d, want 001250/1250", got.Raw, got.Scaled)
			}
			if got.Decimals != int(tt.ai[3]-'0') {
				t.Errorf("Decimals = %d, want %d", got.Decimals, tt.ai[3]-'0')
			}
			if got.Unit != tt.wantUnit {
				t.Errorf("Unit = %q, want %q", got.Unit, tt.wantUnit)
			}
			want := 1250.0
			for i := 0; i < got.Decimals; i++ {
				want /= 10
			}
			if got.Value != want {
				t.Errorf("Value = %v, want %v", got.Value, want)
			}
		})
	}
}

func TestMeasureDecimalPositions(t *testing.T) {
	for decimals := 0; decimals <= 5; decimals++ {
		t.Run(itoa(decimals), func(t *testing.T) {
			b, err := Parse("310" + itoa(decimals) + "001250")
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}
			got, ok := b.NetWeightKg()
			if !ok {
				t.Fatal("NetWeightKg() = not found")
			}
			if got.Decimals != decimals {
				t.Errorf("Decimals = %d, want %d", got.Decimals, decimals)
			}
		})
	}
}

func TestMeasureMissingOrUnsupported(t *testing.T) {
	b, err := Parse("0104150000021126")
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if _, ok := b.Measure("310"); ok {
		t.Error("Measure(310) = found, want missing")
	}
	if _, ok := b.Measure("311"); ok {
		t.Error("Measure(311) = found, want unsupported")
	}
}
