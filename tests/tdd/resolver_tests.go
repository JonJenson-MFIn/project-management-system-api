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

func TestQueryResolvers(t *testing.T) {
	tests := []struct {
		name    string
		query   string
		wantErr bool
	}{
		{
			name:    "Query employees should work",
			query:   `query { employees { id name } }`,
			wantErr: false,
		},
		{
			name:    "Query projects should work",
			query:   `query { projects { id name } }`,
			wantErr: false,
		},
		{
			name:    "Invalid query should fail",
			query:   `query { invalidField }`,
			wantErr: true,
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
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			srv.Config.Handler.ServeHTTP(w, req)

			if (w.Code != 200) != tt.wantErr {
				t.Errorf("got status %d, wantErr=%v", w.Code, tt.wantErr)
			}
		})
	}
}
