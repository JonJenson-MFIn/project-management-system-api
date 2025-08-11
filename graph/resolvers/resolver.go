package resolvers
//go:generate go run github.com/99designs/gqlgen generate
import "gorm.io/gorm"

type Resolver struct {
	DB *gorm.DB
}
