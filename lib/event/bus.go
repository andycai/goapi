package event

import (
	"context"
	"fmt"
	"reflect"
	"sync"

	"github.com/andycai/goapi/lib/work"
)

// Event 定义了事件的基础接口，所有事件都应可被断言为 interface{}
// 实际项目中，可以定义一个更具体的 Event 接口，比如要求有 Name() 方法
type Event interface{}

// EventHandler 定义了事件处理函数的类型签名
// 使用泛型 T 来约束事件的具体类型，确保类型安全
type EventHandler[T Event] func(ctx context.Context, event T) error

// EventBus 结构体
type EventBus struct {
	handlers map[reflect.Type][]any // 存储事件类型到其处理函数列表的映射
	mu       sync.RWMutex           // 用于并发安全的读写锁
	workPool *work.WorkPool         // 用于异步执行事件处理
}

// NewEventBus 创建一个新的 EventBus 实例
func NewEventBus() *EventBus {
	return &EventBus{
		handlers: make(map[reflect.Type][]any),
		workPool: work.NewWorkPool(30, 100),
	}
}

// Subscribe 订阅一个事件类型
// handler 必须是 EventHandler[T] 类型
func Subscribe[T Event](bus *EventBus, handler EventHandler[T]) {
	bus.mu.Lock()
	defer bus.mu.Unlock()

	// 获取事件 T 的真实类型
	// 注意：不能直接用 reflect.TypeOf(T)，因为 T 只是一个类型参数
	// 需要一个 T 类型的实例或者使用 (*T)(nil) 来获取元素类型
	var t T
	eventType := reflect.TypeOf(t)
	if eventType == nil && reflect.TypeOf((*T)(nil)) != nil { // 处理接口类型事件
		eventType = reflect.TypeOf((*T)(nil)).Elem()
	} else if eventType == nil {
		// 这种情况比较少见，通常 T 会是一个具体的 struct 或 interface
		// 如果 T 是 interface{} 本身，这里会是 nil，需要特殊处理或禁止
		// 对于 interface{} 类型的事件，可能需要不同的注册和分发逻辑
		panic("eventbus: cannot subscribe to nil event type or generic interface{} directly without a concrete type assertion helper")
	}

	bus.handlers[eventType] = append(bus.handlers[eventType], handler)
	fmt.Printf("Subscribed handler for event type %s\n", eventType)
}

// Publish 发布一个事件
// 事件将异步分发给所有订阅了该事件类型的处理器
// 你也可以选择同步分发
func Publish[T Event](bus *EventBus, ctx context.Context, event T) {
	bus.mu.RLock() // 使用读锁，因为我们只是读取 handlers 맵
	defer bus.mu.RUnlock()

	eventType := reflect.TypeOf(event)
	if eventType == nil {
		fmt.Println("Warning: Publishing a nil type event. This usually means the event instance itself is nil.")
		return
	}

	fmt.Printf("Publishing event of type %s: %+v\n", eventType, event)

	if handlers, ok := bus.handlers[eventType]; ok {
		for _, h := range handlers {
			// 类型断言回具体的 EventHandler[T]
			// 这是必要的，因为我们在 map 中存储的是 interface{}
			if handler, typeOk := h.(EventHandler[T]); typeOk {
				// 异步执行，避免阻塞 Publish 调用
				// 实际项目中可能需要一个 goroutine 池来管理
				bus.workPool.PostWork("eventbus", &eventWorker[T]{
					ctx:     ctx,
					event:   event,
					handler: handler,
				})
			} else {
				// 这通常不应该发生，如果 Subscribe 逻辑正确的话
				fmt.Printf("Error: Handler for event type %s has incorrect type. Expected EventHandler[%s], got %T\n",
					eventType, eventType, h)
			}
		}
	} else {
		fmt.Printf("No handlers subscribed for event type %s\n", eventType)
	}
}

// eventWorker 实现了 work.PoolWorker 接口
type eventWorker[T Event] struct {
	ctx     context.Context
	event   T
	handler EventHandler[T]
}

func (w *eventWorker[T]) DoWork(workRoutine int) {
	if err := w.handler(w.ctx, w.event); err != nil {
		// 处理错误，例如日志记录
		fmt.Printf("Error executing handler for event %s: %v\n", reflect.TypeOf(w.event), err)
	}
}

// func main() {
// 	bus := event.NewEventBus()

// 	// 订阅事件
// 	event.Subscribe(bus, event.EventHandler[UserCreatedEvent](handleUserCreated))
// 	event.Subscribe(bus, event.EventHandler[UserCreatedEvent](sendWelcomeEmail))
// 	event.Subscribe(bus, event.EventHandler[OrderPlacedEvent](handleOrderPlaced))

// 	// 发布事件
// 	ctx := context.Background()

// 	userEvent := UserCreatedEvent{UserID: 1, Username: "Alice"}
// 	event.Publish(bus, ctx, userEvent)

// 	orderEvent := OrderPlacedEvent{OrderID: "ORD123", Amount: 99.99}
// 	event.Publish(bus, ctx, orderEvent)

// 	failedOrderEvent := OrderPlacedEvent{OrderID: "ORD000", Amount: -10.00}
// 	event.Publish(bus, ctx, failedOrderEvent)

// 	// 给异步处理器一些时间执行 (仅为演示)
// 	time.Sleep(1 * time.Second)
// 	fmt.Println("Done.")
// }
