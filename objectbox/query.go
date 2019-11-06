/*
 * Copyright 2019 ObjectBox Ltd. All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package objectbox

/*
#include <stdlib.h>
#include "objectbox.h"
*/
import "C"
import (
	"errors"
	"fmt"
	"runtime"
	"sync"
	"unsafe"
)

// A Query allows to search for objects matching user defined conditions.
//
// For example, you can find all people whose last name starts with an 'N':
// 		box.Query(Person_.LastName.HasPrefix("N", false)).Find()
// Note that Person_ is a struct generated by ObjectBox allowing to conveniently reference properties.
type Query struct {
	entity          *entity
	objectBox       *ObjectBox
	box             *Box
	cQuery          *C.OBX_query
	closeMutex      sync.Mutex
	offset          uint64
	limit           uint64
	linkedEntityIds []TypeId
}

// Close frees (native) resources held by this Query.
// Note that this is optional and not required because the GC invokes a finalizer automatically.
func (query *Query) Close() error {
	query.closeMutex.Lock()
	defer query.closeMutex.Unlock()

	if query.cQuery != nil {
		return cCall(func() C.obx_err {
			var err = C.obx_query_close(query.cQuery)
			query.cQuery = nil
			return err
		})
	}
	return nil
}

func queryFinalizer(query *Query) {
	err := query.Close()
	if err != nil {
		fmt.Printf("Error while finalizer closed query: %s", err)
	}
}

// The native query object in the ObjectBox core is not tied with other resources.
// Thus timing of the Close call is independent from other resources.
func (query *Query) installFinalizer() {
	runtime.SetFinalizer(query, queryFinalizer)
}

func (query *Query) errorClosed() error {
	return errors.New("illegal state; query was closed")
}

// Find returns all objects matching the query
func (query *Query) Find() (objects interface{}, err error) {
	if query.cQuery == nil {
		return 0, query.errorClosed()
	}

	const existingOnly = true
	if supportsBytesArray {
		var cFn = func() *C.OBX_bytes_array {
			return C.obx_query_find(query.cQuery, C.uint64_t(query.offset), C.uint64_t(query.limit))
		}
		return query.box.readManyObjects(existingOnly, cFn)
	}

	var cFn = func(visitorArg unsafe.Pointer) C.obx_err {
		return C.obx_query_visit(query.cQuery, dataVisitor, visitorArg,
			C.uint64_t(query.offset), C.uint64_t(query.limit))
	}
	return query.box.readUsingVisitor(existingOnly, cFn)
}

// Offset defines the index of the first object to process (how many objects to skip)
func (query *Query) Offset(offset uint64) *Query {
	query.offset = offset
	return query
}

// Limit sets the number of elements to process by the query
func (query *Query) Limit(limit uint64) *Query {
	query.limit = limit
	return query
}

// FindIds returns IDs of all objects matching the query
func (query *Query) FindIds() ([]uint64, error) {
	if query.cQuery == nil {
		return nil, query.errorClosed()
	}

	return cGetIds(func() *C.OBX_id_array {
		return C.obx_query_find_ids(query.cQuery, C.uint64_t(query.offset), C.uint64_t(query.limit))
	})
}

// Count returns the number of objects matching the query
func (query *Query) Count() (uint64, error) {
	// doesn't support offset/limit at this point
	if query.offset != 0 || query.limit != 0 {
		return 0, fmt.Errorf("limit/offset are not supported by Count at this moment")
	}

	if query.cQuery == nil {
		return 0, query.errorClosed()
	}

	var cResult C.uint64_t
	if err := cCall(func() C.obx_err { return C.obx_query_count(query.cQuery, &cResult) }); err != nil {
		return 0, err
	}
	return uint64(cResult), nil
}

// Remove permanently deletes all objects matching the query from the database
func (query *Query) Remove() (count uint64, err error) {
	// doesn't support offset/limit at this point
	if query.offset != 0 || query.limit != 0 {
		return 0, fmt.Errorf("limit/offset are not supported by Remove at this moment")
	}

	if query.cQuery == nil {
		return 0, query.errorClosed()
	}

	var cResult C.uint64_t
	if err := cCall(func() C.obx_err { return C.obx_query_remove(query.cQuery, &cResult) }); err != nil {
		return 0, err
	}
	return uint64(cResult), nil
}

