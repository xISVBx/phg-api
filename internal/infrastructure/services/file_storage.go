package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	appsvc "photogallery/api_go/internal/application/interfaces/services"
)

type LocalFileStorageService struct {
	basePath       string
	maxSizeBytes   int64
	allowedMIMESet map[string]struct{}
}

var fileNameSanitizer = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)

func NewLocalFileStorageService(basePath string, maxMB int64, allowed []string) *LocalFileStorageService {
	set := map[string]struct{}{}
	for _, m := range allowed {
		v := strings.TrimSpace(strings.ToLower(m))
		if v != "" {
			set[v] = struct{}{}
		}
	}
	return &LocalFileStorageService{basePath: basePath, maxSizeBytes: maxMB * 1024 * 1024, allowedMIMESet: set}
}

func (s *LocalFileStorageService) Store(ctx context.Context, fh *multipart.FileHeader, relativeDir string) (*appsvc.StoredFile, error) {
	if s.basePath == "" {
		return nil, errors.New("FILES_BASE_PATH is empty")
	}
	if fh.Size > s.maxSizeBytes {
		return nil, fmt.Errorf("file exceeds max size: %d", fh.Size)
	}
	ct := strings.ToLower(strings.TrimSpace(fh.Header.Get("Content-Type")))
	if _, ok := s.allowedMIMESet[ct]; len(s.allowedMIMESet) > 0 && !ok {
		return nil, fmt.Errorf("mime not allowed: %s", ct)
	}
	cleanDir := filepath.Clean(relativeDir)
	if strings.Contains(cleanDir, "..") {
		return nil, errors.New("invalid relative path")
	}
	fileName := sanitizeFileName(fh.Filename)
	targetDir := filepath.Join(s.basePath, cleanDir)
	if err := os.MkdirAll(targetDir, 0o755); err != nil {
		return nil, err
	}
	stamped := fmt.Sprintf("%d_%s", time.Now().UTC().UnixNano(), fileName)
	abs := filepath.Join(targetDir, stamped)
	f, err := fh.Open()
	if err != nil {
		return nil, err
	}
	defer f.Close()
	out, err := os.Create(abs)
	if err != nil {
		return nil, err
	}
	defer out.Close()
	if _, err := io.Copy(out, f); err != nil {
		return nil, err
	}
	rel := filepath.ToSlash(filepath.Join(cleanDir, stamped))
	if _, err := s.OpenAbsPath(rel); err != nil {
		return nil, err
	}
	_ = ctx
	baseKey := strings.Split(rel, "/")[0]
	return &appsvc.StoredFile{OriginalName: fh.Filename, ContentType: ct, SizeBytes: fh.Size, StorageRelativePath: rel, StorageBaseKey: baseKey}, nil
}

func (s *LocalFileStorageService) Delete(_ context.Context, relativePath string) error {
	abs, err := s.OpenAbsPath(relativePath)
	if err != nil {
		return err
	}
	if err := os.Remove(abs); err != nil && !os.IsNotExist(err) {
		return err
	}
	return nil
}

func (s *LocalFileStorageService) OpenAbsPath(relativePath string) (string, error) {
	if s.basePath == "" {
		return "", errors.New("FILES_BASE_PATH is empty")
	}
	cleanRel := filepath.Clean(relativePath)
	if strings.Contains(cleanRel, "..") {
		return "", errors.New("path traversal detected")
	}
	absBase, err := filepath.Abs(s.basePath)
	if err != nil {
		return "", err
	}
	absPath := filepath.Join(absBase, cleanRel)
	absPath, err = filepath.Abs(absPath)
	if err != nil {
		return "", err
	}
	if !strings.HasPrefix(absPath, absBase) {
		return "", errors.New("invalid file path")
	}
	return absPath, nil
}

func sanitizeFileName(v string) string {
	base := filepath.Base(v)
	base = fileNameSanitizer.ReplaceAllString(base, "_")
	if base == "" || base == "." {
		return "file.bin"
	}
	return base
}
