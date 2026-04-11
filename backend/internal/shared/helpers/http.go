package helpers

func RespondWithError(c interface {
	JSON(int, any)
}, statusCode int, err error) {
	c.JSON(statusCode, map[string]string{"error": err.Error()})
}

func RespondWithErrors(c interface {
	JSON(int, any)
}, statusCode int, errs []error) {
	messages := make([]string, 0, len(errs))
	for _, err := range errs {
		if err != nil {
			messages = append(messages, err.Error())
		}
	}

	c.JSON(statusCode, map[string]any{"errors": messages})
}
