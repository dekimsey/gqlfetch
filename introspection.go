package main

import (
	"encoding/json"
	"log"

	"github.com/vektah/gqlparser/ast"
)

type introspectionResult struct {
	Errors graphqlErrs `json:"errors"`
	Data   struct {
		Schema introspectionSchema `json:"__schema"`
	} `json:"data"`
}

type introspectionSchema struct {
	QueryType    ast.Definition                     `json:"queryType"`
	MutationType ast.Definition                     `json:"mutationType"`
	Types        []introspectionTypeDefinition      `json:"types"`
	Directives   []introspectionDirectiveDefinition `json:"directives"`
}

type introspectionTypeDefinition struct {
	Kind        ast.DefinitionKind `json:"kind"`
	Name        string             `json:"name"`
	Description string             `json:"description"`
	Fields      []struct {
		Name              string            `json:"name"`
		Description       string            `json:"description"`
		Args              []interface{}     `json:"args"`
		Type              *introspectedType `json:"type"`
		IsDeprecated      bool              `json:"isDeprecated"`
		DeprecationReason interface{}       `json:"deprecationReason"`
	} `json:"fields"`
	InputFields   []introspectionInputField `json:"inputFields"`
	Interfaces    []ast.Definition          `json:"interfaces"`
	EnumValues    json.RawMessage           `json:"enumValues"`
	PossibleTypes json.RawMessage           `json:"possibleTypes"`
}

type introspectionDirectiveDefinition struct {
	Name        string                  `json:"name"`
	Description string                  `json:"description"`
	Locations   []ast.DirectiveLocation `json:"locations"`
	Args        []struct {
		Name         string            `json:"name"`
		Description  string            `json:"description"`
		Type         *introspectedType `json:"type"`
		DefaultValue interface{}       `json:"defaultValue"`
	} `json:"args"`
}

type introspectionInputField struct {
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	Type         introspectedType `json:"type"`
	DefaultValue interface{}      `json:"defaultValue"`
}

type introspectedType struct {
	Kind   introspectionTypeKind `json:"kind"`
	Name   *string               `json:"name"`
	OfType *introspectedType     `json:"ofType"`
}

type introspectionTypeKind string

const (
	NON_NULL introspectionTypeKind = "NON_NULL"
	LIST     introspectionTypeKind = "LIST"
)

func introspectionTypeToAstType(typ *introspectedType) *ast.Type {
	var res ast.Type
	if typ.OfType == nil {
		res.NamedType = *typ.Name
		return &res
	}

	switch typ.Kind {
	case NON_NULL:
		res.NonNull = true
		res.Elem = introspectionTypeToAstType(typ.OfType)
		return &res
	case LIST:
		res.Elem = introspectionTypeToAstType(typ.OfType)
		return &res
	default:
		log.Fatalf("type kind unknown: %s", typ.Kind)
    return nil
	}
}
