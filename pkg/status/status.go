package status

import "strings"

type Stage string
type State string
type WorkflowStep string

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

	// Workflow steps
	GIT_CLONE          WorkflowStep = "GIT_CLONE"
	ANDROID_BUILD      WorkflowStep = "APK_BUILD"
	CUSTOM_SCRIPT      WorkflowStep = "CUSTOM_SCRIPT"
	REMOTE_CACHE_SETUP WorkflowStep = "REMOTE_CACHE_SETUP"
	APK_DETECTION      WorkflowStep = "APK_DETECTION"
)

func GenerateBuildMsg(stage Stage, state State, msg string) string {
	statusMsg := []string{string(stage), string(state), msg}
	return strings.Join(statusMsg, "|")
}
