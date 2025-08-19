package benchmark

import (
	"bytes"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/JonJenson-MFIn/project-management-system-api/graph/generated"
	"github.com/JonJenson-MFIn/project-management-system-api/graph/resolvers"
)

func BenchmarkQueryResolvers(b *testing.B) {
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

	body, _ := json.Marshal(GraphQLRequest{
		Query: `query { employees { id name } }`,
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		srv.Config.Handler.ServeHTTP(w, req)
	}
}

func BenchmarkMutationResolvers(b *testing.B) {
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

	body, _ := json.Marshal(GraphQLRequest{
		Query: `mutation { createEmployee(input: { name: "Test User", email: "test@example.com" }) { id name } }`,
	})

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		srv.Config.Handler.ServeHTTP(w, req)
	}
}
