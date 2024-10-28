package constants

const RoleAdmin = "admin"
const RoleDeveloper = "developer"
const RoleUser = "user"

// GPT pricing charge rate per token
const (
	GPT3CompletionCharge = 0.002 / 1000
	GPT3PromptCharge     = 0.002 / 1000
)

const (
	GPT4CompletionCharge = 0.06 / 1000
	GPT4PromptCharge     = 0.03 / 1000
)

const DollarToChineseCentsRate = 1100

const (
	RechargingCardActive   = "active"
	RechargingCardInactive = "inactive"
	RechargingCardUsed     = "used"
)

const (
	TransactionTypeRecharge = "recharge"
	TransactionTypeAdmin    = "admin"
)

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
