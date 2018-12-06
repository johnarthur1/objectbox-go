// Code generated by ObjectBox; DO NOT EDIT.

package iot

import (
	"github.com/google/flatbuffers/go"
	"github.com/objectbox/objectbox-go/objectbox"
	"github.com/objectbox/objectbox-go/objectbox/fbutils"
)

type event_EntityInfo struct {
	Id  objectbox.TypeId
	Uid uint64
}

var EventBinding = event_EntityInfo{
	Id:  1,
	Uid: 1468539308767086854,
}

// Event_ contains type-based Property helpers to facilitate some common operations such as Queries.
var Event_ = struct {
	Id      *objectbox.PropertyUint64
	Uid     *objectbox.PropertyString
	Device  *objectbox.PropertyString
	Date    *objectbox.PropertyInt64
	Picture *objectbox.PropertyByteVector
}{
	Id: &objectbox.PropertyUint64{
		Property: &objectbox.Property{
			Id: 1,
		},
	},
	Uid: &objectbox.PropertyString{
		Property: &objectbox.Property{
			Id: 4,
		},
	},
	Device: &objectbox.PropertyString{
		Property: &objectbox.Property{
			Id: 2,
		},
	},
	Date: &objectbox.PropertyInt64{
		Property: &objectbox.Property{
			Id: 3,
		},
	},
	Picture: &objectbox.PropertyByteVector{
		Property: &objectbox.Property{
			Id: 5,
		},
	},
}

// GeneratorVersion is called by the ObjectBox to verify the compatibility of the generator used to generate this code
func (event_EntityInfo) GeneratorVersion() int {
	return 1
}

// AddToModel is called by the ObjectBox during model build
func (event_EntityInfo) AddToModel(model *objectbox.Model) {
	model.Entity("Event", 1, 1468539308767086854)
	model.Property("Id", objectbox.PropertyType_Long, 1, 3098166604415018001)
	model.PropertyFlags(objectbox.PropertyFlags_ID)
	model.Property("Uid", objectbox.PropertyType_String, 4, 472416569173577818)
	model.PropertyFlags(objectbox.PropertyFlags_UNIQUE)
	model.PropertyIndex(1, 3297791712577314158)
	model.Property("Device", objectbox.PropertyType_String, 2, 1213411729427304641)
	model.Property("Date", objectbox.PropertyType_Date, 3, 5907655274386702697)
	model.Property("Picture", objectbox.PropertyType_ByteVector, 5, 6024563395733984005)
	model.EntityLastPropertyId(5, 6024563395733984005)
}

// GetId is called by the ObjectBox during Put operations to check for existing ID on an object
func (event_EntityInfo) GetId(object interface{}) (uint64, error) {
	return object.(*Event).Id, nil
}

// SetId is called by the ObjectBox during Put to update an ID on an object that has just been inserted
func (event_EntityInfo) SetId(object interface{}, id uint64) error {
	object.(*Event).Id = id
	return nil
}

// Flatten is called by the ObjectBox to transform an object to a FlatBuffer
func (event_EntityInfo) Flatten(object interface{}, fbb *flatbuffers.Builder, id uint64) {
	obj := object.(*Event)
	var offsetUid = fbutils.CreateStringOffset(fbb, obj.Uid)
	var offsetDevice = fbutils.CreateStringOffset(fbb, obj.Device)
	var offsetPicture = fbutils.CreateByteVectorOffset(fbb, obj.Picture)

	// build the FlatBuffers object
	fbb.StartObject(5)
	fbutils.SetUint64Slot(fbb, 0, id)
	fbutils.SetUOffsetTSlot(fbb, 3, offsetUid)
	fbutils.SetUOffsetTSlot(fbb, 1, offsetDevice)
	fbutils.SetInt64Slot(fbb, 2, obj.Date)
	fbutils.SetUOffsetTSlot(fbb, 4, offsetPicture)
}

