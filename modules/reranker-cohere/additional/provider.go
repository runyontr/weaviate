//                           _       _
// __      _____  __ ___   ___  __ _| |_ ___
// \ \ /\ / / _ \/ _` \ \ / / |/ _` | __/ _ \
//  \ V  V /  __/ (_| |\ V /| | (_| | ||  __/
//   \_/\_/ \___|\__,_| \_/ |_|\__,_|\__\___|
//
//  Copyright © 2016 - 2023 Weaviate B.V. All rights reserved.
//
//  CONTACT: hello@weaviate.io
//

package additional

import (
	"context"

	"github.com/tailor-inc/graphql"
	"github.com/tailor-inc/graphql/language/ast"
	"github.com/weaviate/weaviate/entities/modulecapabilities"
	"github.com/weaviate/weaviate/entities/moduletools"
	"github.com/weaviate/weaviate/entities/search"
)

type AdditionalProperty interface {
	AdditionalPropertyFn(ctx context.Context,
		in []search.Result, params interface{}, limit *int,
		argumentModuleParams map[string]interface{}, cfg moduletools.ClassConfig) ([]search.Result, error)
	ExtractAdditionalFn(param []*ast.Argument) interface{}
	AdditionalPropertyDefaultValue() interface{}
	AdditionalFieldFn(classname string) *graphql.Field
}

type GraphQLAdditionalArgumentsProvider struct {
	CrossRankerProvider AdditionalProperty
}

func New(CrossRankerProvider AdditionalProperty) *GraphQLAdditionalArgumentsProvider {
	return &GraphQLAdditionalArgumentsProvider{CrossRankerProvider}
}

func (p *GraphQLAdditionalArgumentsProvider) AdditionalProperties() map[string]modulecapabilities.AdditionalProperty {
	additionalProperties := map[string]modulecapabilities.AdditionalProperty{}
	additionalProperties["rerank"] = p.getCrossRanker()
	return additionalProperties
}

func (p *GraphQLAdditionalArgumentsProvider) getCrossRanker() modulecapabilities.AdditionalProperty {
	return modulecapabilities.AdditionalProperty{
		GraphQLNames:           []string{"rerank"},
		GraphQLFieldFunction:   p.CrossRankerProvider.AdditionalFieldFn,
		GraphQLExtractFunction: p.CrossRankerProvider.ExtractAdditionalFn,
		SearchFunctions: modulecapabilities.AdditionalSearch{
			ExploreGet:  p.CrossRankerProvider.AdditionalPropertyFn,
			ExploreList: p.CrossRankerProvider.AdditionalPropertyFn,
		},
	}
}
