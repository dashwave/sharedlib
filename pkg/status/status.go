package status

import (
	"strconv"
	"strings"
)

type Stage string
type State string
type WorkflowStep string

type AIAgentAction string

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
	GIT_CLONE                WorkflowStep = "GIT_CLONE"
	ANDROID_BUILD            WorkflowStep = "APK_BUILD" // DEPRECATED (use GRADLE_ANDROID_BUILD)
	GRADLE_ANDROID_BUILD     WorkflowStep = "ANDROID_BUILD"
	CUSTOM_SCRIPT            WorkflowStep = "CUSTOM_SCRIPT"
	REMOTE_CACHE_SETUP       WorkflowStep = "REMOTE_CACHE_SETUP"
	APK_DETECTION            WorkflowStep = "APK_DETECTION" // DEPRECATED (use OUTPUT_DETECTION)
	OUTPUT_DETECTION         WorkflowStep = "OUTPUT_DETECTION"
	BRANCH_CHECKOUT          WorkflowStep = "BRANCH_CHECKOUT"
	ARTEFACT_UPLOAD          WorkflowStep = "ARTEFACT_UPLOAD"
	IOS_BUILD                WorkflowStep = "IOS_BUILD"
	COCOAPODS_INSTALL        WorkflowStep = "COCOAPODS_INSTALL"
	XCODE_SELECT             WorkflowStep = "XCODE_SELECT"
	IOS_CERTIFICATE_DOWNLOAD WorkflowStep = "IOS_CERTIFICATE_DOWNLOAD"
	IOS_CERTIFICATE_INSTALL  WorkflowStep = "IOS_CERTIFICATE_INSTALL"
	UPLOAD_IPA_TO_CLOUD      WorkflowStep = "UPLOAD_IPA_TO_CLOUD"
	DOWNLOAD_ASSETS          WorkflowStep = "DOWNLOAD_ASSETS"
	UPLOAD_AAB_TO_PLAYSTORE  WorkflowStep = "UPLOAD_AAB_TO_PLAYSTORE"
	PUBLISH_TO_PLAYSTORE     WorkflowStep = "PUBLISH_TO_PLAYSTORE"
	FLUTTER_BUILD            WorkflowStep = "FLUTTER_BUILD"

	//Workflow error Codes
	BUILD_LOG_FILE_OPEN_ERROR            string = "BLF_01"
	BUILDER_STATUS_UPDATE_ERROR          string = "BSR_01"
	STEP_DETAIL_ERROR                    string = "SD_01"
	STEP_VERIFICATION_ERROR              string = "SV_01"
	STEP_RUNNER_SETUP_ERROR              string = "SRS_01"
	STEP_EXECUTION_ERROR                 string = "SE_01"
	CACHE_MANAGER_SETUP_ERROR            string = "CMS_01"
	OUTPUT_PATH_RETRIEVE_ERROR           string = "OPR_01"
	OUTPUT_DETECTION_ERROR               string = "ODE_01"
	INVALID_OUTPUT_TYPE_ERROR            string = "IOTE_01"
	VALIDATE_CLONE_COMMAND_MISSING_ERROR string = "VCME_01"
	VALIDATE_CLONE_COMMAND_ERROR         string = "VCME_02"
	CLONE_SPACE_COMMAD_ERROR             string = "CSCE_01"
	CLONE_SPACE_ERROR                    string = "CSE_01"
	CLONE_COMMAND_MISSING_ERROR          string = "CCME_01"
	CLEAR_WORKSPACE_ERROR                string = "CWE_01"
	CLONE_ERROR                          string = "CE_01"
	TRANSFORM_DEPS_ERROR                 string = "TDE_01"
	CACHE_SETUP_ERROR                    string = "CSE_01"
	CUSTOM_SCRIPT_ERROR                  string = "CSE_01"
	IOS_CERTIFICATE_DOWNLOAD_ERROR       string = "ICDE_01"
	IOS_CERTIFICATE_INSTALL_ERROR        string = "ICIE_01"
	UPLOAD_IPA_TO_CLOUD_ERROR            string = "UITCE_01"
	PODS_INSTALLATION_ERROR              string = "PIE_01"
	DOWNLOAD_ASSETS_ERROR                string = "DAE_01"
	APK_PUBLISH_NOT_ALLOWED              string = "APN_01"
	UPLOAD_AAB_TO_PLAYSTORE_ERROR        string = "UAP_01"
	PUBLISH_TO_PLAYSTORE_ERROR           string = "PTP_01"

	// branch checkout error codes
	FETCH_ORGIN_ERROR     string = "FOE_01"
	CLEAN_CHANGES_ERROR   string = "CCE_01"
	CHECKOUT_BRANCH_ERROR string = "CBE_01"
	PULL_ORIGIN_ERROR     string = "POE_01"

	// Android build errors
	PREBUILD_BUILD_ERROR      string = "PBE_01"
	BUILD_PROCCESS_FAIL_ERROR string = "BPFE_01"

	//fluttr errors
	DEPENDENCY_INSTALLATION_ERROR  string = "DIE_01"
	PREBUILD_SETUP_ERRROR          string = "PSE_01"
	PREBUILD_COMMAND_ERROR         string = "PCE_01"
	GRADLE_WRAPPER_CONFIGURE_ERROR string = "GWCE_01"

	//reactnative errors
	NPM_INSTALLATION_ERROR             string = "NIE_01"
	PREFERRED_DEPENDENCY_ERROR         string = "PDE_01"
	PREFERRED_DEPENDENCY_COMMAND_ERROR string = "PDCE_01"
	EXPO_OR_VANILLA_DETECTION_ERROR    string = "EVD_01"
	CONFIGURE_EXPO_ERROR               string = "CEE_01"
	ENTRYPOINT_DETECTIOIN_ERROR        string = "EDE_01"
	BUNDLE_COMMAND_ERROR               string = "BCE_01"

	// output detection erros
	APK_SEARCH_ERROR string = "ASE_01"

	// AI Agent Actions
	REQUEST_AGENT AIAgentAction = "REQUEST_AGENT"
	FIX_ISSUE     AIAgentAction = "FIX_ISSUE"
	FIX_PR        AIAgentAction = "FIX_PR"
)

func GenerateBuildMsg(stage Stage, state State, msg string) string {
	statusMsg := []string{string(stage), string(state), msg}
	return strings.Join(statusMsg, "|")
}

func GenerateWorkflowStatusMsg(stepPosition int32, state State, errorCode, msg string) string {
	statusMsg := []string{strconv.Itoa(int(stepPosition)), string(state), errorCode, msg}
	return strings.Join(statusMsg, "|")
}