// ToObject is called by the ObjectBox to load an object from a FlatBuffer
func (event_EntityInfo) ToObject(bytes []byte) interface{} {
	table := &flatbuffers.Table{
		Bytes: bytes,
		Pos:   flatbuffers.GetUOffsetT(bytes),
	}

	return &Event{
		Id:      table.GetUint64Slot(4, 0),
		Uid:     fbutils.GetStringSlot(table, 10),
		Device:  fbutils.GetStringSlot(table, 6),
		Date:    table.GetInt64Slot(8, 0),
		Picture: fbutils.GetByteVectorSlot(table, 12),
	}
}

// MakeSlice is called by the ObjectBox to construct a new slice to hold the read objects
func (event_EntityInfo) MakeSlice(capacity int) interface{} {
	return make([]*Event, 0, capacity)
}

// AppendToSlice is called by the ObjectBox to fill the slice of the read objects
func (event_EntityInfo) AppendToSlice(slice interface{}, object interface{}) interface{} {
	return append(slice.([]*Event), object.(*Event))
}

// Box provides CRUD access to Event objects
type EventBox struct {
	*objectbox.Box
}

// BoxForEvent opens a box of Event objects
func BoxForEvent(ob *objectbox.ObjectBox) *EventBox {
	return &EventBox{
		Box: ob.InternalBox(1),
	}
}

// Put synchronously inserts/updates a single object.
// In case the Id is not specified, it would be assigned automatically (auto-increment).
// When inserting, the Event.Id property on the passed object will be assigned the new ID as well.
func (box *EventBox) Put(object *Event) (uint64, error) {
	return box.Box.Put(object)
}

// PutAsync asynchronously inserts/updates a single object.
// When inserting, the Event.Id property on the passed object will be assigned the new ID as well.
//
// It's executed on a separate internal thread for better performance.
//
// There are two main use cases:
//
// 1) "Put & Forget:" you gain faster puts as you don't have to wait for the transaction to finish.
//
// 2) Many small transactions: if your write load is typically a lot of individual puts that happen in parallel,
// this will merge small transactions into bigger ones. This results in a significant gain in overall throughput.
//
//
// In situations with (extremely) high async load, this method may be throttled (~1ms) or delayed (<1s).
// In the unlikely event that the object could not be enqueued after delaying, an error will be returned.
//
// Note that this method does not give you hard durability guarantees like the synchronous Put provides.
// There is a small time window (typically 3 ms) in which the data may not have been committed durably yet.
func (box *EventBox) PutAsync(object *Event) (uint64, error) {
	return box.Box.PutAsync(object)
}

// PutAll inserts multiple objects in single transaction.
// In case Ids are not set on the objects, they would be assigned automatically (auto-increment).
//
// Returns: IDs of the put objects (in the same order).
// When inserting, the Event.Id property on the objects in the slice will be assigned the new IDs as well.
//
// Note: In case an error occurs during the transaction, some of the objects may already have the Event.Id assigned
// even though the transaction has been rolled back and the objects are not stored under those IDs.
//
// Note: The slice may be empty or even nil; in both cases, an empty IDs slice and no error is returned.
func (box *EventBox) PutAll(objects []*Event) ([]uint64, error) {
	return box.Box.PutAll(objects)
}

// Get reads a single object.
//
// Returns nil (and no error) in case the object with the given ID doesn't exist.
func (box *EventBox) Get(id uint64) (*Event, error) {
	object, err := box.Box.Get(id)
	if err != nil {
		return nil, err
	} else if object == nil {
		return nil, nil
	}
	return object.(*Event), nil
}

// Get reads all stored objects
func (box *EventBox) GetAll() ([]*Event, error) {
	objects, err := box.Box.GetAll()
	if err != nil {
		return nil, err
	}
	return objects.([]*Event), nil
}

// Remove deletes a single object
func (box *EventBox) Remove(object *Event) (err error) {
	return box.Box.Remove(object.Id)
}

