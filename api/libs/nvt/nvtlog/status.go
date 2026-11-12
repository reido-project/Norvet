package nvtlog

import "api/libs/nvt/nvterr"

// resolveStatus picks the status in this order:
//  1. the first positive explicit code from the caller (override)
//  2. the nvterr status (also covers context.Canceled / DeadlineExceeded)
//  3. fallback (0 = unknown for logging, 500 for responses)
func resolveStatus(err error, code []int, fallback int) int {
	if len(code) > 0 && code[0] > 0 {
		return code[0]
	}
	if nvterr.IsKnown(err) {
		return nvterr.StatusOf(err)
	}
	return fallback
}
