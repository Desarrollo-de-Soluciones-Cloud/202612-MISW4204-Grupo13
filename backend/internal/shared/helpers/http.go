package helpers

func RespondWithError(c interface {
	JSON(int, any)
}, statusCode int, err error) {
	c.JSON(statusCode, map[string]string{"error": err.Error()})
}
