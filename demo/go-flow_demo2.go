package main

import (
	"encoding/json"
	"fmt"
	flow "github.com/s8sg/goflow/flow/v1"
	goflow "github.com/s8sg/goflow/v1"
	"math/rand"
	"strconv"
)

func Input2(data []byte, options map[string][]string) ([]byte, error) {
	var input map[string]int
	if err := json.Unmarshal(data, &input); err != nil {
		return nil, err
	}
	outputInt := input["input"]
	fmt.Println("Input data", outputInt)

	return []byte(strconv.Itoa(outputInt)), nil

}

func AddTen(data []byte, options map[string][]string) ([]byte, error) {
	num, _ := strconv.Atoi(string(data))
	outputInt := num + rand.Intn(10) + 1
	fmt.Println("Step1 result", outputInt)

	return []byte(strconv.Itoa(outputInt)), nil
}

func AddTen2(data []byte, options map[string][]string) ([]byte, error) {
	num, _ := strconv.Atoi(string(data))
	outputInt := num + rand.Intn(10) + 1
	fmt.Println("Step1 result", outputInt)

	return []byte(strconv.Itoa(outputInt)), nil
}

func Aggregator(data []byte, options map[string][]string) ([]byte, error) {
	fmt.Println("Aggregator = ", string(data))
	return data, nil
}

func Expand10(data []byte, options map[string][]string) ([]byte, error) {
	num, _ := strconv.Atoi(string(data))
	outputInt := num * 10
	fmt.Println("expand10 result", outputInt)

	return []byte(strconv.Itoa(outputInt)), nil
}

func Expand100(data []byte, options map[string][]string) ([]byte, error) {
	num, _ := strconv.Atoi(string(data))
	outputInt := num * 100
	fmt.Println("expand100 result", outputInt)

	return []byte(strconv.Itoa(outputInt)), nil
}
func Output2(data []byte, options map[string][]string) ([]byte, error) {
	fmt.Println("data", string(data))
	return []byte("ok"), nil
}

func MyFlow2(workFlow *flow.Workflow, context *flow.Context) error {
	dag := workFlow.Dag()
	dag.Node("input", Input2)
	dag.Node("add-ten", AddTen)
	dag.Node("add-ten2", AddTen2)
	dag.Node("aggregator", Aggregator, flow.Aggregator(func(m map[string][]byte) ([]byte, error) {
		a, _ := strconv.Atoi(string(m["add-ten"]))
		b, _ := strconv.Atoi(string(m["add-ten2"]))
		num := a + b
		fmt.Println("aggregator = ", num)
		return []byte(strconv.Itoa(num)), nil
	}))

	branches := dag.ConditionalBranch("judgeTen", []string{"moreThan", "lessThan"}, func(bytes []byte) []string {

		num, _ := strconv.Atoi(string(bytes))
		fmt.Println("ConditionBranch num=", num)
		if num > 10 {
			return []string{"moreThan"}
		}
		return []string{"lessThan"}
	}, flow.Aggregator(func(m map[string][]byte) ([]byte, error) {
		if v, ok := m["moreThan"]; ok {
			return v, nil
		}
		if v, ok := m["lessThan"]; ok {
			return v, nil
		}
		return nil, nil
	}))
	branches["moreThan"].Node("expand-10", Expand10)
	branches["lessThan"].Node("expand-100", Expand100)
	dag.Node("output", Output2)

	dag.Edge("input", "add-ten")
	dag.Edge("input", "add-ten2")
	dag.Edge("add-ten", "aggregator")
	dag.Edge("add-ten2", "aggregator")

	dag.Edge("aggregator", "judgeTen")
	dag.Edge("judgeTen", "output")
	return nil

}

func main() {
	fs := goflow.FlowService{
		Port:              8080,
		RedisURL:          "localhost:6379",
		WorkerConcurrency: 5,
	}

	if err := fs.Register("math2", MyFlow2); err != nil {
		panic(err)

	}
	if err := fs.Start(); err != nil {
		panic(err)
	}
}
