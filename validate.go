package helper

import (
	"fmt"
	"reflect"
	"regexp"
	"sync"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var ValidatorV10 binding.StructValidator = &defaultValidator{}

func init() {
	v, ok := ValidatorV10.Engine().(*validator.Validate)
	if ok {
		// 自定义验证方法 todo more
		err := v.RegisterValidation("checkMobile", CheckPassword)
		if err != nil {
			fmt.Print(err.Error())
		}
	}
}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			return error(err)
		}
	}

	return nil
}

func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("validate")

		// add any custom validations etc. here
	})
}

func kindOfData(data interface{}) reflect.Kind {

	value := reflect.ValueOf(data)
	valueType := value.Kind()

	if valueType == reflect.Ptr {
		valueType = value.Elem().Kind()
	}
	return valueType
}

func CheckPassword(fl validator.FieldLevel) bool {
	value := fl.Field().String()
	if ok, _ := regexp.MatchString("^$", value); !ok {
		return false
	}
	return true
}
