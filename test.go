package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	//"strconv"
	"os"

	"tp-test/RF"

	"time"
	//"math"
)

func main() {

	start := time.Now()
	f, _ := os.Open("diabetesV3.csv")
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

	forest := RF.BuildForest(inputs, targets, 10, 500, len(train_inputs[0])) //100 trees

	test_inputs = train_inputs
	test_targets = train_targets
	err_count := 0.0
	for i := 0; i < len(test_inputs); i++ {
		output := forest.Predicate(test_inputs[i])
		expect := test_targets[i]
		//fmt.Println(output, expect)
		if output != expect {
			err_count += 1
		}
	}
	fmt.Println("success rate:", 1.0-err_count/float64(len(test_inputs)))

	fmt.Println(time.Since(start))

	// Ejecutar el menú para ingresar datos manualmente
	for {
		fmt.Println("---- Menú de Predicción de Diabetes ----")
		fmt.Println("1. Ingresar datos manualmente")
		fmt.Println("2. Salir")
		fmt.Print("Seleccione una opción: ")

		reader := bufio.NewReader(os.Stdin)
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		if choice == "1" {
			input := inputData()
			fmt.Println("Datos ingresados:", input)
			prob := forest.Predicate(input)
			fmt.Printf("Predicción: %s\n", prob)
		} else if choice == "2" {
			break
		} else {
			fmt.Println("Opción inválida, intente nuevamente.")
		}
	}
	//imprimir test
	/*fmt.Print("Test Input (Not diabetic): ")
	fmt.Println("Male, 15.0, 0, 0, never, 30.36, 6.1, 200")
	fmt.Println("Prediction:", forest.Predicate([]interface{}{"Male", "15.0", "0", "0", "never", "30.36", "6.1", "200"}))

	fmt.Print("Test Input (Diabetic): ")
	fmt.Println("Female, 53.0, 0, 0, former, 27.32, 7.0, 159")
	fmt.Println("Prediction:", forest.Predicate([]interface{}{"Female", "53.0", "0", "0", "former", "27.32", "7.0", "159"}))
	//fmt.Println("Prediction:", forest.Predicate(test_inputs[2]))*/
}

// Función para ingresar datos manualmente
func inputData() []interface{} {
	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Ingrese género (F/M): ")
	genderInput, _ := reader.ReadString('\n')
	genderInput = strings.TrimSpace(genderInput)

	var gender string
	if strings.ToUpper(genderInput) == "F" {
		gender = "Female"
	} else if strings.ToUpper(genderInput) == "M" {
		gender = "Male"
	}

	fmt.Print("Ingrese edad: ")
	ageStr, _ := reader.ReadString('\n')
	age := strings.TrimSpace(ageStr)

	fmt.Print("Ingrese hipertensión (0 o 1): ")
	hypertensionStr, _ := reader.ReadString('\n')
	hypertension := strings.TrimSpace(hypertensionStr)

	fmt.Print("Ingrese enfermedad cardíaca (0 o 1): ")
	heartDiseaseStr, _ := reader.ReadString('\n')
	heartDisease := strings.TrimSpace(heartDiseaseStr)

	fmt.Print("Ingrese historial de fumar (0: Nunca, 1: Sin informacion, 2: Actualmente, 3: No actualmente, 4: Anteriormente, 5: Alguna vez): ")
	smokingHistoryStr, _ := reader.ReadString('\n')
	smokingHistory, _ := strconv.Atoi(strings.TrimSpace(smokingHistoryStr))

	var smokingHistoryEng string
	switch smokingHistory {
	case 0:
		smokingHistoryEng = "never"
	case 1:
		smokingHistoryEng = "No Info"
	case 2:
		smokingHistoryEng = "current"
	case 3:
		smokingHistoryEng = "not current"
	case 4:
		smokingHistoryEng = "former"
	case 5:
		smokingHistoryEng = "ever"
	}

	fmt.Print("Ingrese BMI: ")
	bmiStr, _ := reader.ReadString('\n')
	bmi := strings.TrimSpace(bmiStr)

	fmt.Print("Ingrese nivel HbA1c: ")
	hba1cLevelStr, _ := reader.ReadString('\n')
	hba1cLevel := strings.TrimSpace(hba1cLevelStr)

	fmt.Print("Ingrese nivel de glucosa en sangre: ")
	bloodGlucoseLevelStr, _ := reader.ReadString('\n')
	bloodGlucoseLevel := strings.TrimSpace(bloodGlucoseLevelStr)

	// Retorna los datos como un slice de interface{}
	return []interface{}{
		gender,            // Gender (Female/Male)
		age,               // Age
		hypertension,      // Hypertension (0 o 1)
		heartDisease,      // Heart Disease (0 o 1)
		smokingHistoryEng, // Smoking History (never/No Info/current/not current/former/ever)
		bmi,               // BMI
		hba1cLevel,        // HbA1c Level
		bloodGlucoseLevel, // Blood Glucose Level
	}
}
