package helpers

import (
	sharedErrors "backend/internal/shared/errors"
	"strconv"
)

func ParseResourceID(rawID string) (uint, error) {
	id, err := strconv.ParseUint(rawID, 10, 32)
	if err != nil || id == 0 {
		return 0, sharedErrors.ErrInvalidResourceID
	}

	return uint(id), nil
}
