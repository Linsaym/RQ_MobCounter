package stats

import (
	"strings"
	"testing"

	"RQ_MobCounter/parser"
)

func TestCalculate(t *testing.T) {
	entries := []parser.LogEntry{
		{
			Timestamp:   "1/16 06:45:41",
			MonsterName: "Злая шкатулка",
			ExpGained:   2873,
		},
		{
			Timestamp:   "1/16 06:45:53",
			MonsterName: "Часы",
			ExpGained:   17530,
		},
		{
			Timestamp:   "1/16 06:45:55",
			MonsterName: "Злая шкатулка",
			ExpGained:   2873,
		},
		{
			Timestamp:   "1/16 06:51:17",
			MonsterName: "Росинка",
			ExpGained:   0,
		},
		{
			Timestamp:   "1/16 06:51:26",
			MonsterName: "Луговая Жужа",
			ExpGained:   0,
		},
	}

	calculator := NewCalculator(entries)
	result := calculator.Calculate("count", 0)

	if len(result) != 4 {
		t.Errorf("Expected 4 unique monsters, got %d", len(result))
	}

	// Check first entry (should be sorted by kill count descending)
	firstMonster := result[0]
	if firstMonster.Name != "Злая шкатулка" {
		t.Errorf("First monster should be 'Злая шкатулка', got %q", firstMonster.Name)
	}
	if firstMonster.KillCount != 2 {
		t.Errorf("Злая шкатулка kill count: expected 2, got %d", firstMonster.KillCount)
	}
	if firstMonster.TotalExp != 5746 {
		t.Errorf("Злая шкатулка total exp: expected 5746, got %d", firstMonster.TotalExp)
	}

	// Check monster without exp
	for _, monster := range result {
		if monster.Name == "Росинка" {
			if monster.KillCount != 1 {
				t.Errorf("Росинка kill count: expected 1, got %d", monster.KillCount)
			}
			if monster.TotalExp != 0 {
				t.Errorf("Росинка exp: expected 0, got %d", monster.TotalExp)
			}
			break
		}
	}
}

func TestFormatTableWithExp(t *testing.T) {
	stats := []MonsterStats{
		{
			Name:      "Злая шкатулка",
			KillCount: 3,
			TotalExp:  8619,
		},
		{
			Name:      "Росинка",
			KillCount: 1,
			TotalExp:  0,
		},
	}

	output := FormatTable(stats, true)

	// Check if output contains expected elements
	if !strings.Contains(output, "Злая шкатулка") {
		t.Errorf("Output should contain 'Злая шкатулка'")
	}
	if !strings.Contains(output, "8,619") {
		t.Errorf("Output should contain '8,619'")
	}
	if !strings.Contains(output, "Монстр") {
		t.Errorf("Output should contain header 'Монстр'")
	}
	if !strings.Contains(output, "Суммарный опыт") {
		t.Errorf("Output should contain 'Суммарный опыт' when exp flag is true")
	}
}

func TestFormatTableWithoutExp(t *testing.T) {
	stats := []MonsterStats{
		{
			Name:      "Злая шкатулка",
			KillCount: 3,
			TotalExp:  8619,
		},
		{
			Name:      "Росинка",
			KillCount: 1,
			TotalExp:  0,
		},
	}

	output := FormatTable(stats, false)

	// Check if output contains expected elements
	if !strings.Contains(output, "Злая шкатулка") {
		t.Errorf("Output should contain 'Злая шкатулка'")
	}
	if strings.Contains(output, "8619") {
		t.Errorf("Output should not contain exp values when exp flag is false")
	}
	if !strings.Contains(output, "Монстр") {
		t.Errorf("Output should contain header 'Монстр'")
	}
	if !strings.Contains(output, "3") {
		t.Errorf("Output should contain kill count '3'")
	}
}

func TestEmptyStats(t *testing.T) {
	entries := []parser.LogEntry{}
	calculator := NewCalculator(entries)
	result := calculator.Calculate("count", 0)

	if len(result) != 0 {
		t.Errorf("Expected 0 stats for empty entries, got %d", len(result))
	}

	output := FormatTable(result, true)
	if !strings.Contains(output, "Нет данных") {
		t.Errorf("Output should contain 'Нет данных' for empty stats")
	}
}

func TestSortingByKillCount(t *testing.T) {
	entries := []parser.LogEntry{
		{Timestamp: "1", MonsterName: "A", ExpGained: 100},
		{Timestamp: "2", MonsterName: "B", ExpGained: 100},
		{Timestamp: "3", MonsterName: "B", ExpGained: 100},
		{Timestamp: "4", MonsterName: "C", ExpGained: 100},
		{Timestamp: "5", MonsterName: "C", ExpGained: 100},
		{Timestamp: "6", MonsterName: "C", ExpGained: 100},
	}

	calculator := NewCalculator(entries)
	result := calculator.Calculate("count", 0)

	// Should be sorted by kill count descending
	if result[0].Name != "C" || result[0].KillCount != 3 {
		t.Errorf("First should be C with 3 kills, got %s with %d", result[0].Name, result[0].KillCount)
	}
	if result[1].Name != "B" || result[1].KillCount != 2 {
		t.Errorf("Second should be B with 2 kills, got %s with %d", result[1].Name, result[1].KillCount)
	}
	if result[2].Name != "A" || result[2].KillCount != 1 {
		t.Errorf("Third should be A with 1 kill, got %s with %d", result[2].Name, result[2].KillCount)
	}
}

