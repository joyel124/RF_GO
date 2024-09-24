package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	//"strconv"
	"os"

	"tp-test/RF"

	"time"
	//"math"
)

func main() {

	start := time.Now()
	f, _ := os.Open("diabetes.csv")
	defer f.Close()
	content, _ := ioutil.ReadAll(f)
	s_content := string(content)
	lines := strings.Split(s_content, "\n")

	inputs := make([][]interface{}, 0)
	targets := make([]string, 0)
	for _, line := range lines {

		line = strings.TrimRight(line, "\r\n")

		if len(line) == 0 {
			continue
		}
		tup := strings.Split(line, ",")
		pattern := tup[:len(tup)-1]
		target := tup[len(tup)-1]
		X := make([]interface{}, 0)
		for _, x := range pattern {
			X = append(X, x)
		}
		inputs = append(inputs, X)

		targets = append(targets, target)
	}
	train_inputs := make([][]interface{}, 0)

	train_targets := make([]string, 0)

	test_inputs := make([][]interface{}, 0)
	test_targets := make([]string, 0)

	for i, x := range inputs {
		if i%2 == 1 {
			test_inputs = append(test_inputs, x)
		} else {
			train_inputs = append(train_inputs, x)
		}
	}

	for i, y := range targets {
		if i%2 == 1 {
			test_targets = append(test_targets, y)
		} else {
			train_targets = append(train_targets, y)
		}
	}

	fmt.Println("train size:", len(train_inputs))
	fmt.Println("feature size:", len(train_inputs[0]))
	/*fmt.Println("feature:", train_inputs)
	fmt.Println("Inputs:", inputs)
	fmt.Println("Targets:", targets)*/

	forest := RF.BuildForest(inputs, targets, 10, 500, len(train_inputs[0])) //100 trees

	test_inputs = train_inputs
	test_targets = train_targets
	err_count := 0.0
	for i := 0; i < len(test_inputs); i++ {
		output := forest.Predicate(test_inputs[i])
		expect := test_targets[i]
		fmt.Println(output, expect)
		if output != expect {
			err_count += 1
		}
	}
	fmt.Println("success rate:", 1.0-err_count/float64(len(test_inputs)))

	fmt.Println(time.Since(start))

	//fmt.Println("Test Inputs:", test_inputs[1000])
	//Probar introduciendo un nuevo dato Female,53.0,0,0,former,27.32,7.0,159,1
	//imprimir test
	fmt.Print("Test Input: ")
	fmt.Println("Female, 53.0, 0, 0, former, 27.32, 7.0, 159")
	fmt.Println("Prediction:", forest.Predicate([]interface{}{"Female", "53.0", "0", "0", "former", "27.32", "7.0", "159"}))
	//fmt.Println("Prediction:", forest.Predicate(test_inputs[2]))
}
