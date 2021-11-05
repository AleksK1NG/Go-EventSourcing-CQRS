package es

// HandleCommand Aggregate commands' handler method
// Example
//
// func (a *OrderAggregate) HandleCommand(command interface{}) error {
//	switch c := command.(type) {
//	case *CreateOrderCommand:
//		return a.handleCreateOrderCommand(c)
//	case *OrderPaidCommand:
//		return a.handleOrderPaidCommand(c)
//	case *SubmitOrderCommand:
//		return a.handleSubmitOrderCommand(c)
//	default:
//		return errors.New("invalid command type")
//	}
//}
type HandleCommand interface {
	HandleCommand(command Command) error
}

// When process and update aggregate state on specified es.Event type
// Example:
//
//func (a *OrderAggregate) When(evt es.Event) error {
//
//	switch evt.GetEventType() {
//
//	case events.OrderCreated:
//		var eventData events.OrderCreatedData
//		if err := json.Unmarshal(evt.GetData(), &eventData); err != nil {
//			return err
//		}
//		a.Order.ItemsIDs = eventData.ItemsIDs
//		a.Order.Created = true
//		return nil
//
//	default:
//		return errors.New("invalid event type")
//	}
//}
type When interface {
	When(event Event) error
}

type when func(event Event) error

// Apply process Aggregate Event
type Apply interface {
	Apply(event Event) error
}

// Load create Aggregate state from Event's.
type Load interface {
	Load(events []Event) error
}

type Aggregate interface {
	When
	AggregateRoot
	HandleCommand
}

// AggregateRoot contains all methods of AggregateBase
type AggregateRoot interface {
	GetUncommittedEvents() []Event
	GetID() string
	SetID(id string) *AggregateBase
	GetVersion() uint64
	ClearUncommittedEvents()
	ToSnapshot()
	SetType(aggregateType AggregateType)
	GetType() AggregateType
	SetAppliedEvents(events []Event)
	GetAppliedEvents() []Event
	Load
	Apply
}

// AggregateType type of the Aggregate
type AggregateType string

// AggregateBase base aggregate contains all main necessary fields
type AggregateBase struct {
	ID                string
	Version           uint64
	AppliedEvents     []Event
	UncommittedEvents []Event
	Type              AggregateType
	when              when
}

// NewAggregateBase AggregateBase constructor, contains all main fields and methods,
// main aggregate must realize When interface and pass as argument to constructor
// Example of recommended aggregate constructor method:
//
// func NewOrderAggregate() *OrderAggregate {
//	orderAggregate := &OrderAggregate{
//		Order: models.NewOrder(),
//	}
//	base := es.NewAggregateBase(orderAggregate.When)
//	base.SetType(OrderAggregateType)
//	orderAggregate.AggregateBase = base
//	return orderAggregate
//}
func NewAggregateBase(when when) *AggregateBase {
	if when == nil {
		return nil
	}

	return &AggregateBase{
		Version:           0,
		AppliedEvents:     make([]Event, 0, 10),
		UncommittedEvents: make([]Event, 0, 10),
		when:              when,
	}
}

// SetID set AggregateBase ID
func (a *AggregateBase) SetID(id string) *AggregateBase {
	a.ID = id
	return a
}

// GetID get AggregateBase ID
func (a *AggregateBase) GetID() string {
	return a.ID
}

// SetType set AggregateBase AggregateType
func (a *AggregateBase) SetType(aggregateType AggregateType) {
	a.Type = aggregateType
}

// GetType get AggregateBase AggregateType
func (a *AggregateBase) GetType() AggregateType {
	return a.Type
}

// GetVersion get AggregateBase version
func (a *AggregateBase) GetVersion() uint64 {
	return a.Version
}

// ClearUncommittedEvents clear AggregateBase uncommitted Event's
func (a *AggregateBase) ClearUncommittedEvents() {
	a.UncommittedEvents = make([]Event, 0)
}

// GetAppliedEvents get AggregateBase applied Event's
func (a *AggregateBase) GetAppliedEvents() []Event {
	return a.AppliedEvents
}

// SetAppliedEvents set AggregateBase applied Event's
func (a *AggregateBase) SetAppliedEvents(events []Event) {
	a.AppliedEvents = events
}

// GetUncommittedEvents get AggregateBase uncommitted Event's
func (a *AggregateBase) GetUncommittedEvents() []Event {
	return a.UncommittedEvents
}

// Load add existing events from event store to aggregate using When interface method
func (a *AggregateBase) Load(events []Event) error {

	for _, evt := range events {
		if evt.GetAggregateID() != a.GetID() {
			return ErrInvalidAggregate
		}

		if err := a.when(evt); err != nil {
			return err
		}

		a.AppliedEvents = append(a.AppliedEvents, evt)
		a.Version++
	}

	return nil
}

// Apply push event to aggregate uncommitted events using When method
func (a *AggregateBase) Apply(event Event) error {
	if event.GetAggregateID() != a.GetID() {
		return ErrInvalidAggregateID
	}

	event.SetVersion(a.GetVersion())
	event.SetAggregateType(a.GetType())

	if err := a.when(event); err != nil {
		return err
	}

	a.UncommittedEvents = append(a.UncommittedEvents, event)
	a.Version++
	return nil
}

// ToSnapshot prepare AggregateBase for saving Snapshot.
func (a *AggregateBase) ToSnapshot() {
	a.AppliedEvents = append(a.AppliedEvents, a.UncommittedEvents...)
	a.ClearUncommittedEvents()
}
