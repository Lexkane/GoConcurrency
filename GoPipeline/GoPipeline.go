package main

func lineListSource(ctx context.Context, lines ...string) (
	<-chan string, <-chan error, error) {
	if len(lines) == 0 {
		// Handle an error that occurs before the goroutine begins.
		return nil, nil, errors.Errorf("no lines provided")
	}
	out := make(chan string)
	errc := make(chan error, 1)
	go func() {
		defer close(out)
		defer close(errc)
		for lineIndex, line := range lines {
			if line == "" {
				// Handle an error that occurs during the goroutine.
				errc <- errors.Errorf("line %v is empty", lineIndex+1)
				return
			}
			// Send the data to the output channel but return early
			// if the context has been cancelled.
			select {
			case out <- line:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out, errc, nil
}

func sink(ctx context.Context, in <-chan int64) (
	<-chan error, error) {
	errc := make(chan error, 1)
	go func() {
		defer close(errc)
		for n := range in {
			if n >= 100 {
				// Handle an error that occurs during the goroutine.
				errc <- errors.Errorf("number %v is too large", n)
				return
			}
			fmt.Printf("sink: %v\n", n)
		}
	}()
	return errc, nil
}

func runSimplePipeline(base int, lines []string) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	var errcList []<-chan error
	// Source pipeline stage.
	linec, errc, err := lineListSource(ctx, lines...)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)
	// Transformer pipeline stage.
	numberc, errc, err := lineParser(ctx, base, linec)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)
	// Sink pipeline stage.
	errc, err = sink(ctx, numberc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)
	fmt.Println("Pipeline started. Waiting for pipeline to complete.")
	return WaitForPipeline(errcList...)
}

// WaitForPipeline waits for results from all error channels.
// It returns early on the first error.
func WaitForPipeline(errs ...<-chan error) error {
	errc := MergeErrors(errs...)
	for err := range errc {
		if err != nil {
			return err
		}
	}
	return nil
}

// MergeErrors merges multiple channels of errors.
// Based on https://blog.golang.org/pipelines.
func MergeErrors(cs ...<-chan error) <-chan error {
	var wg sync.WaitGroup
	// We must ensure that the output channel has the capacity to
	// hold as many errors
	// as there are error channels.
	// This will ensure that it never blocks, even
	// if WaitForPipeline returns early.
	out := make(chan error, len(cs))
	// Start an output goroutine for each input channel in cs.  output
	// copies values from c to out until c is closed, then calls
	// wg.Done.
	output := func(c <-chan error) {
		for n := range c {
			out <- n
		}
		wg.Done()
	}
	wg.Add(len(cs))
	for _, c := range cs {
		go output(c)
	}
	// Start a goroutine to close out once all the output goroutines
	// are done.  This must start after the wg.Add call.
	go func() {
		wg.Wait()
		close(out)
	}()
	return out
}

func splitter(ctx context.Context, in <-chan int64) (
	<-chan int64, <-chan int64, <-chan error, error) {
	out1 := make(chan int64)
	out2 := make(chan int64)
	errc := make(chan error, 1)
	go func() {
		defer close(out1)
		defer close(out2)
		defer close(errc)
		for n := range in {
			// Send the data to the output channel 1 but return early
			// if the context has been cancelled.
			select {
			case out1 <- n:
			case <-ctx.Done():
				return
			}
			// Send the data to the output channel 2 but return early
			// if the context has been cancelled.
			select {
			case out2 <- n:
			case <-ctx.Done():
				return
			}
		}
	}()
	return out1, out2, errc, nil
}

func runComplexPipeline(base int, lines []string) error {
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	var errcList []<-chan error
	// Source pipeline stage.
	linec, errc, err := lineListSource(ctx, lines...)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)
	// Transformer pipeline stage 1.
	numberc, errc, err := lineParser(ctx, base, linec)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)
	// Transformer pipeline stage 2.
	numberc1, numberc2, errc, err := splitter(ctx, numberc)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)
	// Transformer pipeline stage 3.
	numberc3, errc, err := squarer(ctx, numberc1)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)
	// Sink pipeline stage 1.
	errc, err = sink(ctx, numberc3)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)
	// Sink pipeline stage 2.
	errc, err = sink(ctx, numberc2)
	if err != nil {
		return err
	}
	errcList = append(errcList, errc)
	fmt.Println("Pipeline started. Waiting for pipeline to complete.")
	return WaitForPipeline(errcList...)
}

func main() {

}
