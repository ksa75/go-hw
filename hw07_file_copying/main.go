package main

import (
	"errors"
	"flag"
	"fmt"
	"log"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse() // проанализировать аргументы
	// теперь в from,to,offset,limit есть значения

	// Проверка, что обязательные флаги указаны
	if from == "" || to == "" {
		log.Fatal("Both source and destination file paths must be provided")
	}

	// Выполняем копирование
	err := Copy(from, to, offset, limit)
	if err != nil {
		if errors.Is(err, ErrOffsetExceedsFileSize) {
			log.Fatalf("Error: %v\n", ErrOffsetExceedsFileSize)
			if errors.Is(err, ErrUnsupportedFile) {
				log.Fatalf("Error: %v\n", ErrUnsupportedFile)
			} else {
				log.Fatalf("Error: %v\n", err)
			}
		}
	} else {
		fmt.Println("File copied successfully!")
	}
}
