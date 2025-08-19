package tdd

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/JonJenson-MFIn/project-management-system-api/graph/generated"
	"github.com/JonJenson-MFIn/project-management-system-api/graph/resolvers"
)

func TestMiddlewareIntegration(t *testing.T) {
	tests := []struct {
		name        string
		query       string
		contentType string
		wantErr     bool
	}{
		{
			name:        "Valid JSON content type should work",
			query:       `query { employees { id } }`,
			contentType: "application/json",
			wantErr:     false,
		},
		{
			name:        "Invalid content type should fail",
			query:       `query { employees { id } }`,
			contentType: "text/plain",
			wantErr:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(
				handler.NewDefaultServer(
					generated.NewExecutableSchema(
						generated.Config{
							Resolvers: &resolvers.Resolver{},
						},
					),
				),
			)
			defer srv.Close()

			body, _ := json.Marshal(GraphQLRequest{Query: tt.query})
			req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", tt.contentType)

			w := httptest.NewRecorder()
			srv.Config.Handler.ServeHTTP(w, req)

			if (w.Code != 200) != tt.wantErr {
				t.Errorf("got status %d, wantErr=%v", w.Code, tt.wantErr)
			}
		})
	}
}
