package service

import "testing"

func TestConvert(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{"text sample", "Привет", ".--. .-. .. .-- . -", false},
		{"Morse sample", ".--. .-. .. .-- . -", "ПРИВЕТ", false},
		{"text words", "Привет мир", ".--. .-. .. .-- . -   -- .. .-.", false},
		{"Morse words", ".--. .-. .. .-- . -   -- .. .-.", "ПРИВЕТ МИР", false},
		{"text line breaks", "  Привет\r\nмир\t ", ".--. .-. .. .-- . -   -- .. .-.", false},
		{"Morse line breaks", ".--. .-. .. .-- . -\r\n-- .. .-.\r\n", "ПРИВЕТ МИР", false},
		{"Morse spacing", ".--.  .-.\t.. .-- . -      -- .. .-.", "ПРИВЕТ МИР", false},
		{"digits and punctuation", "1, 2?", ".---- .-.-.-   ..--- ..--..", false},
		{"single dot", ".", "Е", false},
		{"empty", "", "", true},
		{"whitespace", " \r\n\t", "", true},
		{"invalid Morse", ".......", "", true},
		{"partly invalid Morse", ".- .......", "", true},
		{"unsupported text", "Hello", "", true},
		{"partly unsupported text", "Привет🙂", "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Convert(tt.input)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Convert() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Convert() = %q, want %q", got, tt.want)
			}
		})
	}
}