// Creates a query with the given conditions. Use the fields of the Event_ struct to create conditions.
// Keep the *EventQuery if you intend to execute the query multiple times.
// Note: this function panics if you try to create illegal queries; e.g. use properties of an alien type.
// This is typically a programming error. Use QueryOrError instead if you want the explicit error check.
func (box *EventBox) Query(conditions ...objectbox.Condition) *EventQuery {
	return &EventQuery{
		box.Box.Query(conditions...),
	}
}

// Creates a query with the given conditions. Use the fields of the Event_ struct to create conditions.
// Keep the *EventQuery if you intend to execute the query multiple times.
func (box *EventBox) QueryOrError(conditions ...objectbox.Condition) (*EventQuery, error) {
	if query, err := box.Box.QueryOrError(conditions...); err != nil {
		return nil, err
	} else {
		return &EventQuery{query}, nil
	}
}

// Query provides a way to search stored objects
//
// For example, you can find all Event which Id is either 42 or 47:
// 		box.Query(Event_.Id.In(42, 47)).Find()
type EventQuery struct {
	*objectbox.Query
}

// Find returns all objects matching the query
func (query *EventQuery) Find() ([]*Event, error) {
	objects, err := query.Query.Find()
	if err != nil {
		return nil, err
	}
	return objects.([]*Event), nil
}

type reading_EntityInfo struct {
	Id  objectbox.TypeId
	Uid uint64
}

var ReadingBinding = reading_EntityInfo{
	Id:  2,
	Uid: 5284076134434938613,
}

// Reading_ contains type-based Property helpers to facilitate some common operations such as Queries.
var Reading_ = struct {
	Id              *objectbox.PropertyUint64
	Date            *objectbox.PropertyInt64
	EventId         *objectbox.PropertyUint64
	ValueName       *objectbox.PropertyString
	ValueString     *objectbox.PropertyString
	ValueInteger    *objectbox.PropertyInt64
	ValueFloating   *objectbox.PropertyFloat64
	ValueInt32      *objectbox.PropertyInt32
	ValueFloating32 *objectbox.PropertyFloat32
}{
	Id: &objectbox.PropertyUint64{
		Property: &objectbox.Property{
			Id: 1,
		},
	},
	Date: &objectbox.PropertyInt64{
		Property: &objectbox.Property{
			Id: 2,
		},
	},
	EventId: &objectbox.PropertyUint64{
		Property: &objectbox.Property{
			Id: 3,
		},
	},
	ValueName: &objectbox.PropertyString{
		Property: &objectbox.Property{
			Id: 4,
		},
	},
	ValueString: &objectbox.PropertyString{
		Property: &objectbox.Property{
			Id: 5,
		},
	},
	ValueInteger: &objectbox.PropertyInt64{
		Property: &objectbox.Property{
			Id: 6,
		},
	},
	ValueFloating: &objectbox.PropertyFloat64{
		Property: &objectbox.Property{
			Id: 7,
		},
	},
	ValueInt32: &objectbox.PropertyInt32{
		Property: &objectbox.Property{
			Id: 8,
		},
	},
	ValueFloating32: &objectbox.PropertyFloat32{
		Property: &objectbox.Property{
			Id: 9,
		},
	},
}

// GeneratorVersion is called by the ObjectBox to verify the compatibility of the generator used to generate this code
func (reading_EntityInfo) GeneratorVersion() int {
	return 1
}

// AddToModel is called by the ObjectBox during model build
func (reading_EntityInfo) AddToModel(model *objectbox.Model) {
	model.Entity("Reading", 2, 5284076134434938613)
	model.Property("Id", objectbox.PropertyType_Long, 1, 3968063745680890327)
	model.PropertyFlags(objectbox.PropertyFlags_ID)
	model.Property("Date", objectbox.PropertyType_Date, 2, 4852407661923085028)
	model.Property("EventId", objectbox.PropertyType_Relation, 3, 1403806151574554320)
	model.PropertyRelation("Event", 2, 2642563953244304959)
	model.Property("ValueName", objectbox.PropertyType_String, 4, 5626221656121286670)
	model.Property("ValueString", objectbox.PropertyType_String, 5, 7303099924122013060)
	model.Property("ValueInteger", objectbox.PropertyType_Long, 6, 1404333021836291657)
	model.Property("ValueFloating", objectbox.PropertyType_Double, 7, 7102253623343671118)
	model.Property("ValueInt32", objectbox.PropertyType_Int, 8, 7566830186276557216)
	model.Property("ValueFloating32", objectbox.PropertyType_Float, 9, 6040892611651481730)
	model.EntityLastPropertyId(9, 6040892611651481730)
}

