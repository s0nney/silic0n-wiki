package database

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
)

func RunMigrations(migrationsPath string) error {
	files, err := os.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var sqlFiles []string
	for _, f := range files {
		if filepath.Ext(f.Name()) == ".sql" {
			sqlFiles = append(sqlFiles, f.Name())
		}
	}
	sort.Strings(sqlFiles)

	for _, file := range sqlFiles {
		path := filepath.Join(migrationsPath, file)
		content, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration %s: %w", file, err)
		}

		_, err = DB.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}

		fmt.Printf("Applied migration: %s\n", file)
	}

	return nil
}
