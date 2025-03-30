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

func Copy(fromPath, toPath string, offset, limit int64) error {
	// Открываем исходный файл
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", ErrFile)
	}
	defer fromFile.Close()

	// Получаем информацию о файле
	fileInfo, _ := fromFile.Stat()

	// Проверяем, является ли файл обычным
	if fileInfo.Mode()&os.ModeType == os.ModeNamedPipe || fileInfo.Mode()&os.ModeDevice != 0 {
		// Если это именованный канал или устройство (например, /dev/urandom), длина может быть неизвестна
		return fmt.Errorf("undefined src file length: %w", ErrUnsupportedFile)
	}
	// Проверяем доступность длины для обычных файлов
	fromFileLen := fileInfo.Size()
	if fromFileLen == -1 {
		// программа может НЕ обрабатывать файлы, у которых неизвестна длина;
		return fmt.Errorf("undefined src file length: %w", ErrUnsupportedFile)
	}

	// offset больше, чем размер файла - невалидная ситуация;
	if offset > fromFileLen {
		return fmt.Errorf("wrong offset: %w", ErrOffsetExceedsFileSize)
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
