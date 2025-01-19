package hotels

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/google/uuid"
	"github.com/studio-b12/gowebdav"
)

type WebDAVService struct {
	client  *gowebdav.Client
	baseURL string
}

func NewWebDAVService(url, user, password string) *WebDAVService {
	client := gowebdav.NewClient("https://webdav.cloud.mail.ru", user, password)
	err := client.Connect()
	if err != nil {
		log.Printf("Ошибка подключения к Mail.ru Cloud: %v", err)
		return nil
	}
	return &WebDAVService{
		client:  client,
		baseURL: os.Getenv("WEVDAV_URL"), // Вставьте сюда хэш из публичной ссылки
	}
}

func (s *WebDAVService) UploadImage(data []byte, filename string) (string, error) {
	err := s.client.MkdirAll("hotel-images", 0644)
	if err != nil {
		return "", err
	}
	filepath := path.Join("/hotel-images", filename)
	err = s.client.Write(filepath, data, 0644)
	if err != nil {
		return "", err
	}

	publicURL := s.baseURL + "/" + filename
	return publicURL, nil
}

func (s *WebDAVService) DeleteImage(filename string) error {
	filepath := path.Join("/hotel-images", filename)
	return s.client.Remove(filepath)
}

// Вспомогательные функции
func isImageFile(filename string) bool {
	ext := strings.ToLower(path.Ext(filename))
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif"
}

func generateUniqueFilename(originalFilename string) string {
	ext := path.Ext(originalFilename)
	return fmt.Sprintf("%s%s", uuid.New().String(), ext)
}
