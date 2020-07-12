package helper

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"

	"github.com/gin-gonic/gin/binding"
	zhongwen "github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/go-playground/validator/v10/translations/zh"
)

type defaultValidator struct {
	once     sync.Once
	validate *validator.Validate
}

var ValidatorV10 binding.StructValidator = &defaultValidator{}

func (v *defaultValidator) ValidateStruct(obj interface{}) error {

	if kindOfData(obj) == reflect.Struct {

		v.lazyinit()

		if err := v.validate.Struct(obj); err != nil {
			errs := err.(validator.ValidationErrors)
			//自动翻译
			return errors.New(Translate(errs))
		}
	}

	return nil
}

func (v *defaultValidator) Engine() interface{} {
	v.lazyinit()
	return v.validate
}

func Translate(errs validator.ValidationErrors) string {
	var errList []string
	for _, e := range errs {
		// can translate each error one at a time.
		errList = append(errList, e.Translate(trans))
	}
	return strings.Join(errList, "|")
}

func (v *defaultValidator) lazyinit() {
	v.once.Do(func() {
		v.validate = validator.New()
		v.validate.SetTagName("validate")

		zhs := zhongwen.New()
		uni := ut.New(zhs, zhs)
		trans, _ := uni.GetTranslator("zh")
		err := zh.RegisterDefaultTranslations(v.validate, trans)
		if err != nil {
			panic(err)
		}

		// add any custom validations etc. here
		// 自定义验证方法 todo more
		err = v.validate.RegisterValidation("checkMobile", CheckPassword)
		if err != nil {
			fmt.Print(err.Error())
		}
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
