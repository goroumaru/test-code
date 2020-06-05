package pipeline_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/goroumaru/test-code/pipeline"
)

func TestPipeline(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	intStream := pipeline.Generator(ctx, 1, 2, 3, 4) // int型可変スライスの入力値
	// out = {(in * 2) + 1}*2
	pipeline := pipeline.Multiply(ctx, pipeline.Add(ctx, pipeline.Multiply(ctx, intStream, 2), 1), 2)

	for i := range pipeline {
		fmt.Println(i)
	}
}
