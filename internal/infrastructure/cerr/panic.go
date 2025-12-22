package cerr

import (
	"fmt"
	"my-go-playground/internal/infrastructure/logging"

	"github.com/rs/zerolog"
)

func HandlePanic(r any, stack []byte, logger *zerolog.Logger) {
	if r == nil {
		return
	}
	err, ok := r.(error)
	if !ok {
		err = fmt.Errorf("%v", r)
	}
	logger.Fatal().
		Bytes(logging.FieldStack, stack).
		Err(err).
		Msgf("panic")
}
