package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/hazadus/gh-commitmsg/internal/git"
)

func main() {
	stagedChanges, err := git.GetStagedChanges()
	if err != nil {
		fmt.Printf("Ошибка при получении staged изменений: %v\n", err)
		os.Exit(1)
	}

	if stagedChanges == "" {
		fmt.Println("Нет staged изменений в репозитории.")
		return
	}

	fmt.Println("Staged изменения:")
	fmt.Println(strings.Repeat("-", 50))
	fmt.Print(stagedChanges)
}

