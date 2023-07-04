package utils

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Dur("duration", duration).
		Msg("received a gRPC request")

	return result, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (rec *ResponseRecorder) WriteHeader(code int) {
	rec.StatusCode = code
	rec.ResponseWriter.WriteHeader(code)
}

func (rec *ResponseRecorder) Write(b []byte) (int, error) {
	rec.Body = b
	return rec.ResponseWriter.Write(b)
}

func HttpLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		rec := &ResponseRecorder{
			ResponseWriter: c.Writer,
			StatusCode:     http.StatusOK,
		}

		startTime := time.Now()
		duration := time.Since(startTime)

		logger := log.Info().
			Str("protocol", "http").
			Str("method", c.Request.Method).
			Str("path", c.Request.URL.Path).
			Int("status_code", rec.StatusCode).
			Str("status_text", http.StatusText(rec.StatusCode)).
			Dur("duration", duration)

		if rec.StatusCode >= 500 {
			logger = logger.Bytes("body", rec.Body)
		}

		logger.Msg("received an HTTP request")
	}
}
