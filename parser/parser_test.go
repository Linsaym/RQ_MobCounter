package parser

import (
	"os"
	"testing"
)

func TestParseLogEntry(t *testing.T) {
	tests := []struct {
		name       string
		timestamp  string
		content    string
		wantName   string
		wantExp    int
		wantNil    bool
	}{
		{
			name:      "Simple entry with exp",
			timestamp: "1/16 06:45:41",
			content:   "Злая шкатулка погибает. Получено опыта: 2873.",
			wantName:  "Злая шкатулка",
			wantExp:   2873,
			wantNil:   false,
		},
		{
			name:      "Entry without exp",
			timestamp: "1/16 06:51:17",
			content:   "Росинка погибает.",
			wantName:  "Росинка",
			wantExp:   0,
			wantNil:   false,
		},
		{
			name:      "Entry with spaces in monster name",
			timestamp: "1/16 06:51:26",
			content:   "Луговая Жужа погибает.",
			wantName:  "Луговая Жужа",
			wantExp:   0,
			wantNil:   false,
		},
		{
			name:      "System message without monster death",
			timestamp: "1/16 06:58:30",
			content:   "Вы достигли 2 уровня!",
			wantNil:   true,
		},
		{
			name:      "Empty content",
			timestamp: "1/16 06:45:41",
			content:   "",
			wantNil:   true,
		},
		{
			name:      "Entry with HTML tags",
			timestamp: "1/16 06:45:41",
			content:   "<b>Часы</b> погибает. Получено опыта: 17530.",
			wantName:  "Часы",
			wantExp:   17530,
			wantNil:   false,
		},
		{
			name:      "Small exp value",
			timestamp: "1/16 06:58:01",
			content:   "Росинка погибает. Получено опыта: 3.",
			wantName:  "Росинка",
			wantExp:   3,
			wantNil:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseLogEntry(tt.timestamp, tt.content)

			if tt.wantNil {
				if result != nil {
					t.Errorf("Expected nil, got %v", result)
				}
				return
			}

			if result == nil {
				t.Errorf("Expected non-nil result")
				return
			}

			if result.MonsterName != tt.wantName {
				t.Errorf("Monster name: got %q, want %q", result.MonsterName, tt.wantName)
			}

			if result.ExpGained != tt.wantExp {
				t.Errorf("Exp gained: got %d, want %d", result.ExpGained, tt.wantExp)
			}

			if result.Timestamp != tt.timestamp {
				t.Errorf("Timestamp: got %q, want %q", result.Timestamp, tt.timestamp)
			}
		})
	}
}

func TestParseFile(t *testing.T) {
	// Создаем временный HTML файл для тестирования
	tmpFile := "test_parse_file.htm"

	htmlContent := `<HTML>
<HEAD><meta http-equiv="content-type" content="text/html; charset=UTF-8" /><TITLE>exp</TITLE></HEAD>
<BODY>
<TABLE width=800 align=center bgcolor=#333333 style='white-space:pre-wrap'><TR style='color:#4A92D3' valign=top title='1/16 06:45:41'><TD colspan=2>Злая шкатулка погибает. Получено опыта: 2873.
<TR style='color:#4A92D3' valign=top title='1/16 06:45:53'><TD colspan=2>Часы погибает. Получено опыта: 17530.
<TR style='color:#4A92D3' valign=top title='1/16 06:51:17'><TD colspan=2>Росинка погибает.
<TR style='color:#4A92D3' valign=top title='1/16 06:51:26'><TD colspan=2>Луговая Жужа погибает.
<TR style='color:#4A92D3' valign=top title='1/16 06:58:30'><TD colspan=2>Вы достигли 2 уровня!
</TABLE>
</BODY>
</HTML>`

	// Пишем контент в файл
	err := writeTempFile(tmpFile, htmlContent)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer deleteTempFile(tmpFile)

	// Парсим файл
	entries, err := ParseFile(tmpFile)
	if err != nil {
		t.Fatalf("ParseFile failed: %v", err)
	}

	// Проверяем результаты
	if len(entries) != 4 {
		t.Errorf("Expected 4 entries, got %d", len(entries))
	}

	// Проверяем первую запись
	if entries[0].MonsterName != "Злая шкатулка" {
		t.Errorf("First entry name: got %q, want %q", entries[0].MonsterName, "Злая шкатулка")
	}
	if entries[0].ExpGained != 2873 {
		t.Errorf("First entry exp: got %d, want %d", entries[0].ExpGained, 2873)
	}

	// Проверяем запись без опыта
	for _, entry := range entries {
		if entry.MonsterName == "Росинка" && entry.ExpGained != 0 {
			t.Errorf("Росинка should have 0 exp, got %d", entry.ExpGained)
		}
	}
}

func TestMonsterNameExtraction(t *testing.T) {
	tests := []struct {
		content  string
		expected string
	}{
		{"Злая шкатулка погибает. Получено опыта: 2873.", "Злая шкатулка"},
		{"Луговая Жужа погибает.", "Луговая Жужа"},
		{"  Часы  погибает. Получено опыта: 17530.", "Часы"},
		{"Стрекоскоп погибает. Получено опыта: 4070.", "Стрекоскоп"},
	}

	for _, tt := range tests {
		result := parseLogEntry("test", tt.content)
		if result == nil {
			t.Errorf("Failed to parse: %s", tt.content)
			continue
		}
		if result.MonsterName != tt.expected {
			t.Errorf("Content: %s\nGot: %q, Expected: %q", tt.content, result.MonsterName, tt.expected)
		}
	}
}

func TestExpExtraction(t *testing.T) {
	tests := []struct {
		content  string
		expected int
	}{
		{"Злая шкатулка погибает. Получено опыта: 2873.", 2873},
		{"Часы погибает. Получено опыта: 17530.", 17530},
		{"Росинка погибает.", 0},
		{"Крупье погибает. Получено опыта: 11469.", 11469},
		{"Росинка погибает. Получено опыта: 3.", 3},
	}

	for _, tt := range tests {
		result := parseLogEntry("test", tt.content)
		if result == nil {
			if tt.expected != 0 {
				t.Errorf("Failed to parse: %s", tt.content)
			}
			continue
		}
		if result.ExpGained != tt.expected {
			t.Errorf("Content: %s\nGot: %d, Expected: %d", tt.content, result.ExpGained, tt.expected)
		}
	}
}

// Helper functions for testing
func writeTempFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}

func deleteTempFile(filename string) error {
	return os.Remove(filename)
}
