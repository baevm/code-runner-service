package logger

import "go.uber.org/zap"

func New() (*zap.SugaredLogger, error) {
	logger, err := zap.NewProduction()
	defer logger.Sync()

	if err != nil {
		return nil, err
	}

	sugared := logger.Sugar()

	

	return sugared, nil
}
