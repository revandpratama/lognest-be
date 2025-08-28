package azurestorage

import (
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	"math/big"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob/sas"
	"github.com/revandpratama/lognest/config"
)

// Pre-compile regex patterns for efficiency.
var (
	// Matches any character that is not a letter, number, or space.
	unsafeChars = regexp.MustCompile(`[^a-z0-9 ]+`)
	// Matches one or more consecutive spaces.
	multipleSpaces = regexp.MustCompile(` +`)
)

const randomStringLength = 5
const randomCharset = "abcdefghijklmnopqrstuvwxyz0123456789"

// SanitizeFileName cleans and formats a filename to be URL and filesystem-friendly.
// The final format is "sanitized-name-timestamp-random.extension".
func SanitizeFileName(fileName string) string {
	// 1. Separate the base name and the extension.
	extension := filepath.Ext(fileName)
	baseName := strings.TrimSuffix(fileName, extension)

	// 2. Sanitize the base name.
	//    a. Convert to lowercase.
	sanitizedBase := strings.ToLower(baseName)
	//    b. Remove all unsafe characters.
	sanitizedBase = unsafeChars.ReplaceAllString(sanitizedBase, "")
	//    c. Replace one or more spaces with a single hyphen.
	sanitizedBase = multipleSpaces.ReplaceAllString(sanitizedBase, "-")
	//    d. Trim any leading/trailing hyphens from edge cases.
	sanitizedBase = strings.Trim(sanitizedBase, "-")

	// Truncate the sanitized name to a reasonable length (e.g., 50 chars).
	if len(sanitizedBase) > 50 {
		sanitizedBase = sanitizedBase[:50]
	}

	// Handle cases where the name becomes empty after sanitization.
	if sanitizedBase == "" {
		sanitizedBase = "file"
	}

	// 3. Generate the timestamp.
	timestamp := time.Now().Unix()

	// 4. Generate a truly random string using crypto/rand.
	randomBytes := make([]byte, randomStringLength)
	for i := range randomBytes {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(randomCharset))))
		if err != nil {
			// Fallback to a simple timestamp-based string if crypto/rand fails
			return fmt.Sprintf("%s-%d%s", sanitizedBase, timestamp, extension)
		}
		randomBytes[i] = randomCharset[num.Int64()]
	}
	randomString := string(randomBytes)

	// 5. Combine all parts and return the final filename.
	return fmt.Sprintf("%s-%d-%s%s", sanitizedBase, timestamp, randomString, extension)
}

func UploadFile(ctx context.Context, azblobClient *azblob.Client, containerName string, pathName string, fileName string, file io.Reader) (string, error) {
	log.Printf("Azure client USED in Usecase. Memory Address: %p", azblobClient)
	containerClient := azblobClient.ServiceClient().NewContainerClient(containerName)

	fileFullPath := fmt.Sprintf("%s/%s", pathName, fileName)
	// fileFullPath := fileName

	blobClient := containerClient.NewBlockBlobClient(fileFullPath)

	_, err := blobClient.UploadStream(ctx, file, nil)
	if err != nil {
		return "", err
	}

	return blobClient.URL(), nil
}

func GetFileURL(ctx context.Context, azblobClient *azblob.Client, containerName string, filePath string) (string, error) {

	permissions := sas.BlobPermissions{
		Read: true,
	}

	duration, err := strconv.Atoi(config.ENV.AZURE_STORAGE_URL_EXPIRY_DURATION_IN_MINUTES)
	if err != nil {
		duration = 15
	}

	expiry := time.Now().Add(time.Duration(duration) * time.Minute)
	containerClient := azblobClient.ServiceClient().NewContainerClient(containerName)
	blobClient := containerClient.NewBlockBlobClient(filePath)

	sasURL, err := blobClient.GetSASURL(permissions, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate SAS URL: %w", err)
	}

	return sasURL, nil
}

func DeleteFile(ctx context.Context, azblobClient *azblob.Client, containerName string, filePath string) error {

	containerClient := azblobClient.ServiceClient().NewContainerClient(containerName)
	blobClient := containerClient.NewBlockBlobClient(filePath)

	_, err := blobClient.Delete(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
