package actions

const (
	SuccessResponse            = 200
	ErrorResponse              = iota + 1001 // 1002
	InvalidRequest                           // 1003
	UnauthorizedError                        // 1004
	ServerError                              // 1005
	UserExistedError                         // 1006
	PasswordMethodExistedError               // 1007
	FileExistedError                         // 1008
	ProjectNotExistError                     // 1009
)
