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

package modrerankertransformers

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/weaviate/weaviate/entities/modulecapabilities"
	"github.com/weaviate/weaviate/entities/moduletools"
	rerankeradditional "github.com/weaviate/weaviate/modules/reranker-transformers/additional"
	rerankeradditionalrank "github.com/weaviate/weaviate/modules/reranker-transformers/additional/rank"
	client "github.com/weaviate/weaviate/modules/reranker-transformers/clients"
	"github.com/weaviate/weaviate/modules/reranker-transformers/ent"
)

const Name = "reranker-transformers"

func New() *ReRankerModule {
	return &ReRankerModule{}
}

type ReRankerModule struct {
	reranker                     ReRankerClient
	additionalPropertiesProvider modulecapabilities.AdditionalProperties
}

type ReRankerClient interface {
	Rank(ctx context.Context, property string, query string) (*ent.RankResult, error)
	MetaInfo() (map[string]interface{}, error)
}

func (m *ReRankerModule) Name() string {
	return Name
}

func (m *ReRankerModule) Type() modulecapabilities.ModuleType {
	return modulecapabilities.Text2TextReranker
}

func (m *ReRankerModule) Init(ctx context.Context,
	params moduletools.ModuleInitParams,
) error {
	if err := m.initAdditional(ctx, params.GetLogger()); err != nil {
		return errors.Wrap(err, "init re encoder")
	}

	return nil
}

func (m *ReRankerModule) initAdditional(ctx context.Context,
	logger logrus.FieldLogger,
) error {
	uri := os.Getenv("RERANKER_INFERENCE_API")
	if uri == "" {
		return nil
	}

	client := client.New(uri, logger)

	m.reranker = client
	if err := client.WaitForStartup(ctx, 1*time.Second); err != nil {
		return errors.Wrap(err, "init remote sum module")
	}

	rerankerProvider := rerankeradditionalrank.New(m.reranker)
	m.additionalPropertiesProvider = rerankeradditional.New(rerankerProvider)
	return nil
}

func (m *ReRankerModule) MetaInfo() (map[string]interface{}, error) {
	return m.reranker.MetaInfo()
}

func (m *ReRankerModule) RootHandler() http.Handler {
	// TODO: remove once this is a capability interface
	return nil
}

func (m *ReRankerModule) AdditionalProperties() map[string]modulecapabilities.AdditionalProperty {
	return m.additionalPropertiesProvider.AdditionalProperties()
}

// verify we implement the modules.Module interface
var (
	_ = modulecapabilities.Module(New())
	_ = modulecapabilities.AdditionalProperties(New())
	_ = modulecapabilities.MetaProvider(New())
)
