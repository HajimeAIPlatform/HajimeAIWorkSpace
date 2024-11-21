package constants

const RoleAdmin = "admin"
const RoleEditor = "editor"
const RoleUser = "normal"

// file

var (
	// IMAGE_EXTENSIONS contains the list of allowed image file extensions
	IMAGE_EXTENSIONS = []string{"jpg", "jpeg", "png", "webp", "gif", "svg"}

	// ALLOWED_EXTENSIONS contains the list of allowed file extensions for structured data
	ALLOWED_EXTENSIONS = []string{"txt", "markdown", "md", "pdf", "html", "htm", "xlsx", "xls", "docx", "csv"}

	// UNSTRUCTURED_ALLOWED_EXTENSIONS contains the list of allowed file extensions for unstructured data
	UNSTRUCTURED_ALLOWED_EXTENSIONS = []string{"txt", "markdown", "md", "pdf", "html", "htm", "xlsx", "xls", "docx", "csv", "eml", "msg", "pptx", "ppt", "xml", "epub"}
)

const SizeMB = 1024 * 1024
