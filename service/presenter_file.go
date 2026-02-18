package service

import (
	"fmt"
	"log/slog"
	"os"
)

type Presenter interface {
	Present([]string) error
}
type filepresent struct {
	path string
}

func (pres filepresent) Present(lines []string) error {
	file, err := os.Create(pres.path)
	if err != nil {
		return err
	}
	defer func() {
		if err := file.Close(); err != nil {
			slog.Error("Ошибка закрытия файла", "Error", err)
		}
	}()
	for _, line := range lines {
		if _, err := file.WriteString(line + "\n"); err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}
	}
	return nil
}
func FilePresenter(path string) Presenter {
	var pr filepresent
	pr.path = path
	return pr
}
