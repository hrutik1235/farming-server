package middleware

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/hrutik1235/farming-server/utils"
	"github.com/hrutik1235/farming-server/validators"
)

type ValidateConfig struct {
	// ContinueValidate indicates whether to continue validation even if the first validation fails.
	ContinueValidate bool // defaults to true
}

func ConfigureVC() *ValidateConfig {
	return &ValidateConfig{
		ContinueValidate: true,
	}
}

func (v *ValidateConfig) SetContinueValidate(val bool) *ValidateConfig {
	v.ContinueValidate = val
	return v
}

func validationErroHandler(err error, payload any, config *ValidateConfig) []string {
	errorsMap := []string{}
	if errs, ok := err.(validator.ValidationErrors); ok {
		tType := reflect.TypeOf(payload)
		if tType.Kind() == reflect.Ptr {
			tType = tType.Elem()
		}
		for _, e := range errs {
			field := e.Field()
			// array of fields like "field[0]" should be handled
			if idx := strings.Index(field, "["); idx != -1 {
				field = field[:idx]
			}
			_field, _ := tType.FieldByName(field)

			switch e.Tag() {
			case "required":
				errorsMap = append(errorsMap, fmt.Sprintf("%s is required", _field.Tag.Get("name")))
			default:
				customMessage := _field.Tag.Get("message")
				if customMessage != "" {
					errorsMap = append(errorsMap, customMessage)
				} else {
					fmt.Printf("Field: %s, Tag: %s, Value: %s\n", field, _field.Tag.Get("name"), e.Value())
					errorsMap = append(errorsMap, fmt.Sprintf("%s is invalid", _field.Tag.Get("name")))
				}
			}
			if !config.ContinueValidate && len(errorsMap) > 0 {
				break
			}
		}
	}
	return errorsMap
}

func ValidateRequest[BodyT any, QueryT any, ParamT any](cfg ...*ValidateConfig) gin.HandlerFunc {
	return func(c *gin.Context) {
		var config *ValidateConfig

		if len(cfg) > 0 {
			config = cfg[0]
		} else {
			config = ConfigureVC()
		}

		validate := validator.New()
		validate.RegisterValidation("indianphone", validators.IndianPhoneValidator)

		// Validate Body
		var body BodyT
		if hasFields[BodyT]() && c.Request.Body != nil {
			buf, _ := io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(buf))

			if len(buf) > 0 {
				if err := c.ShouldBindJSON(&body); err != nil {
					errs := validationErroHandler(err, body, config)
					if len(errs) > 0 {
						c.JSON(http.StatusBadRequest, utils.NewHttpError(c, errs, http.StatusBadRequest))
						c.Abort()
					}
					return
				}
				if err := validate.Struct(body); err != nil {
					errs := validationErroHandler(err, body, config)
					if len(errs) > 0 {
						c.JSON(http.StatusBadRequest, utils.NewHttpError(c, errs, http.StatusBadRequest))
						c.Abort()
					}
					return
				}
				c.Set("body", body)
			}
		}

		// Validate Query
		var query QueryT
		if hasFields[QueryT]() {
			fmt.Printf("Query Params: %+v\n", c.Request.URL.Query())
			if err := c.ShouldBindQuery(&query); err != nil {
				fmt.Printf("Error binding query: %v\n", err)
				errs := validationErroHandler(err, query, config)
				if len(errs) > 0 {
					c.JSON(http.StatusBadRequest, utils.NewHttpError(c, errs, http.StatusBadRequest))
					c.Abort()
				}
				return
			}

			if err := validate.Struct(query); err != nil {
				errs := validationErroHandler(err, query, config)
				if len(errs) > 0 {
					c.JSON(http.StatusBadRequest, utils.NewHttpError(c, errs, http.StatusBadRequest))
					c.Abort()
				}
				return
			}
			c.Set("query", query)
		}

		// Validate Params
		var params ParamT
		if hasFields[ParamT]() {
			if err := c.ShouldBindUri(&params); err != nil {
				errs := validationErroHandler(err, params, config)
				if len(errs) > 0 {
					c.JSON(http.StatusBadRequest, utils.NewHttpError(c, errs, http.StatusBadRequest))
					c.Abort()
				}
				return
			}

			if err := validate.Struct(params); err != nil {
				errs := validationErroHandler(err, params, config)
				if len(errs) > 0 {
					c.JSON(http.StatusBadRequest, utils.NewHttpError(c, errs, http.StatusBadRequest))
					c.Abort()
				}
				return
			}
			c.Set("params", params)
		}
		val, _ := c.Get("params")
		fmt.Printf("Request Params: %+v\n", val)

		c.Next()
	}
}

func hasFields[T any]() bool {
	t := reflect.TypeOf((*T)(nil)).Elem()
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return t.Kind() == reflect.Struct && t.NumField() > 0
}
