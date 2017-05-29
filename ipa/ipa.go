package ipa

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/bitrise-io/go-utils/ziputil"
)

func findFileInPayloadAppDir(payloadPth, preferedAppName, fileName string) (string, error) {
	appDir := filepath.Join(payloadPth, preferedAppName+".app")

	filePth := filepath.Join(appDir, fileName)
	if exist, err := pathutil.IsPathExists(filePth); err != nil {
		return "", err
	} else if exist {
		return filePth, nil
	}
	// ---

	// It's somewhere else - let's find it!
	pattern := filepath.Join(payloadPth, "*.app", fileName)
	if filePths, err := filepath.Glob(pattern); err != nil {
		return "", err
	} else if len(filePths) > 0 {
		return filePths[0], nil
	}
	// ---

	return "", fmt.Errorf("failed to find %s", fileName)
}

func unwrapFileEmbeddedInPayloadAppDir(ipaPth, fileName string) (string, error) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__ipa__")
	if err != nil {
		return "", err
	}

	if err := ziputil.UnZip(ipaPth, tmpDir); err != nil {
		return "", err
	}

	payloadPth := filepath.Join(tmpDir, "Payload")
	ipaName := strings.TrimSuffix(filepath.Base(ipaPth), filepath.Ext(ipaPth))

	return findFileInPayloadAppDir(payloadPth, ipaName, fileName)
}

// UnwrapEmbeddedMobileProvision ...
func UnwrapEmbeddedMobileProvision(ipaPth string) (string, error) {
	return unwrapFileEmbeddedInPayloadAppDir(ipaPth, "embedded.mobileprovision")
}

// UnwrapEmbeddedInfoPlist ...
func UnwrapEmbeddedInfoPlist(ipaPth string) (string, error) {
	return unwrapFileEmbeddedInPayloadAppDir(ipaPth, "Info.plist")
}
