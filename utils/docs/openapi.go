package docs

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3gen"
	"gopkg.in/yaml.v3"
)

type ApiWrapper struct {
	Path                 string          `json:"path"`
	Method               string          `json:"method"`
	Name                 string          `json:"name"`
	Description          string          `json:"description"`
	Tags                 []string        `json:"tags,omitempty"`
	RequestBody          *ApiRequestBody `json:"requestBody,omitempty"`
	Response             *ApiResponse    `json:"response,omitempty"`
	Parameters           []ApiParameter  `json:"parameters,omitempty"`
	UnAuthenticated      bool            `json:"unauthenticated,omitempty"`
	ProjectIDNotRequired bool            `json:"projectIdNotRequired,omitempty"`
}

type ApiParameter struct {
	Name        string `json:"name"`
	In          string `json:"in"` // e.g., "query", "path", "header", "cookie"
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required,omitempty"`
}

type ApiRequestBody struct {
	Description string      `json:"description,omitempty"`
	Content     interface{} `json:"content,omitempty"`
}

// Response represents the response structure for an API
type ApiResponse struct {
	Description string      `json:"description,omitempty"`
	Content     interface{} `json:"content,omitempty"`
}

var apis = map[string][]ApiWrapper{}

func RegisterApi(api ApiWrapper) {
	apis[api.Path] = append(apis[api.Path], api)
}

func GenerateOpenAPI() ([]byte, error) {
	if len(apis) == 0 {
		return nil, errors.New("no APIs registered")
	}

	swagger := &openapi3.T{
		OpenAPI: "3.0.3",
		Info: &openapi3.Info{
			Title:       "Go IAM APIs",
			Description: intoDescription,
			Version:     "1.0.0",
		},
	}

	paths := []openapi3.NewPathsOption{}
	for _, api := range apis {
		pathItem, err := generateDocsForApi(api)
		if err != nil {
			return nil, fmt.Errorf("error generating docs for API %s: %w", api[0].Name, err)
		}
		// replace fiber path params like ":id" with OpenAPI style "{id}"
		path := extractPathParams(api[0].Path)
		paths = append(paths, openapi3.WithPath(path, pathItem))
	}

	swagger.Paths = openapi3.NewPaths(paths...)
	com := openapi3.NewComponents()
	com.SecuritySchemes = openapi3.SecuritySchemes{
		"BearerAuth": &openapi3.SecuritySchemeRef{
			Value: &openapi3.SecurityScheme{
				Type:         "http",
				Description:  "Auth is required if and only if go IAM is configured to use auth providers. By default only `/me/v1` requires authentication.",
				Scheme:       "bearer",
				BearerFormat: "JWT",
			},
		},
	}
	swagger.Components = &com

	out, err := yaml.Marshal(swagger)
	if err != nil {
		return nil, fmt.Errorf("error marshaling OpenAPI document: %w", err)
	}
	return out, nil
}

func extractPathParams(path string) string {
	var params []string
	for _, part := range strings.Split(path, "/") {
		if len(part) > 0 && part[0] == ':' {
			params = append(params, fmt.Sprintf("{%s}", part[1:]))
		} else {
			params = append(params, part)
		}
	}
	return strings.Join(params, "/")
}

func generateDocsForApi(apiss []ApiWrapper) (*openapi3.PathItem, error) {
	if len(apis) == 0 {
		return nil, errors.New("API path and method must be specified")
	}

	doc := &openapi3.PathItem{}
	for _, api := range apiss {
		operation := &openapi3.Operation{
			Summary:     api.Name,
			Description: api.Description,
			Tags:        api.Tags,
		}

		// Add request body if provided
		if api.RequestBody != nil {
			ref, err := generateDocsForRequestBody(*api.RequestBody)
			if err != nil {
				return nil, fmt.Errorf("error generating request body for API %s: %w", api.Name, err)
			}
			operation.RequestBody = ref
		}

		// Add response if provided
		if api.Response != nil {
			ref, err := generateDocsForResponse(*api.Response)
			if err != nil {
				return nil, fmt.Errorf("error generating response for API %s: %w", api.Name, err)
			}
			operation.Responses = ref
		}

		// Add parameters if provided
		if len(api.Parameters) > 0 {
			params, err := generateDocsForParameters(api.Parameters)
			if err != nil {
				return nil, fmt.Errorf("error generating parameters for API %s: %w", api.Name, err)
			}
			operation.Parameters = params
		}

		if api.UnAuthenticated {
			operation.Security = nil // Unauthenticated APIs do not require security
		} else {
			operation.Security = &openapi3.SecurityRequirements{
				{
					"BearerAuth": []string{},
				},
			}
		}

		if !api.ProjectIDNotRequired {
			operation.Parameters = append(operation.Parameters, &openapi3.ParameterRef{
				Value: &openapi3.Parameter{
					In:          "header",
					Name:        "X-Project-Ids",
					Description: "Comma-separated list of project IDs. If not provided, the API will use the default project ID configured in the server.",
					Required:    true,
					Schema:      openapi3.NewStringSchema().NewRef(),
				},
			})
		}

		switch api.Method {
		case "GET":
			doc.Get = operation
		case "POST":
			doc.Post = operation
		case "PUT":
			doc.Put = operation
		case "DELETE":
			doc.Delete = operation
		default:
			return nil, fmt.Errorf("unsupported HTTP method %s", api.Method)
		}
	}

	return doc, nil
}

func generateDocsForRequestBody(requestBody ApiRequestBody) (*openapi3.RequestBodyRef, error) {
	if requestBody.Content == nil {
		return nil, errors.New("request body content must be provided")
	}
	bodyRef, err := openapi3gen.NewSchemaRefForValue(requestBody.Content, nil)
	if err != nil {
		return nil, fmt.Errorf("error converting sdk to json schema %w", err)
	}
	return &openapi3.RequestBodyRef{
		Value: &openapi3.RequestBody{
			Description: requestBody.Description,
			Content: map[string]*openapi3.MediaType{
				"application/json": {
					Schema: bodyRef,
				},
			},
		},
	}, nil
}

func generateDocsForResponse(response ApiResponse) (*openapi3.Responses, error) {
	if response.Content == nil {
		return nil, errors.New("response content must be provided")
	}
	resRef, err := openapi3gen.NewSchemaRefForValue(response.Content, nil)
	if err != nil {
		return nil, fmt.Errorf("error converting sdk to json schema %w", err)
	}
	defaultRes := openapi3.WithName("default", &openapi3.Response{
		Description: &response.Description,
		Content: map[string]*openapi3.MediaType{
			"application/json": {
				Schema: resRef,
			},
		},
	})
	return openapi3.NewResponses(defaultRes), nil
}

func generateDocsForParameters(parameters []ApiParameter) ([]*openapi3.ParameterRef, error) {
	if len(parameters) == 0 {
		return nil, nil
	}

	var paramRefs []*openapi3.ParameterRef
	for _, param := range parameters {
		paramRef := &openapi3.ParameterRef{
			Value: &openapi3.Parameter{
				Name:        param.Name,
				In:          param.In,
				Description: param.Description,
				Required:    param.Required,
			},
		}
		paramRefs = append(paramRefs, paramRef)
	}
	return paramRefs, nil
}

func CreateOpenApiDoc(fileName string) error {
	doc, err := GenerateOpenAPI()
	if err != nil {
		return fmt.Errorf("error generating OpenAPI document: %w", err)
	}
	return os.WriteFile(fileName, doc, 0644)
}
