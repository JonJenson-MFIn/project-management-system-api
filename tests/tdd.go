package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql"
	"github.com/JonJenson-MFIn/project-management-system-api/graph/generated"
	"github.com/JonJenson-MFIn/project-management-system-api/graph/model"
	"github.com/JonJenson-MFIn/project-management-system-api/graph/resolvers"
)

type GraphQLRequest struct {
	Query string `json:"query"`
}

func TestAuthDirective(t *testing.T) {
	tests := []struct {
		name    string
		role    string
		query   string
		wantErr bool
	}{
		{"ADMIN can delete employee", "ADMIN", `mutation { deleteEmployee(id: "1") }`, false},
		{"USER cannot delete employee", "USER", `mutation { deleteEmployee(id: "1") }`, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
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

			body, _ := json.Marshal(GraphQLRequest{Query: tt.query})
			req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")
			req = req.WithContext(context.WithValue(req.Context(), "role", tt.role))

			w := httptest.NewRecorder()
			srv.Config.Handler.ServeHTTP(w, req)

			if (w.Code != 200) != tt.wantErr {
				t.Errorf("got status %d, wantErr=%v", w.Code, tt.wantErr)
			}
		})
	}
}