// DescribeParams returns a string representation of the query conditions
func (query *Query) DescribeParams() (string, error) {
	if query.cQuery == nil {
		return "", query.errorClosed()
	}
	// no need to free, it's handled by the cQuery internally
	cResult := C.obx_query_describe_params(query.cQuery)

	return C.GoString(cResult), nil
}

func (query *Query) checkIdentifier(identifier propertyOrAlias) error {
	// NOTE: maybe validate if the alias was previously used in this query?
	if identifier.alias() != nil {
		return nil
	}

	var entityId = identifier.entityId()

	if query.entity.id == entityId {
		return nil
	}

	if query.linkedEntityIds != nil {
		for _, id := range query.linkedEntityIds {
			if id == entityId {
				return nil
			}
		}

		return fmt.Errorf("property from a different entity %d passed, expected one of %v",
			entityId, append([]TypeId{query.entity.id}, query.linkedEntityIds...))
	}

	return fmt.Errorf("property from a different entity %d passed, expected %d", entityId, query.entity.id)
}

// propertyOrAlias is used to identify a condition in a query.
// You can use either a BaseProperty (or any Property* type embedding it), or an Alias("str") call result.
type propertyOrAlias interface {
	propertyId() TypeId
	entityId() TypeId
	alias() *string
}

// SetStringParams changes query parameter values on the given property
func (query *Query) SetStringParams(identifier propertyOrAlias, values ...string) error {
	if err := query.checkIdentifier(identifier); err != nil {
		return err
	}

	if len(values) == 0 {
		return fmt.Errorf("no values given")
	}

	var cAlias *C.char
	if alias := identifier.alias(); alias != nil {
		cAlias = C.CString(*alias)
		defer C.free(unsafe.Pointer(cAlias))
	}

	if len(values) == 1 {
		return cCall(func() C.obx_err {
			cString := C.CString(values[0])
			defer C.free(unsafe.Pointer(cString))

			if cAlias != nil {
				return C.obx_query_string_param_alias(query.cQuery, cAlias, cString)
			}
			return C.obx_query_string_param(query.cQuery, C.obx_schema_id(identifier.entityId()), C.obx_schema_id(identifier.propertyId()), cString)
		})
	}

	return fmt.Errorf("too many values given")
}

// SetStringParamsIn changes query parameter values on the given property
func (query *Query) SetStringParamsIn(identifier propertyOrAlias, values ...string) error {
	if err := query.checkIdentifier(identifier); err != nil {
		return err
	}

	if len(values) == 0 {
		return fmt.Errorf("no values given")
	}

	var cAlias *C.char
	if alias := identifier.alias(); alias != nil {
		cAlias = C.CString(*alias)
		defer C.free(unsafe.Pointer(cAlias))
	}

	cStringArray := goStringArrayToC(values)
	defer cStringArray.free()

	return cCall(func() C.obx_err {
		if cAlias != nil {
			return C.obx_query_string_params_in_alias(query.cQuery, cAlias, cStringArray.cArray, C.int(cStringArray.size))
		}
		return C.obx_query_string_params_in(query.cQuery, C.obx_schema_id(identifier.entityId()), C.obx_schema_id(identifier.propertyId()), cStringArray.cArray, C.int(cStringArray.size))
	})
}

// SetInt64Params changes query parameter values on the given property
func (query *Query) SetInt64Params(identifier propertyOrAlias, values ...int64) error {
	if err := query.checkIdentifier(identifier); err != nil {
		return err
	}

	if len(values) == 0 {
		return fmt.Errorf("no values given")
	}

	var cAlias *C.char
	if alias := identifier.alias(); alias != nil {
		cAlias = C.CString(*alias)
		defer C.free(unsafe.Pointer(cAlias))
	}

	if len(values) == 1 {
		return cCall(func() C.obx_err {
			if cAlias != nil {
				return C.obx_query_int_param_alias(query.cQuery, cAlias, C.int64_t(values[0]))
			}
			return C.obx_query_int_param(query.cQuery, C.obx_schema_id(identifier.entityId()), C.obx_schema_id(identifier.propertyId()), C.int64_t(values[0]))
		})

	} else if len(values) == 2 {
		return cCall(func() C.obx_err {
			if cAlias != nil {
				return C.obx_query_int_params_alias(query.cQuery, cAlias, C.int64_t(values[0]), C.int64_t(values[1]))
			}
			return C.obx_query_int_params(query.cQuery, C.obx_schema_id(identifier.entityId()), C.obx_schema_id(identifier.propertyId()), C.int64_t(values[0]), C.int64_t(values[1]))
		})
	}

	return fmt.Errorf("too many values given")
}

