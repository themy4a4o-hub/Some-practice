package service

import (
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
	defer file.Close()
	for _, line := range lines {
		file.WriteString(line + "\n")
	}
	return nil
}
func FilePresenter(path string) Presenter {
	var pr filepresent
	pr.path = path
	return pr
}
