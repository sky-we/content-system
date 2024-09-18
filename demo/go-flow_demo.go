package main

import (
	"encoding/json"
	"fmt"
	flow "github.com/s8sg/goflow/flow/v1"
	goflow "github.com/s8sg/goflow/v1"
	"math/rand"
	"strconv"
)

func Input(data []byte, options map[string][]string) ([]byte, error) {
	var input map[string]int
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, err
	}
	outputInt := input["input"]
	fmt.Println("Input data", outputInt)

	return []byte(strconv.Itoa(outputInt)), nil

}

func Step1(data []byte, options map[string][]string) ([]byte, error) {
	num, _ := strconv.Atoi(string(data))
	outputInt := num + rand.Intn(10) + 1
	fmt.Println("Step1 result", outputInt)

	return []byte(strconv.Itoa(outputInt)), nil
}

func Step2(data []byte, options map[string][]string) ([]byte, error) {
	num, _ := strconv.Atoi(string(data))
	outputInt := num + rand.Intn(100) + 1
	fmt.Println("Step2 result", outputInt)

	return []byte(strconv.Itoa(outputInt)), nil
}
func Output(data []byte, options map[string][]string) ([]byte, error) {
	fmt.Println("output data", string(data))
	return []byte("ok"), nil
}

func MyFlow(flow *flow.Workflow, context *flow.Context) error {
	dag := flow.Dag()
	dag.Node("input", Input)
	dag.Node("step1", Step1)
	dag.Node("step2", Step2)
	dag.Node("output", Output)

	dag.Edge("input", "step1")
	dag.Edge("step1", "step2")
	dag.Edge("step2", "output")

	return nil
}
func main() {
	fs := goflow.FlowService{
		Port:              8080,
		RedisURL:          "localhost:6379",
		WorkerConcurrency: 5,
	}

	if err := fs.Register("math", MyFlow); err != nil {
		panic(err)

	}
	if err := fs.Start(); err != nil {
		panic(err)
	}
}