func TestTruncateString(t *testing.T) {
	tests := []struct {
		input    string
		maxLen   int
		expected string
	}{
		{"Short", 10, "Short"},
		{"This is a very long monster name", 10, "This is..."},
		{"Exactly", 7, "Exactly"},
		{"Exactly", 6, "Exa..."},
	}

	for _, tt := range tests {
		result := truncateString(tt.input, tt.maxLen)
		if result != tt.expected {
			t.Errorf("truncateString(%q, %d): got %q, want %q", tt.input, tt.maxLen, result, tt.expected)
		}
	}
}

func TestDuplicateMonsters(t *testing.T) {
	entries := []parser.LogEntry{
		{Timestamp: "1", MonsterName: "Слизь", ExpGained: 100},
		{Timestamp: "2", MonsterName: "Слизь", ExpGained: 100},
		{Timestamp: "3", MonsterName: "Слизь", ExpGained: 100},
	}

	calculator := NewCalculator(entries)
	result := calculator.Calculate("count", 0)

	if len(result) != 1 {
		t.Errorf("Expected 1 unique monster, got %d", len(result))
	}

	if result[0].Name != "Слизь" {
		t.Errorf("Expected 'Слизь', got %q", result[0].Name)
	}

	if result[0].KillCount != 3 {
		t.Errorf("Expected 3 kills, got %d", result[0].KillCount)
	}

	if result[0].TotalExp != 300 {
		t.Errorf("Expected 300 total exp, got %d", result[0].TotalExp)
	}
}

func TestMixedExpAndNoExp(t *testing.T) {
	entries := []parser.LogEntry{
		{Timestamp: "1", MonsterName: "A", ExpGained: 100},
		{Timestamp: "2", MonsterName: "A", ExpGained: 0},
		{Timestamp: "3", MonsterName: "A", ExpGained: 100},
	}

	calculator := NewCalculator(entries)
	result := calculator.Calculate("count", 0)

	if len(result) != 1 {
		t.Errorf("Expected 1 monster, got %d", len(result))
	}

	monster := result[0]
	if monster.KillCount != 3 {
		t.Errorf("Expected 3 kills, got %d", monster.KillCount)
	}

	if monster.TotalExp != 200 {
		t.Errorf("Expected 200 total exp (100+0+100), got %d", monster.TotalExp)
	}
}

func TestFormatNumberForDisplay(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{0, "0"},
		{1, "1"},
		{100, "100"},
		{1000, "1,000"},
		{10000, "10,000"},
		{100000, "100,000"},
		{1000000, "1,000,000"},
		{10000000, "10,000,000"},
		{102413, "102,413"},
		{22984, "22,984"},
	}

	for _, tt := range tests {
		result := FormatNumberForDisplay(tt.input)
		if result != tt.expected {
			t.Errorf("FormatNumberForDisplay(%d) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestSortingByExp(t *testing.T) {
	entries := []parser.LogEntry{
		{Timestamp: "1", MonsterName: "A", ExpGained: 100},
		{Timestamp: "2", MonsterName: "B", ExpGained: 200},
		{Timestamp: "3", MonsterName: "B", ExpGained: 200},
		{Timestamp: "4", MonsterName: "C", ExpGained: 300},
		{Timestamp: "5", MonsterName: "C", ExpGained: 300},
		{Timestamp: "6", MonsterName: "C", ExpGained: 300},
	}

	calculator := NewCalculator(entries)
	result := calculator.Calculate("exp", 0)

	// Should be sorted by total exp descending
	if result[0].Name != "C" || result[0].TotalExp != 900 {
		t.Errorf("First should be C with 900 exp, got %s with %d", result[0].Name, result[0].TotalExp)
	}
	if result[1].Name != "B" || result[1].TotalExp != 400 {
		t.Errorf("Second should be B with 400 exp, got %s with %d", result[1].Name, result[1].TotalExp)
	}
	if result[2].Name != "A" || result[2].TotalExp != 100 {
		t.Errorf("Third should be A with 100 exp, got %s with %d", result[2].Name, result[2].TotalExp)
	}
}

func TestLimit(t *testing.T) {
	entries := []parser.LogEntry{
		{Timestamp: "1", MonsterName: "A", ExpGained: 100},
		{Timestamp: "2", MonsterName: "B", ExpGained: 200},
		{Timestamp: "3", MonsterName: "C", ExpGained: 300},
		{Timestamp: "4", MonsterName: "D", ExpGained: 400},
		{Timestamp: "5", MonsterName: "E", ExpGained: 500},
	}

	calculator := NewCalculator(entries)
	result := calculator.Calculate("count", 3)

	if len(result) != 3 {
		t.Errorf("Expected 3 monsters with limit 3, got %d", len(result))
	}

	// Should be top 3 by kill count (all have 1 kill, so alphabetical)
	expected := []string{"A", "B", "C"}
	for i, exp := range expected {
		if result[i].Name != exp {
			t.Errorf("Position %d should be %s, got %s", i, exp, result[i].Name)
		}
	}
}
