package stats

import (
	"fmt"
	"sort"
	"strings"

	"RQ_MobCounter/parser"
)

type MonsterStats struct {
	Name      string
	KillCount int
	TotalExp  int
}

type Calculator struct {
	entries []parser.LogEntry
}

func NewCalculator(entries []parser.LogEntry) *Calculator {
	return &Calculator{
		entries: entries,
	}
}

func (c *Calculator) Calculate(sortBy string, limit int) []MonsterStats {
	statsMap := make(map[string]*MonsterStats)

	for _, entry := range c.entries {
		if entry.MonsterName == "" {
			continue
		}

		if _, exists := statsMap[entry.MonsterName]; !exists {
			statsMap[entry.MonsterName] = &MonsterStats{
				Name: entry.MonsterName,
			}
		}

		stats := statsMap[entry.MonsterName]
		stats.KillCount++
		stats.TotalExp += entry.ExpGained
	}

	var result []MonsterStats
	for _, stat := range statsMap {
		result = append(result, *stat)
	}

	sort.Slice(result, func(i, j int) bool {
		if sortBy == "exp" {
			if result[i].TotalExp != result[j].TotalExp {
				return result[i].TotalExp > result[j].TotalExp
			}
		} else {
			if result[i].KillCount != result[j].KillCount {
				return result[i].KillCount > result[j].KillCount
			}
		}
		return result[i].Name < result[j].Name
	})

	if limit > 0 && len(result) > limit {
		result = result[:limit]
	}

	return result
}

func FormatTable(stats []MonsterStats, showExp bool) string {
	if len(stats) == 0 {
		return "Нет данных для отображения\n"
	}

	output := ""

	if showExp {
		output += fmt.Sprintf("%-40s | %15s | %15s\n",
			"Монстр", "Количество", "Суммарный опыт")
		output += strings.Repeat("-", 75) + "\n"

		for _, s := range stats {
			output += fmt.Sprintf("%-40s | %15d | %15s\n",
				truncateString(s.Name, 40),
				s.KillCount,
				FormatNumberForDisplay(s.TotalExp))
		}
	} else {
		output += fmt.Sprintf("%-40s | %15s\n", "Монстр", "Количество")
		output += strings.Repeat("-", 60) + "\n"

		for _, s := range stats {
			output += fmt.Sprintf("%-40s | %15d\n",
				truncateString(s.Name, 40),
				s.KillCount)
		}
	}

	return output
}

func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func FormatNumberForDisplay(n int) string {
	s := fmt.Sprintf("%d", n)
	if len(s) <= 3 {
		return s
	}

	var result strings.Builder
	for i, c := range s {
		if i > 0 && (len(s)-i)%3 == 0 {
			result.WriteRune(',')
		}
		result.WriteRune(c)
	}
	return result.String()
}
