// Package git предоставляет функции для работы с git репозиториями
package git

import (
	"fmt"
	"os/exec"
)

// GetStagedChanges выполняет команду git diff --staged и возвращает её вывод
func GetStagedChanges() (string, error) {
	// Проверяем, что мы находимся в git репозитории
	if !isGitRepository() {
		return "", fmt.Errorf("текущая директория не является git репозиторием")
	}

	// Выполняем команду git diff --staged
	cmd := exec.Command("git", "diff", "--staged")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("ошибка при выполнении git diff --staged: %v", err)
	}

	return string(output), nil
}

// isGitRepository проверяет, является ли текущая директория git репозиторием
func isGitRepository() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	err := cmd.Run()
	return err == nil
}