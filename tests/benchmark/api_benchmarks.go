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

func BenchmarkFullAPIEndpoint(b *testing.B) {
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

	// Test multiple operations in sequence
	queries := []string{
		`query { employees { id name email } }`,
		`query { projects { id name description } }`,
		`mutation { createEmployee(input: { name: "Benchmark User", email: "benchmark@example.com" }) { id name } }`,
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		for _, query := range queries {
			body, _ := json.Marshal(GraphQLRequest{Query: query})
			req := httptest.NewRequest("POST", "/", bytes.NewBuffer(body))
			req.Header.Set("Content-Type", "application/json")

			w := httptest.NewRecorder()
			srv.Config.Handler.ServeHTTP(w, req)
		}
	}
}
