package status

import "strings"

type Stage string
type State string

const (
	// Stages of build process
	CONFIGURE        Stage = "CONFIGURE"
	APK_SEARCH       Stage = "APK-SEARCH"
	EMULATOR_CONNECT Stage = "EMULATOR-CONNECT"
	APK_INSTALL      Stage = "APK-INSTALL"
	BUILD            Stage = "BUILD"

	// States of build's stages
	FAILED     State = "FAILED"
	INPROGRESS State = "INPROGRESS"
	DONE       State = "DONE"
	CANCELLED  State = "CANCELLED"
)

func GenerateBuildMsg(stage Stage, state State, msg string) string {
	statusMsg := []string{string(stage), string(state), msg}
	return strings.Join(statusMsg, "|")
}
