package input

type FileStatus int

const (
	FILE_STATUS_INIT          FileStatus = 1
	FILE_STATUS_PROCESSING    FileStatus = 2
	FILE_STATUS_BODY_READ     FileStatus = 3
	FILE_STATUS_HEADERS_READ  FileStatus = 4
	FILE_STATUS_HEADERS_PARSE FileStatus = 5
	FILE_STATUS_DONE          FileStatus = 99
	FILE_STATUS_ERROR         FileStatus = 0

	HeaderToKey          = "To"
	HeaderFromKey        = "From"
	HeaderSubjectKey     = "Subject"
	HeaderContentTypeKey = "Content-Type"
	HeaderMsgIDKey       = "Message-ID"
)
