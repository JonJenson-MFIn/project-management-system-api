package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/JonJenson-MFIn/project-management-system-api/graph/generated"
	"github.com/JonJenson-MFIn/project-management-system-api/graph/model"
	"github.com/JonJenson-MFIn/project-management-system-api/graph/resolvers"
)

func BenchmarkAuthDirective(b *testing.B) {
	srv := httptest.NewServer(
		handler.NewDefaultServer(
			generated.NewExecutableSchema(
				generated.Config{
					Resolvers: &resolvers.Resolver{},
					Directives: generated.DirectiveRoot{
						Auth: func(ctx context.Context, obj interface{}, next graphql.Resolver, role model.Role) (interface{}, error) {
							if ctx.Value("role") != role {
								return nil, fmt.Errorf("unauthorized")
							}
							return next(ctx)
						},
					},
				},
			),
		),
	)
	defer srv.Close()

	body, _ := json.Marshal(GraphQLRequest{
		Query: `mutation { deleteEmployee(id: 1) }`,
	})

	b.ResetTimer() 

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		req = req.WithContext(context.WithValue(req.Context(), "role", model.RoleAdmin)) 

		w := httptest.NewRecorder()
		srv.Config.Handler.ServeHTTP(w, req)
	}
}