// SetInt64ParamsIn changes query parameter values on the given property
func (query *Query) SetInt64ParamsIn(identifier propertyOrAlias, values ...int64) error {
	if err := query.checkIdentifier(identifier); err != nil {
		return err
	}

	if len(values) == 0 {
		return fmt.Errorf("no values given")
	}

	var cAlias *C.char
	if alias := identifier.alias(); alias != nil {
		cAlias = C.CString(*alias)
		defer C.free(unsafe.Pointer(cAlias))
	}

	return cCall(func() C.obx_err {
		if cAlias != nil {
			return C.obx_query_int64_params_in_alias(query.cQuery, cAlias, (*C.int64_t)(unsafe.Pointer(&values[0])), C.int(len(values)))
		}
		return C.obx_query_int64_params_in(query.cQuery, C.obx_schema_id(identifier.entityId()), C.obx_schema_id(identifier.propertyId()), (*C.int64_t)(unsafe.Pointer(&values[0])), C.int(len(values)))
	})
}

// SetInt32ParamsIn changes query parameter values on the given property
func (query *Query) SetInt32ParamsIn(identifier propertyOrAlias, values ...int32) error {
	if err := query.checkIdentifier(identifier); err != nil {
		return err
	}

	if len(values) == 0 {
		return fmt.Errorf("no values given")
	}

	var cAlias *C.char
	if alias := identifier.alias(); alias != nil {
		cAlias = C.CString(*alias)
		defer C.free(unsafe.Pointer(cAlias))
	}

	return cCall(func() C.obx_err {
		if cAlias != nil {
			return C.obx_query_int32_params_in_alias(query.cQuery, cAlias, (*C.int32_t)(unsafe.Pointer(&values[0])), C.int(len(values)))
		}
		return C.obx_query_int32_params_in(query.cQuery, C.obx_schema_id(identifier.entityId()), C.obx_schema_id(identifier.propertyId()), (*C.int32_t)(unsafe.Pointer(&values[0])), C.int(len(values)))
	})
}

// SetFloat64Params changes query parameter values on the given property
func (query *Query) SetFloat64Params(identifier propertyOrAlias, values ...float64) error {
	if err := query.checkIdentifier(identifier); err != nil {
		return err
	}

	if len(values) == 0 {
		return fmt.Errorf("no values given")
	}

	var cAlias *C.char
	if alias := identifier.alias(); alias != nil {
		cAlias = C.CString(*alias)
		defer C.free(unsafe.Pointer(cAlias))
	}

	if len(values) == 1 {
		return cCall(func() C.obx_err {
			if cAlias != nil {
				return C.obx_query_double_param_alias(query.cQuery, cAlias, C.double(values[0]))
			}
			return C.obx_query_double_param(query.cQuery, C.obx_schema_id(identifier.entityId()), C.obx_schema_id(identifier.propertyId()), C.double(values[0]))
		})

	} else if len(values) == 2 {
		return cCall(func() C.obx_err {
			if cAlias != nil {
				return C.obx_query_double_params_alias(query.cQuery, cAlias, C.double(values[0]), C.double(values[1]))
			}
			return C.obx_query_double_params(query.cQuery, C.obx_schema_id(identifier.entityId()), C.obx_schema_id(identifier.propertyId()), C.double(values[0]), C.double(values[1]))
		})

	}

	return fmt.Errorf("too many values given")
}

// SetBytesParams changes query parameter values on the given property
func (query *Query) SetBytesParams(identifier propertyOrAlias, values ...[]byte) error {
	if err := query.checkIdentifier(identifier); err != nil {
		return err
	}

	if len(values) == 0 {
		return fmt.Errorf("no values given")

	} else if len(values) > 1 {
		return fmt.Errorf("too many values given")
	}

	var cAlias *C.char
	if alias := identifier.alias(); alias != nil {
		cAlias = C.CString(*alias)
		defer C.free(unsafe.Pointer(cAlias))
	}

	return cCall(func() C.obx_err {
		var value = values[0]
		if value == nil {
			value = []byte{}
		}
		if cAlias != nil {
			return C.obx_query_bytes_param_alias(query.cQuery, cAlias, cBytesPtr(value), C.size_t(len(value)))
		}
		return C.obx_query_bytes_param(query.cQuery, C.obx_schema_id(identifier.entityId()), C.obx_schema_id(identifier.propertyId()), cBytesPtr(value), C.size_t(len(value)))
	})
}
