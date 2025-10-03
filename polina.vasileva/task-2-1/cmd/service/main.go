package main

import "fmt"

const (
	MINTEMP = 15
	MAXTEMP = 30
)

type ConditionerT struct {
	minTemp int
	maxTemp int
	status  bool
}

func main() {
	var departNum int
	_, err := fmt.Scan(&departNum)

	if err != nil {
		fmt.Println("Invalid input")

		return
	}

	for range departNum {
		var emplCount int

		_, err := fmt.Scan(&emplCount)
		if err != nil {
			fmt.Println("Invalid input")

			return
		}

		conditioner := ConditionerT{MINTEMP, MAXTEMP, true}

		for range emplCount {
			var sign string

			_, err = fmt.Scan(&sign)
			if err != nil {
				fmt.Println("Invalid input")

				return
			}

			var degrees int

			_, err = fmt.Scan(&degrees)
			if err != nil {
				fmt.Println("Invalid input")

				return
			}

			if !conditioner.status {
				continue
			}

			changeTemp(&conditioner, sign, degrees)

			if conditioner.status {
				fmt.Println(conditioner.minTemp)
			} else {
				fmt.Println("-1")

				return
			}
		}
	}
}

func changeTemp(cond *ConditionerT, sign string, degrees int) {
	switch sign {
	case ">=":
		handleGreaterEqual(cond, degrees)
	case "<=":
		handleLessEqual(cond, degrees)
	default:
		cond.status = false
	}
}

func handleGreaterEqual(cond *ConditionerT, degrees int) {
	if degrees < cond.minTemp {
		return
	}

	if degrees > cond.maxTemp {
		cond.status = false

		return
	}

	cond.minTemp = degrees
}

func handleLessEqual(cond *ConditionerT, degrees int) {
	if degrees > cond.maxTemp {
		return
	}

	if degrees < cond.minTemp {
		cond.status = false

		return
	}

	cond.maxTemp = degrees
}
