package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrFile                  = errors.New("problem with file specified")
)

// Helper function to check if a byte is an invisible character.
func nonCountable(b byte) bool {
	return b == '\r' // || b == '\r' //|| b == '\t' //|| b == ' '
}

// Function to fix the offset and limit for excluding invisible characters.
func adjLen(src io.Reader, offset, limit int64) (int64, int64, error) {
	var adjustedOffset, adjustedLimit int64
	var readPos, readLimit int64

	for i := int64(0); readPos <= offset; i++ {
		var b [1]byte
		_, err := src.Read(b[:])
		if err != nil && err != io.EOF {
			return 0, 0, err
		}
		if nonCountable(b[0]) {
			continue
		}
		readPos++
		adjustedOffset = i
		if err == io.EOF {
			break
		}
	}

	for i := int64(0); readLimit <= limit; i++ {
		var b [1]byte
		_, err := src.Read(b[:])
		if err != nil && err != io.EOF {
			return 0, 0, err
		}
		if nonCountable(b[0]) {
			continue
		}
		adjustedLimit = i
		if err == io.EOF {
			adjustedLimit++
			break
		}
		readLimit++
	}

	return adjustedOffset, adjustedLimit, nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Открываем исходный файл
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", ErrFile)
	}
	defer fromFile.Close()

	// Получаем информацию о файле (размер)
	fileInfo, err := fromFile.Stat()
	if err != nil {
		// программа может НЕ обрабатывать файлы, у которых неизвестна длина (например, /dev/urandom);
		return fmt.Errorf("undefined src file length: %w", ErrUnsupportedFile)
	}
	fromFileLen := fileInfo.Size()

	// offset больше, чем размер файла - невалидная ситуация;
	if offset > fromFileLen {
		return fmt.Errorf("wrong offset: %w", ErrOffsetExceedsFileSize)
	}

	// считаются все кроме \n
	offset, limit, err = adjLen(fromFile, offset, limit)
	if err != nil {
		return fmt.Errorf("error adjusting offset, limit: %w", err)
	}

	// Если limit не задан, копируем до конца файла
	// limit больше, чем размер файла - копируется исходный файл до его EOF
	if limit == 0 || offset+limit > fromFileLen {
		limit = fromFileLen - offset
	}

	// Позиционируем указатель на начало копирования
	_, err = fromFile.Seek(offset, io.SeekStart)
	if err != nil {
		return fmt.Errorf("failed to seek to offset: %w", err)
	}

	// Открываем целевой файл
	toFile, err := os.Create(toPath)
	if err != nil {
		return fmt.Errorf("failed to create target file: %w", ErrFile)
	}
	defer toFile.Close()

	// Создаем прогресс-бар
	bar := pb.New(int(limit)).SetWidth(50)
	bar.Start()

	N := int(limit)
	buffer := make([]byte, N)
	var totalBytesCopied int64

	// Копируем данные с прогрессом
	for totalBytesCopied < limit {
		// Читаем из исходного файла
		bytesRead, err := fromFile.Read(buffer)
		if err != nil && !errors.Is(err, io.EOF) {
			return fmt.Errorf("failed to read from source file: %w", ErrFile)
		}

		// Пишем в целевой файл
		bytesWritten, err := toFile.Write(buffer[:bytesRead])
		if err != nil {
			return fmt.Errorf("failed to write to target file: %w", ErrFile)
		}

		// Обновляем прогресс-бар
		totalBytesCopied += int64(bytesWritten)
		bar.Add(bytesWritten)
	}

	bar.Finish()
	return nil
}
