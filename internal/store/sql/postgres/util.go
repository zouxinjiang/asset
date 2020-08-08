package postgres

import (
	"fmt"
	"github.com/jinzhu/gorm"
	sentity "github.com/zouxinjiang/axes/internal/store/sql/entity"
	"math"
	"reflect"
	"strings"
)

func inStringSlice(s string, list []string) bool {
	for _, e := range list {
		if s == e {
			return true
		}
	}
	return false
}

func resolvedGormTagColumn(s string) string {
	columnTag := s
	tmp := strings.Split(s, ";")
	for _, v := range tmp {
		if strings.Index(strings.ToLower(strings.Trim(v, "\r\n\t ")), "column") == 0 {
			columnTag = v
		}
	}

	a := strings.Split(columnTag, ":")
	if len(a) >= 2 {
		return a[1]
	}
	return ""
}

func isReflectValueZero(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.String:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return math.Float64bits(v.Float()) == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Complex64, reflect.Complex128:
		c := v.Complex()
		return math.Float64bits(real(c)) == 0 && math.Float64bits(imag(c)) == 0
	case reflect.Array:
		for i := 0; i < v.Len(); i++ {
			if !isReflectValueZero(v.Index(i)) {
				return false
			}
		}
		return true
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return v.IsNil()
	case reflect.Struct:
		for i := 0; i < v.NumField(); i++ {
			if !isReflectValueZero(v.Field(i)) {
				return false
			}
		}
		return true
	default:
		panic(fmt.Sprintf("reflect.Value.IsZero %v", v.Kind()))
	}
}

// WrapQueryOrder 为query添加Order
func WrapQueryOrder(query *gorm.DB, filter sentity.Filter, keymap map[string]string) *gorm.DB {
	if query == nil {
		return query
	}
	order := genOrderSQL(filter, keymap)
	if order != "" {
		query = query.Order(order)
	}
	return query
}

func generateKeyMapFromModel(m interface{}) (keymap map[string]string) {
	keymap = map[string]string{}

	parseTag := func(gormTag string) (fieldName string, ok bool) {
		items := strings.Split(gormTag, ";")
		for _, v := range items {
			kv := strings.Split(v, ":")
			if len(kv) == 2 && strings.ToLower(strings.Trim(kv[0], "\r\n\t\v ")) == "column" {
				return strings.Trim(kv[1], "\r\t\v\n "), true
			}
		}
		return "", false
	}

	t := reflect.TypeOf(m)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() != reflect.Struct {
		return keymap
	}
	for i := 0; i < t.NumField(); i++ {
		tagValue := t.Field(i).Tag.Get("gorm")
		fieldName, ok := parseTag(tagValue)
		if !ok {
			continue
		}
		keymap[t.Field(i).Name] = fieldName
	}
	return keymap
}

func outerFields2innerFields(outer map[string]interface{}, keymap map[string]string) (inner map[string]interface{}) {
	for k, v := range outer {
		ki, ok := keymap[k]
		if ok {
			inner[ki] = v
		}
	}
	return
}

func outerField2innerField(outer string, keymap map[string]string) (inner string, ok bool) {
	inner, ok = keymap[outer]
	return
}

func wrapWithTableName(keymap map[string]string, table string) map[string]string {
	var res = map[string]string{}
	for k, v := range keymap {
		res[k] = `"` + table + `".` + v
	}
	return res
}