// GetId is called by the ObjectBox during Put operations to check for existing ID on an object
func (reading_EntityInfo) GetId(object interface{}) (uint64, error) {
	return object.(*Reading).Id, nil
}

// SetId is called by the ObjectBox during Put to update an ID on an object that has just been inserted
func (reading_EntityInfo) SetId(object interface{}, id uint64) error {
	object.(*Reading).Id = id
	return nil
}

// Flatten is called by the ObjectBox to transform an object to a FlatBuffer
func (reading_EntityInfo) Flatten(object interface{}, fbb *flatbuffers.Builder, id uint64) {
	obj := object.(*Reading)
	var offsetValueName = fbutils.CreateStringOffset(fbb, obj.ValueName)
	var offsetValueString = fbutils.CreateStringOffset(fbb, obj.ValueString)

	// build the FlatBuffers object
	fbb.StartObject(9)
	fbutils.SetUint64Slot(fbb, 0, id)
	fbutils.SetInt64Slot(fbb, 1, obj.Date)
	fbutils.SetUint64Slot(fbb, 2, obj.EventId)
	fbutils.SetUOffsetTSlot(fbb, 3, offsetValueName)
	fbutils.SetUOffsetTSlot(fbb, 4, offsetValueString)
	fbutils.SetInt64Slot(fbb, 5, obj.ValueInteger)
	fbutils.SetFloat64Slot(fbb, 6, obj.ValueFloating)
	fbutils.SetInt32Slot(fbb, 7, obj.ValueInt32)
	fbutils.SetFloat32Slot(fbb, 8, obj.ValueFloating32)
}

// ToObject is called by the ObjectBox to load an object from a FlatBuffer
func (reading_EntityInfo) ToObject(bytes []byte) interface{} {
	table := &flatbuffers.Table{
		Bytes: bytes,
		Pos:   flatbuffers.GetUOffsetT(bytes),
	}

	return &Reading{
		Id:              table.GetUint64Slot(4, 0),
		Date:            table.GetInt64Slot(6, 0),
		EventId:         table.GetUint64Slot(8, 0),
		ValueName:       fbutils.GetStringSlot(table, 10),
		ValueString:     fbutils.GetStringSlot(table, 12),
		ValueInteger:    table.GetInt64Slot(14, 0),
		ValueFloating:   table.GetFloat64Slot(16, 0),
		ValueInt32:      table.GetInt32Slot(18, 0),
		ValueFloating32: table.GetFloat32Slot(20, 0),
	}
}

// MakeSlice is called by the ObjectBox to construct a new slice to hold the read objects
func (reading_EntityInfo) MakeSlice(capacity int) interface{} {
	return make([]*Reading, 0, capacity)
}

// AppendToSlice is called by the ObjectBox to fill the slice of the read objects
func (reading_EntityInfo) AppendToSlice(slice interface{}, object interface{}) interface{} {
	return append(slice.([]*Reading), object.(*Reading))
}

// Box provides CRUD access to Reading objects
type ReadingBox struct {
	*objectbox.Box
}

// BoxForReading opens a box of Reading objects
func BoxForReading(ob *objectbox.ObjectBox) *ReadingBox {
	return &ReadingBox{
		Box: ob.InternalBox(2),
	}
}

