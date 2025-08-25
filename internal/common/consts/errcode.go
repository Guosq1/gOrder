package consts

const (
	ErrcodeSuccess      = 0
	ErrcodeUnknownError = 1
	//para
)

var ErrMsg = map[int]string{
	ErrcodeSuccess:      "success",
	ErrcodeUnknownError: "unknown error",
}
