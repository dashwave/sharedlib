package status

type Stage string
type State string

const (
	CONFIGURE Stage = "configure"
	APK_INSTALL Stage = "apk-install"
	BUILD Stage = "build"
)

const (
	FAILED State = "failed"
	INPROGRESS State = "InProgress"
	SUCCESS State = "success"
)