// Put synchronously inserts/updates a single object.
// In case the Id is not specified, it would be assigned automatically (auto-increment).
// When inserting, the Reading.Id property on the passed object will be assigned the new ID as well.
func (box *ReadingBox) Put(object *Reading) (uint64, error) {
	return box.Box.Put(object)
}

// PutAsync asynchronously inserts/updates a single object.
// When inserting, the Reading.Id property on the passed object will be assigned the new ID as well.
//
// It's executed on a separate internal thread for better performance.
//
// There are two main use cases:
//
// 1) "Put & Forget:" you gain faster puts as you don't have to wait for the transaction to finish.
//
// 2) Many small transactions: if your write load is typically a lot of individual puts that happen in parallel,
// this will merge small transactions into bigger ones. This results in a significant gain in overall throughput.
//
//
// In situations with (extremely) high async load, this method may be throttled (~1ms) or delayed (<1s).
// In the unlikely event that the object could not be enqueued after delaying, an error will be returned.
//
// Note that this method does not give you hard durability guarantees like the synchronous Put provides.
// There is a small time window (typically 3 ms) in which the data may not have been committed durably yet.
func (box *ReadingBox) PutAsync(object *Reading) (uint64, error) {
	return box.Box.PutAsync(object)
}

// PutAll inserts multiple objects in single transaction.
// In case Ids are not set on the objects, they would be assigned automatically (auto-increment).
//
// Returns: IDs of the put objects (in the same order).
// When inserting, the Reading.Id property on the objects in the slice will be assigned the new IDs as well.
//
// Note: In case an error occurs during the transaction, some of the objects may already have the Reading.Id assigned
// even though the transaction has been rolled back and the objects are not stored under those IDs.
//
// Note: The slice may be empty or even nil; in both cases, an empty IDs slice and no error is returned.
func (box *ReadingBox) PutAll(objects []*Reading) ([]uint64, error) {
	return box.Box.PutAll(objects)
}

// Get reads a single object.
//
// Returns nil (and no error) in case the object with the given ID doesn't exist.
func (box *ReadingBox) Get(id uint64) (*Reading, error) {
	object, err := box.Box.Get(id)
	if err != nil {
		return nil, err
	} else if object == nil {
		return nil, nil
	}
	return object.(*Reading), nil
}

// Get reads all stored objects
func (box *ReadingBox) GetAll() ([]*Reading, error) {
	objects, err := box.Box.GetAll()
	if err != nil {
		return nil, err
	}
	return objects.([]*Reading), nil
}

// Remove deletes a single object
func (box *ReadingBox) Remove(object *Reading) (err error) {
	return box.Box.Remove(object.Id)
}

// Creates a query with the given conditions. Use the fields of the Reading_ struct to create conditions.
// Keep the *ReadingQuery if you intend to execute the query multiple times.
// Note: this function panics if you try to create illegal queries; e.g. use properties of an alien type.
// This is typically a programming error. Use QueryOrError instead if you want the explicit error check.
func (box *ReadingBox) Query(conditions ...objectbox.Condition) *ReadingQuery {
	return &ReadingQuery{
		box.Box.Query(conditions...),
	}
}

// Creates a query with the given conditions. Use the fields of the Reading_ struct to create conditions.
// Keep the *ReadingQuery if you intend to execute the query multiple times.
func (box *ReadingBox) QueryOrError(conditions ...objectbox.Condition) (*ReadingQuery, error) {
	if query, err := box.Box.QueryOrError(conditions...); err != nil {
		return nil, err
	} else {
		return &ReadingQuery{query}, nil
	}
}

// Query provides a way to search stored objects
//
// For example, you can find all Reading which Id is either 42 or 47:
// 		box.Query(Reading_.Id.In(42, 47)).Find()
type ReadingQuery struct {
	*objectbox.Query
}

// Find returns all objects matching the query
func (query *ReadingQuery) Find() ([]*Reading, error) {
	objects, err := query.Query.Find()
	if err != nil {
		return nil, err
	}
	return objects.([]*Reading), nil
}