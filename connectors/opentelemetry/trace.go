package opentelemetry

import (
	"fmt"

	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func Wrapf(span trace.Span, err error, template string, args ...interface{}) error {
	if err == nil {
		return nil
	}

	template = fmt.Sprintf("%s; traceID: %s", template, span.SpanContext().TraceID().String())
	err = errors.Wrapf(err, template, args...)
	span.SetStatus(codes.Error, err.Error())

	return err
}
