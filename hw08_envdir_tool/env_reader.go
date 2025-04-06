package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read dir: %s, %v", dir, err)
	}

	for _, file := range files {
		// имя `S` не должно содержать `=`;
		if file.IsDir() || strings.Contains(file.Name(), "=") {
			continue
		}

		fPath := dir + "/" + file.Name()

		f, err := os.Open(fPath)
		if err != nil {
			return nil, fmt.Errorf("failed to open file: %s, %v", fPath, err)
		}
		defer f.Close()

		// если файл полностью пустой (длина - 0 байт), то `envdir` удаляет переменную окружения с именем `S`.
		if stat, _ := f.Stat(); stat.Size() == 0 {
			env[file.Name()] = EnvValue{NeedRemove: true}
		}
		scanner := bufio.NewScanner(f)
		if scanner.Scan() {
			// пробелы и табуляция в конце `T` удаляются;
			value := strings.TrimRight(scanner.Text(), " ")
			value = strings.TrimRight(value, "\t")
			// терминальные нули (`0x00`) заменяются на перевод строки (`\n`);
			value = strings.ReplaceAll(value, "\x00", "\n")
			if value == "" {
				env[file.Name()] = EnvValue{NeedRemove: true}
			} else {
				env[file.Name()] = EnvValue{Value: value}
			}
		} else if err := scanner.Err(); err != nil {
			return nil, fmt.Errorf("error reading file %s: %v", fPath, err)
		}
	}
	return env, nil
}
