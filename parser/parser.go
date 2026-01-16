package parser

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

type LogEntry struct {
	Timestamp   string
	MonsterName string
	ExpGained   int
}

func ParseFile(filepath string) ([]LogEntry, error) {
	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %w", err)
	}

	content := string(data)
	var entries []LogEntry

	trRegex := regexp.MustCompile(`<TR[^>]*title='([^']+)'[^>]*><TD[^>]*>([^\n]+)`)

	matches := trRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}

		timestamp := match[1]
		contentStr := match[2]

		entry := parseLogEntry(timestamp, contentStr)
		if entry != nil {
			entries = append(entries, *entry)
		}
	}

	return entries, nil
}

func parseLogEntry(timestamp string, content string) *LogEntry {
	content = regexp.MustCompile(`<[^>]*>`).ReplaceAllString(content, "")
	content = strings.TrimSpace(content)

	if content == "" {
		return nil
	}

	if !strings.Contains(content, "погибает") {
		return nil
	}

	entry := &LogEntry{
		Timestamp: timestamp,
		ExpGained: 0,
	}

	parts := strings.Split(content, "погибает")
	if len(parts) > 0 {
		monsterName := strings.TrimSpace(parts[0])
		if monsterName != "" && !strings.HasPrefix(monsterName, "Вы") {
			entry.MonsterName = monsterName
		}
	}

	if entry.MonsterName == "" {
		return nil
	}

	expRegex := regexp.MustCompile(`Получено опыт[а]?:\s*(\d+)`)
	expMatch := expRegex.FindStringSubmatch(content)

	if len(expMatch) > 1 {
		if exp, err := strconv.Atoi(expMatch[1]); err == nil {
			entry.ExpGained = exp
		}
	}

	return entry
}
