// Copyright 2019 The Terraformer Authors.
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

package heroku

import (
	"errors"
	"os"

	"github.com/GoogleCloudPlatform/terraformer/terraform_utils"
)

const herokuProviderVersion = "~> 2.2.1"

type HerokuProvider struct {
	terraform_utils.Provider
	email  string
	apiKey string
}

func (p *HerokuProvider) Init(args []string) error {
	if os.Getenv("HEROKU_EMAIL") == "" {
		return errors.New("set HEROKU_EMAIL env var")
	}
	p.email = os.Getenv("HEROKU_EMAIL")

	if os.Getenv("HEROKU_API_KEY") == "" {
		return errors.New("set HEROKU_API_KEY env var")
	}
	p.apiKey = os.Getenv("HEROKU_API_KEY")

	return nil
}

func (p *HerokuProvider) GetName() string {
	return "heroku"
}

func (p *HerokuProvider) GetProviderData(arg ...string) map[string]interface{} {
	return map[string]interface{}{
		"provider": map[string]interface{}{
			"heroku": map[string]interface{}{
				"version": herokuProviderVersion,
				"email":   p.email,
				"api_key": p.apiKey,
			},
		},
	}
}

func (HerokuProvider) GetResourceConnections() map[string]map[string][]string {
	return map[string]map[string][]string{}
}

func (p *HerokuProvider) GetSupportedService() map[string]terraform_utils.ServiceGenerator {
	return map[string]terraform_utils.ServiceGenerator{
		"app":   &AppGenerator{},
		"addon": &AddonGenerator{},
	}
}

func (p *HerokuProvider) InitService(serviceName string) error {
	var isSupported bool
	if _, isSupported = p.GetSupportedService()[serviceName]; !isSupported {
		return errors.New("heroku: " + serviceName + " not supported service")
	}
	p.Service = p.GetSupportedService()[serviceName]
	p.Service.SetName(serviceName)
	p.Service.SetProviderName(p.GetName())
	p.Service.SetArgs(map[string]interface{}{
		"email":   p.email,
		"api_key": p.apiKey,
	})
	return nil
}
