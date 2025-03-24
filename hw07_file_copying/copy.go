package main

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {

	// Открываем исходный файл
	fromFile, err := os.Open(fromPath)
	if err != nil {
		return fmt.Errorf("failed to open source file: %w", err)
	}
	defer fromFile.Close()

	// Получаем информацию о файле (размер)
	fileInfo, err := fromFile.Stat()
	if err != nil {
		// * программа может НЕ обрабатывать файлы, у которых неизвестна длина (например, /dev/urandom);
		return fmt.Errorf("undefined src file length: %w", ErrUnsupportedFile)
	}
	fromFileLen := int64(fileInfo.Size())

	//* offset больше, чем размер файла - невалидная ситуация;
	if offset > fromFileLen {
		return fmt.Errorf("wrong offset: %w", ErrOffsetExceedsFileSize)
	}

	//* limit больше, чем размер файла - валидная ситуация, копируется исходный файл до его EOF;
	if limit > fromFileLen {
		limit = fromFileLen
	}
	// Если limit не задан, копируем до конца файла
	if limit == 0 {
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
		return fmt.Errorf("failed to create target file: %w", err)
	}
	defer toFile.Close()

	// Создаем прогресс-бар
	bar := pb.New(int(limit)).SetWidth(50)
	bar.Start()

	// Буфер для копирования
	buffer := make([]byte, 1024)
	var totalBytesCopied int64

	// Копируем данные с прогрессом
	for {
		// Читаем из исходного файла
		bytesRead, err := fromFile.Read(buffer)
		if err != nil && err != io.EOF {
			return fmt.Errorf("failed to read from source file: %w", err)
		}
		//????
		if bytesRead == 0 {
			break
		}

		// Пишем в целевой файл
		bytesWritten, err := toFile.Write(buffer[:bytesRead])
		if err != nil {
			return fmt.Errorf("failed to write to target file: %w", err)
		}

		// Обновляем прогресс-бар
		totalBytesCopied += int64(bytesWritten)
		// if err := bar.Add(bytesWritten); err != nil {
		// return fmt.Errorf("failed to update progress: %w", err)
		// }
		bar.Add(bytesWritten)

		// Если достигнут лимит, завершаем копирование
		if totalBytesCopied >= limit {
			break
		}
	}

	bar.Finish()
	return nil
}
