package delivery

import (
	"errors"
	"strconv"
)

var ErrInvalidUserID = errors.New("invalid user id")

func parseUserIDParam(rawID string) (uint, error) {
	id, err := strconv.ParseUint(rawID, 10, 32)
	if err != nil || id == 0 {
		return 0, ErrInvalidUserID
	}

	return uint(id), nil
}
