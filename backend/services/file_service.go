package services

import (
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"

	"github.com/nfnt/resize"
	"gorm.io/gorm"

	"sirine-go/backend/models"
)

// FileService merupakan service untuk handling file operations
// yang mencakup upload, resize, dan storage management
type FileService struct {
	db              *gorm.DB
	uploadDir       string
	maxFileSize     int64
	allowedFormats  []string
}

// NewFileService membuat instance baru dari FileService
// dengan konfigurasi upload directory dan file constraints
func NewFileService(db *gorm.DB) *FileService {
	// Buat upload directory jika belum ada
	uploadDir := "./public/uploads/profiles"
	os.MkdirAll(uploadDir, 0755)

	return &FileService{
		db:             db,
		uploadDir:      uploadDir,
		maxFileSize:    5 * 1024 * 1024, // 5MB
		allowedFormats: []string{".jpg", ".jpeg", ".png", ".webp"},
	}
}

// UploadProfilePhoto meng-upload dan resize profile photo untuk user
// dengan validation format, size, dan automatic image optimization
func (s *FileService) UploadProfilePhoto(userID uint64, fileHeader *multipart.FileHeader) (string, error) {
	// Validate file size
	if fileHeader.Size > s.maxFileSize {
		return "", errors.New("ukuran file terlalu besar, maksimal 5MB")
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	if !s.isAllowedFormat(ext) {
		return "", errors.New("format file tidak didukung, gunakan JPG, PNG, atau WebP")
	}

	// Open uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Decode image
	img, format, err := image.Decode(file)
	if err != nil {
		return "", errors.New("gagal membaca file gambar")
	}

	// Resize image to 200x200px untuk optimize storage dan performance
	resizedImg := resize.Thumbnail(200, 200, img, resize.Lanczos3)

	// Generate filename: userID.jpg
	filename := fmt.Sprintf("%d.jpg", userID)
	filePath := filepath.Join(s.uploadDir, filename)

	// Create output file
	outFile, err := os.Create(filePath)
	if err != nil {
		return "", err
	}
	defer outFile.Close()

	// Encode dan save image sebagai JPEG dengan quality optimization
	if format == "png" {
		// Convert PNG to JPEG
		if err := jpeg.Encode(outFile, resizedImg, &jpeg.Options{Quality: 90}); err != nil {
			return "", err
		}
	} else {
		// Save as JPEG
		if err := jpeg.Encode(outFile, resizedImg, &jpeg.Options{Quality: 90}); err != nil {
			return "", err
		}
	}

	// Generate public URL
	photoURL := fmt.Sprintf("/uploads/profiles/%s", filename)

	// Update user profile photo URL di database
	if err := s.db.Model(&models.User{}).
		Where("id = ?", userID).
		Update("profile_photo_url", photoURL).Error; err != nil {
		// Rollback: delete uploaded file jika database update gagal
		os.Remove(filePath)
		return "", err
	}

	return photoURL, nil
}

// DeleteProfilePhoto menghapus profile photo user dari storage dan database
func (s *FileService) DeleteProfilePhoto(userID uint64) error {
	// Get user untuk mendapatkan current photo URL
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return err
	}

	if user.ProfilePhotoURL == "" {
		return nil // Tidak ada foto untuk dihapus
	}

	// Delete file dari storage
	filename := filepath.Base(user.ProfilePhotoURL)
	filePath := filepath.Join(s.uploadDir, filename)
	
	// Ignore error jika file tidak ada
	os.Remove(filePath)

	// Update database untuk clear photo URL
	if err := s.db.Model(&user).
		Update("profile_photo_url", "").Error; err != nil {
		return err
	}

	return nil
}

// isAllowedFormat memeriksa apakah file extension diperbolehkan
func (s *FileService) isAllowedFormat(ext string) bool {
	for _, allowed := range s.allowedFormats {
		if ext == allowed {
			return true
		}
	}
	return false
}

// ValidateImage melakukan validation tambahan untuk image file
func (s *FileService) ValidateImage(file io.Reader) error {
	// Decode untuk validate bahwa ini adalah valid image
	_, _, err := image.Decode(file)
	if err != nil {
		return errors.New("file bukan gambar yang valid")
	}
	return nil
}

// GetProfilePhotoPath mendapatkan full path dari profile photo untuk serving
func (s *FileService) GetProfilePhotoPath(userID uint64) (string, error) {
	var user models.User
	if err := s.db.First(&user, userID).Error; err != nil {
		return "", err
	}

	if user.ProfilePhotoURL == "" {
		return "", errors.New("user belum memiliki foto profil")
	}

	filename := filepath.Base(user.ProfilePhotoURL)
	return filepath.Join(s.uploadDir, filename), nil
}

// ResizeAndOptimize melakukan resize dan optimization untuk existing image
// (utility function untuk batch processing atau migration)
func (s *FileService) ResizeAndOptimize(inputPath, outputPath string, width, height uint) error {
	// Open input file
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Decode image
	img, _, err := image.Decode(file)
	if err != nil {
		return err
	}

	// Resize
	resizedImg := resize.Thumbnail(width, height, img, resize.Lanczos3)

	// Create output file
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	// Encode as JPEG dengan optimization
	return jpeg.Encode(outFile, resizedImg, &jpeg.Options{Quality: 90})
}
