package directives

import (
	"context"
	"fmt"

	"github.com/99designs/gqlgen/graphql"
	"github.com/JonJenson-MFIn/project-management-system-api/middleware"
	"github.com/JonJenson-MFIn/project-management-system-api/graph/model"
)

func AuthDirective(ctx context.Context, obj interface{} , next graphql.Resolver, role model.Role) (interface{}, error) {
	user := middleware.GetUserFromContext(ctx)
	if user == nil {
		return nil, fmt.Errorf("unauthenticated")
	}

	if user.Role != role {
		return nil, fmt.Errorf("forbidden: requires %s, got %s", role, user.Role)
	}

	return next(ctx)
}
