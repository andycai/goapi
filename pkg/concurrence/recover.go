package concurrence

var (
	ReallyCrash = true
)

var PanicHandlers = []func(any){
	func(v any) {
	},
}

func HandleCrash(additionalHandlers ...func(any)) {
	if r := recover(); r != nil {
		for _, fn := range PanicHandlers {
			fn(r)
		}
		for _, fn := range additionalHandlers {
			fn(r)
		}
		if ReallyCrash {
			panic(r)
		}
	}
}

// usage:
// func Go(fn func()) {
// 	go func() {
// 		defer HandleCrash()
// 		fn()
// 	}()
// }
