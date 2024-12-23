package constants

const RoleAdmin = "admin"
const RoleEditor = "editor"
const RoleUser = "normal"

const RoleAdminMaxCodeAmount = 100
const RoleEditorMaxCodeAmount = 5
const RoleUserMaxCodeAmount = 2

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

const (
	TransactionTypeRecharge = "recharge"
	TransactionTypeGifted   = "gifted"
	TransactionTypeAdmin    = "admin"
	TransactionTypeWallet   = "wallet"
	TransactionTypeUseAgent = "use_agent"
)

const (
	GiftedPoints = 100
)

const DifyServerPrefix = "/hajime_federation"

const (
	// 日常
	RegisterPoints            = 1000
	WalletLinkPoints          = 100
	DailySignInPoints         = 10
	UseBotAgentWorkflowPoints = 20

	// 控制成本
	ChatCostPerToken = -0.5
	FreeChatSessions = 3

	// 鼓励Devs
	CreateChatbotPoints   = 10000
	CreateWorkflowPoints  = 30000
	UseVariablesPoints    = 2000
	UseToolsPoints        = 3000
	UploadKnowledgePoints = 1000
	UseKnowledgePoints    = 2000
	UsageEarnPerToken     = 0.5

	// 运营 / 市场
	HajimeBotHolderMultiplier        = 1.5
	InvitationBonusRate              = 0.1 // 10% of invite user amount (balance*(1 + 0.1*n))
	InvitationUserGetBalanceRate     = 0.2 // 20% of the balance of the invited user (0.2*balance)
)
