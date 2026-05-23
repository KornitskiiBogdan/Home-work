package main

import (
	"errors"
	"io"
	"os"

	pb "github.com/cheggaaa/pb/v3"
)

var (
	bufferSize               = 1024
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	var fileSize int64
	var err error
	if fileSize, err = validateFile(fromPath); err != nil {
		return err
	}

	bytesToCopy, err := getBytesToCopy(fileSize, offset, limit)
	if err != nil {
		return err
	}

	inFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer inFile.Close()

	_, err = inFile.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	outFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	bar := pb.Start64(bytesToCopy)
	defer bar.Finish()

	return copyWithBuffer(outFile, inFile, bytesToCopy, bar)
}

func copyWithBuffer(dest io.Writer, src io.Reader, totalSize int64, progressBar *pb.ProgressBar) error {
	buf := make([]byte, bufferSize)
	var written int64

	for written < totalSize {
		toRead := int64(len(buf))
		remaining := totalSize - written

		if remaining < toRead {
			toRead = remaining
		}

		n, err := src.Read(buf[:toRead])

		if err != nil && !errors.Is(err, io.EOF) && !errors.Is(err, io.ErrUnexpectedEOF) {
			return err
		}

		if n == 0 {
			break
		}

		if _, err := dest.Write(buf[:n]); err != nil {
			return err
		}
		written += int64(n)

		if progressBar != nil {
			progressBar.Add64(int64(n))
		}
	}

	return nil
}

func getBytesToCopy(fileSize, offset, limit int64) (int64, error) {
	if offset > fileSize {
		return 0, ErrOffsetExceedsFileSize
	}
	available := fileSize - offset
	var bytesToCopy int64
	switch {
	case limit == 0:
		bytesToCopy = available
	case limit > available:
		bytesToCopy = available
	default:
		bytesToCopy = limit
	}
	return bytesToCopy, nil
}

func validateFile(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, err
	}

	if !info.Mode().IsRegular() {
		return 0, ErrUnsupportedFile
	}
	return info.Size(), nil
}
