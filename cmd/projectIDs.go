package cmd

import (
	"fmt"
	"strconv"
)

// projectIDs converts the arguments of a command into Gitlab project IDs.
func projectIDs(args []string) ([]int, error) {
	ids := make([]int, 0, len(args))

	for _, arg := range args {
		id, err := strconv.Atoi(arg)
		if err != nil || id <= 0 {
			return nil, fmt.Errorf("%q is not a valid project ID", arg)
		}
		ids = append(ids, id)
	}

	return ids, nil
}
