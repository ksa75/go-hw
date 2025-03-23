package hw06pipelineexecution

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	ch := in
	for _, stage := range stages {
		ch = stage(readCh(ch, done))
	}
	return ch
}

func readCh(in In, done In) Out {
	nextCh := make(Bi)
	go func() {
		defer close(nextCh)
		for {
			select {
			case <-done:
				<-in // Вычитываем из канала значение чтобы разблокировать писателя
				return
			case data, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done: // Проверка done перед отправкой
				case nextCh <- data:
				}
			}
		}
	}()
	return nextCh
}

// 			result := stages[3](stages[2](stages[1](stages[0](in))))
// 			return result
// }
