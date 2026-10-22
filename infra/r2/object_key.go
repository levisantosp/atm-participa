package r2

import (
	"fmt"
	"mime"
	"strings"
)

func ImageExtension(contentType string) (string, error) {
	mediaType, _, err := mime.ParseMediaType(contentType)
	if err != nil {
		return "", err
	}

	switch mediaType {
	case "image/gif":
		return "gif", nil
	case "image/jpeg":
		return "jpeg", nil
	case "image/png":
		return "png", nil
	case "image/webp":
		return "webp", nil
	default:
		return "", fmt.Errorf("unsupported image content type %q", mediaType)
	}
}

func IssueFileObjectKey(issueID int64, fileID, extension string) string {
	key := fmt.Sprintf("issues/%d/%s", issueID, fileID)
	return key + "." + strings.TrimPrefix(extension, ".")
}

func IssueFileObjectPrefix(issueID int64, fileID string) string {
	return fmt.Sprintf("issues/%d/%s", issueID, fileID)
}
