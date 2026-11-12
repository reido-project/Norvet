package nvterr

// Domain "file": uploads and stored files.

var (
	ErrFileNotFound       = defineChild(ErrNotFound, "file_not_found", "file not found")
	ErrFileAlreadyExists  = defineChild(ErrAlreadyExists, "file_already_exists", "a file with this name already exists")
	ErrFileAlreadyDeleted = defineChild(ErrAlreadyDeleted, "file_already_deleted", "file was already deleted")
	ErrFileExpired        = defineChild(ErrExpired, "file_expired", "file has expired")

	ErrFileEmpty = defineChild(ErrEmptyBody, "file_empty", "uploaded file is empty")

	// 413
	ErrFileTooLarge  = defineChild(ErrPayloadTooLarge, "file_too_large", "file is too large")
	ErrImageTooLarge = defineChild(ErrPayloadTooLarge, "image_too_large", "image is too large")
	ErrVideoTooLarge = defineChild(ErrPayloadTooLarge, "video_too_large", "video is too large")

	// 415
	ErrUnsupportedFileType      = defineChild(ErrUnsupportedMediaType, "unsupported_file_type", "file type is not allowed")
	ErrUnsupportedFileExtension = defineChild(ErrUnsupportedMediaType, "unsupported_file_extension", "file extension is not allowed")
	ErrUnsupportedImageFormat   = defineChild(ErrUnsupportedMediaType, "unsupported_image_format", "image format is not supported")

	// 422
	ErrInvalidFileName        = defineChild(ErrValidationFailed, "invalid_file_name", "file name is invalid")
	ErrImageDimensionsInvalid = defineChild(ErrValidationFailed, "image_dimensions_invalid", "image dimensions are not allowed")
	ErrFileCorrupted          = defineChild(ErrUnprocessable, "file_corrupted", "file is corrupted or unreadable")
	ErrFileMalwareDetected    = defineChild(ErrUnprocessable, "file_malware_detected", "file was rejected by the security scan")

	ErrStorageQuotaExceeded = defineChild(ErrQuotaExceeded, "storage_quota_exceeded", "storage quota exceeded")

	ErrFileProcessingFailed = defineChild(ErrInternal, "file_processing_failed", "internal server error")
)
