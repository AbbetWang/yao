package table

import (
	"fmt"

	"github.com/yaoapp/gou"
)

// Exec execute the hook
func (hook *BeforeHookActionDSL) Exec(args []interface{}, sid string, global map[string]interface{}) ([]interface{}, error) {

	p, err := gou.ProcessOf(hook.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("%s %s", hook.String(), err.Error())
	}

	res, err := p.WithGlobal(global).WithSID(sid).Exec()
	if err != nil {
		return nil, fmt.Errorf("[%s] %s", hook.String(), err.Error())
	}

	newArgs, ok := res.([]interface{})
	if !ok {
		return nil, fmt.Errorf("%s return value is not an array", hook.String())
	}

	if len(newArgs) != len(args) {
		return nil, fmt.Errorf("%s return value is not correct. should: array[%d], got: array[%d]", hook.String(), len(args), len(newArgs))
	}

	return newArgs, nil
}

// Exec execute the hook
func (hook *AfterHookActionDSL) Exec(value interface{}, sid string, global map[string]interface{}) (interface{}, error) {

	args := []interface{}{}
	switch value.(type) {
	case []interface{}:
		args = value.([]interface{})
	default:
		args = append(args, value)
	}

	p, err := gou.ProcessOf(hook.String(), args...)
	if err != nil {
		return nil, fmt.Errorf("[%s] %s", hook.String(), err.Error())
	}

	res, err := p.WithGlobal(global).WithSID(sid).Exec()
	if err != nil {
		return nil, fmt.Errorf("[%s] %s", hook.String(), err.Error())
	}

	return res, nil
}

// String cast to string
func (hook *BeforeHookActionDSL) String() string {
	return string(*hook)
}

// String cast to string
func (hook *AfterHookActionDSL) String() string {
	return string(*hook)
}