func simpleFilterGenSQL(cond sentity.Filter, keymap map[string]string) (string, []interface{}) {
	var andStr = ""
	var vals = []interface{}{}

	for _, item := range cond.AndCond {
		if ik, ok := outerField2innerField(item.OuterName, keymap); ok {
			if andStr != "" {
				andStr += " AND "
			}
			switch item.Cond {
			case sentity.IsNotNull, sentity.IsNull:
				andStr += fmt.Sprintf(" %s %s", ik, item.Cond)
			case sentity.In, sentity.NotIn:
				andStr += fmt.Sprintf(" %s %s (?)", ik, item.Cond)
				vals = append(vals, item.Value)
			case sentity.Like, sentity.NotLike:
				andStr += fmt.Sprintf(" %s %s ?", ik, item.Cond)
				vals = append(vals, fmt.Sprintf("%%%v%%", item.Value))
			default:
				andStr += fmt.Sprintf(" %s %s ?", ik, item.Cond)
				vals = append(vals, item.Value)
			}
		}
	}
	var orStr = ""
	if len(cond.OrCond) > 1 {
		for _, item := range cond.OrCond {
			if ik, ok := outerField2innerField(item.OuterName, keymap); ok {
				if orStr != "" {
					orStr += " OR "
				}
				switch item.Cond {
				case sentity.IsNotNull, sentity.IsNull:
					orStr += fmt.Sprintf(" %s %s", ik, item.Cond)
				case sentity.In, sentity.NotIn:
					orStr += fmt.Sprintf(" %s %s (?)", ik, item.Cond)
					vals = append(vals, item.Value)
				case sentity.Like, sentity.NotLike:
					orStr += fmt.Sprintf(" %s %s ?", ik, item.Cond)
					vals = append(vals, fmt.Sprintf("%%%v%%", item.Value))
				default:
					orStr += fmt.Sprintf(" %s %s ?", ik, item.Cond)
					vals = append(vals, item.Value)
				}
			}
		}
	}
	sqlStr := ""
	if andStr != "" {
		if orStr != "" {
			sqlStr = andStr + " AND (" + orStr + ") "
		} else {
			sqlStr = andStr
		}
	} else {
		if orStr != "" {
			sqlStr = orStr
		}
	}
	return sqlStr, vals
}

func genOrderSQL(filter sentity.Filter, keymap map[string]string) (order string) {
	if len(filter.Odr) > 0 {
		commaFlag := false
		for _, odr := range filter.Odr {
			if commaFlag {
				order = order + ", "
			}
			for k, v := range odr {
				if ik, ok := outerField2innerField(k, keymap); ok {
					if v {
						order = order + ik + " DESC "
					} else {
						order = order + ik + " ASC "
					}
					commaFlag = true
				}
			}
		}
	}
	return
}

func genSqlWithTable(tableName string, cond sentity.Filter, keymap map[string]string) (string, []interface{}) {
	var andStr = ""
	var vals = []interface{}{}

	for _, item := range cond.AndCond {
		if ik, ok := outerField2innerField(item.OuterName, keymap); ok {
			if andStr != "" {
				andStr += " AND "
			}
			switch item.Cond {
			case sentity.IsNotNull, sentity.IsNull:
				andStr += fmt.Sprintf(" %s %s", ik, item.Cond)
			case sentity.In:
				andStr += fmt.Sprintf(" %s %s (?)", ik, item.Cond)
				vals = append(vals, item.Value)
			case sentity.NotIn:
				andStr += fmt.Sprintf(" %s %s (SELECT %s FROM %s WHERE %s IN (?))", ik, item.Cond, ik, tableName, ik)
				vals = append(vals, item.Value)
			case sentity.Like, sentity.NotLike:
				andStr += fmt.Sprintf(" %s %s ?", ik, item.Cond)
				vals = append(vals, fmt.Sprintf("%%%v%%", item.Value))
			default:
				andStr += fmt.Sprintf(" %s %s ?", ik, item.Cond)
				vals = append(vals, item.Value)
			}
		}
	}
	var orStr = ""
	if len(cond.OrCond) > 1 {
		for _, item := range cond.OrCond {
			if ik, ok := outerField2innerField(item.OuterName, keymap); ok {
				if orStr != "" {
					orStr += " OR "
				}
				switch item.Cond {
				case sentity.IsNotNull, sentity.IsNull:
					orStr += fmt.Sprintf(" %s %s", ik, item.Cond)
				case sentity.In:
					orStr += fmt.Sprintf(" %s %s (?)", ik, item.Cond)
					vals = append(vals, item.Value)
				case sentity.NotIn:
					orStr += fmt.Sprintf(" %s %s (SELECT %s FROM %s WHERE %s IN (?))", ik, item.Cond, ik, tableName, ik)
					vals = append(vals, item.Value)
				case sentity.Like, sentity.NotLike:
					orStr += fmt.Sprintf(" %s %s ?", ik, item.Cond)
					vals = append(vals, fmt.Sprintf("%%%v%%", item.Value))
				default:
					orStr += fmt.Sprintf(" %s %s ?", ik, item.Cond)
					vals = append(vals, item.Value)
				}
			}
		}
	}
	sqlStr := ""
	if andStr != "" {
		if orStr != "" {
			sqlStr = andStr + " AND (" + orStr + ") "
		} else {
			sqlStr = andStr
		}
	} else {
		if orStr != "" {
			sqlStr = orStr
		}
	}
	return sqlStr, vals

}
