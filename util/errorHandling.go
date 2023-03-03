package util

import "go.uber.org/zap"

func HandleErr(err error, isFatal bool, logger *zap.SugaredLogger) {
    if err != nil {
        if isFatal {
            logger.Fatalw("Catastrophic failure", "error", err.Error())
        } else {
            logger.Infow("failure", "error", err.Error())
        }
    }
}
