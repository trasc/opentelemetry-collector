// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package configmapprovider // import "go.opentelemetry.io/collector/config/configmapprovider"

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v2"

	"go.opentelemetry.io/collector/config"
)

type fileMapProvider struct {
	fileName string
}

// NewFile returns a new Provider that reads the configuration from the given file.
func NewFile(fileName string) Provider {
	return &fileMapProvider{
		fileName: fileName,
	}
}

func (fmp *fileMapProvider) Retrieve(_ context.Context, _ func(*ChangeEvent)) (Retrieved, error) {
	if fmp.fileName == "" {
		return nil, errors.New("config file not specified")
	}

	content, err := ioutil.ReadFile(fmp.fileName)
	if err != nil {
		return nil, fmt.Errorf("unable to read the file %v: %w", fmp.fileName, err)
	}

	var data map[string]interface{}
	if err = yaml.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("unable to parse yaml: %w", err)
	}

	return &simpleRetrieved{confMap: config.NewMapFromStringMap(data)}, nil
}

func (*fileMapProvider) Shutdown(context.Context) error {
	return nil
}
