package cmds

import "prick/internal/prick/common"

type AccountInfoFetchedMsg struct {
	AccountInfo common.AzAccountShowOutput
}
