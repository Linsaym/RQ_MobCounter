package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"RQ_MobCounter/config"
	"RQ_MobCounter/parser"
	"RQ_MobCounter/stats"
)

func main() {
	showExp := flag.Bool("exp", false, "показывать опыт")
	month := flag.String("month", "", "анализ конкретного месяца (YYYY.MM)")
	all := flag.Bool("all", false, "обработка всех файлов")

	flag.Parse()

	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("ошибка загрузки конфига: %v", err)
	}

	if _, err := os.Stat(cfg.LogPath); err != nil {
		log.Fatalf("путь к логам не найден: %s", cfg.LogPath)
	}

	var filesToProcess []string

	if *all {
		files, err := os.ReadDir(cfg.LogPath)
		if err != nil {
			log.Fatalf("ошибка чтения директории: %v", err)
		}

		for _, file := range files {
			if !file.IsDir() && strings.HasSuffix(file.Name(), ".htm") {
				filesToProcess = append(filesToProcess, filepath.Join(cfg.LogPath, file.Name()))
			}
		}
	} else if *month != "" {
		fileName := fmt.Sprintf("exp (%s).htm", *month)
		filePath := filepath.Join(cfg.LogPath, fileName)

		if _, err := os.Stat(filePath); err != nil {
			log.Fatalf("файл для месяца %s не найден", *month)
		}

		filesToProcess = append(filesToProcess, filePath)
	} else {
		currentMonth := getLastMonth()
		fileName := fmt.Sprintf("exp (%s).htm", currentMonth)
		filePath := filepath.Join(cfg.LogPath, fileName)

		if _, err := os.Stat(filePath); err != nil {
			fmt.Printf("файл для текущего месяца %s не найден. Доступные файлы:\n", currentMonth)
			listAvailableFiles(cfg.LogPath)
			return
		}

		filesToProcess = append(filesToProcess, filePath)
	}

	if len(filesToProcess) == 0 {
		fmt.Println("нет файлов для обработки")
		return
	}

	var allEntries []parser.LogEntry

	for _, filePath := range filesToProcess {
		fmt.Printf("обработка: %s\n", filepath.Base(filePath))

		entries, err := parser.ParseFile(filePath)
		if err != nil {
			log.Printf("ошибка при парсинге %s: %v", filePath, err)
			continue
		}

		allEntries = append(allEntries, entries...)
		fmt.Printf("найдено записей: %d\n\n", len(entries))
	}

	calculator := stats.NewCalculator(allEntries)
	monsterStats := calculator.Calculate()

	if len(filesToProcess) > 1 {
		fmt.Println("=== ОБЩАЯ СТАТИСТИКА ===")
		fmt.Println()
	}

	fmt.Print(stats.FormatTable(monsterStats, *showExp))

	fmt.Printf("\nВсего записей: %d\n", len(allEntries))
	totalKills := 0
	totalExp := 0
	for _, m := range monsterStats {
		totalKills += m.KillCount
		totalExp += m.TotalExp
	}
	if *showExp {
		fmt.Printf("Всего опыта: %d\n", totalExp)
	}
}

func getLastMonth() string {
	now := time.Now()
	return fmt.Sprintf("%d.%02d", now.Year(), now.Month())
}

func listAvailableFiles(logPath string) {
	files, err := os.ReadDir(logPath)
	if err != nil {
		fmt.Printf("ошибка чтения директории: %v\n", err)
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".htm") {
			monthRegex := regexp.MustCompile(`exp \((.+?)\)\.htm`)
			matches := monthRegex.FindStringSubmatch(file.Name())
			if len(matches) > 1 {
				fmt.Printf("  - %s\n", matches[1])
			}
		}
	}
}
