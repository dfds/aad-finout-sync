package handler

import "github.com/joomcode/errorx"

var (
	HandlerError              = errorx.NewNamespace("handler")
	VirtualTagDoesNotExist    = HandlerError.NewType("tag_does_not_exist")
	VirtualTagDoesNotExistMsg = "Unable to find virtual tag"
)
