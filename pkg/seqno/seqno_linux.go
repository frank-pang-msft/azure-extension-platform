package seqno

import (
	"fmt"
	"github.com/D1v38om83r/azure-extension-platform/pkg/constants"
	"io/ioutil"
	"os"
	"strconv"
)

var mostRecentSequenceFileName = "mrseq"

// sequence number for the extension from the registry
func getSequenceNumberInternal(name, version string) (uint, error) {
	mrseqStr, err := ioutil.ReadFile(mostRecentSequenceFileName)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, fmt.Errorf("failed to read mrseq file : %s", err)
	}

	seqNum, err := strconv.Atoi(string(mrseqStr))
	if err != nil {
		return 0, err
	}
	return uint(seqNum), nil

}

func setSequenceNumberInternal(extName, extVersion string, seqNo uint) error {
	b := []byte(fmt.Sprintf("%v", seqNo))
	err := ioutil.WriteFile(mostRecentSequenceFileName, b, constants.FilePermissions_UserOnly_ReadWrite)
	if err != nil {
		return fmt.Errorf("could not write sequence number file %s, error: %v", mostRecentSequenceFileName, err)
	}
	return nil
}
