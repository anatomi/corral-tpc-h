// Code generated by "enumer -type=customerFields"; DO NOT EDIT.

//
package queries

import (
	"fmt"
)

const _customerFieldsName = "C_CUSTKEYC_NAMEC_ADDRESSC_NATIONKEYC_PHONEC_ACCTBALC_MKTSEGMENTC_COMMENT"

var _customerFieldsIndex = [...]uint8{0, 9, 15, 24, 35, 42, 51, 63, 72}

func (i customerFields) String() string {
	if i < 0 || i >= customerFields(len(_customerFieldsIndex)-1) {
		return fmt.Sprintf("customerFields(%d)", i)
	}
	return _customerFieldsName[_customerFieldsIndex[i]:_customerFieldsIndex[i+1]]
}

var _customerFieldsValues = []customerFields{0, 1, 2, 3, 4, 5, 6, 7}

var _customerFieldsNameToValueMap = map[string]customerFields{
	_customerFieldsName[0:9]:   0,
	_customerFieldsName[9:15]:  1,
	_customerFieldsName[15:24]: 2,
	_customerFieldsName[24:35]: 3,
	_customerFieldsName[35:42]: 4,
	_customerFieldsName[42:51]: 5,
	_customerFieldsName[51:63]: 6,
	_customerFieldsName[63:72]: 7,
}

// customerFieldsString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func customerFieldsString(s string) (customerFields, error) {
	if val, ok := _customerFieldsNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to customerFields values", s)
}

// customerFieldsValues returns all values of the enum
func customerFieldsValues() []customerFields {
	return _customerFieldsValues
}

// IsAcustomerFields returns "true" if the value is listed in the enum definition. "false" otherwise
func (i customerFields) IsAcustomerFields() bool {
	for _, v := range _customerFieldsValues {
		if i == v {
			return true
		}
	}
	return false
}
