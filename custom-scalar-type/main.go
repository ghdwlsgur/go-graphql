package main

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
)

type CustomID struct {
	value string
}

func (id *CustomID) String() string {
	return id.value
}

func NewCustomID(v string) *CustomID {
	return &CustomID{value: v}
}

var CustomScalarType = graphql.NewScalar(graphql.ScalarConfig{
	Name:        "CustomScalarType",
	Description: "The `CustomScalarType` scalar type represents an ID Object.",

	// Serialize serializes `CustomID` to string.
	Serialize: func(value interface{}) interface{} {
		switch value := value.(type) {
		case CustomID:
			return value.String()
		case *CustomID:
			v := *value
			return v.String()
		default:
			return nil
		}
	},
	// ParseValue parses GraphQL variables from `string` to `CustomID`.
	ParseValue: func(value interface{}) interface{} {
		switch value := value.(type) {
		case string:
			return NewCustomID(value)
		case *string:
			return NewCustomID(*value)
		default:
			return nil
		}
	},
	// ParseLiteral parses GraphQL AST value to `CustomID`.
	ParseLiteral: func(valueAST ast.Value) interface{} {
		switch valueAST := valueAST.(type) {
		case *ast.StringValue:
			return NewCustomID(valueAST.Value)
		default:
			return nil
		}
	},
})
