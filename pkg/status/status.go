package status

import "strings"

type Stage string
type State string

const (
	// Stages of build process
	CONFIGURE        Stage = "configure"
	APK_SEARCH       Stage = "apk-search"
	EMULATOR_CONNECT Stage = "emulator-connect"
	APK_INSTALL      Stage = "apk-install"
	BUILD            Stage = "build"

	// States of build's stages
	FAILED     State = "FAILED"
	INPROGRESS State = "INPROGRESS"
	DONE       State = "DONE"
)

func GenerateBuildMsg(stage Stage, state State, msg string) string {
	statusMsg := []string{string(stage), string(state), msg}
	return strings.Join(statusMsg, "|")
}
