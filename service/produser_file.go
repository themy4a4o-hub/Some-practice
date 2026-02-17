package service

import (
	"bufio"
	"os"
)

type Produser interface {
	Produce() ([]string, error)
}
type fileprod struct {
	path string
}

func (prod fileprod) Produce() ([]string, error) {
	file, err := os.Open(prod.path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}

func Fileproduser(path string) Produser {
	var p fileprod
	p.path = path
	return p
}
