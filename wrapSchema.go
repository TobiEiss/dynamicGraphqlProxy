package dynamicGraphqlProxy

import (
	"context"
	"fmt"
	"net/http"

	"github.com/TobiEiss/schemaToRest"
	"github.com/graphql-go/handler"
	"github.com/labstack/echo"
)

// ContextKey is the type to define keys for context-values
type ContextKey string

// EchoContext is a key to get echoContext from request-context
const EchoContext ContextKey = "echoContext"

func wrapSchema(delination Delineation) echo.HandlerFunc {
	var setEchoContext = func(ctx echo.Context) *http.Request {
		return ctx.Request().WithContext(context.WithValue(ctx.Request().Context(), EchoContext, ctx))
	}

	switch delination.DelineationType {
	case Rest:
		return func(ctx echo.Context) error {
			ctx.SetRequest(setEchoContext(ctx))
			return schemaToRest.WrapSchema(delination.Schema)(ctx)
		}
	case Graphql:
		return func(ctx echo.Context) error {
			ctx.SetRequest(setEchoContext(ctx))
			handler.New(&handler.Config{
				Schema:   delination.Schema,
				Pretty:   true,
				GraphiQL: true,
			}).ServeHTTP(ctx.Response(), ctx.Request())
			return nil
		}
	}
	return func(ctx echo.Context) error {
		return ctx.JSON(http.StatusInternalServerError, fmt.Sprintf("Can't find delinationType: %s", delination.DelineationType))
	}
}
