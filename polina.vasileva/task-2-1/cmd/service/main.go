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

			changeTemp(&conditioner, sign, degrees)

			if conditioner.status {
				fmt.Println(conditioner.minTemp)
			} else {
				fmt.Println("-1")
			}
		}
	}
}

func changeTemp(cond *ConditionerT, sign string, degrees int) {
	switch sign {
	case ">=":
		if degrees > cond.maxTemp {
			cond.status = false
		}

		if degrees >= cond.minTemp {
			cond.minTemp = degrees
		}

	case "<=":
		if degrees < cond.minTemp {
			cond.status = false
		}

		if degrees <= cond.maxTemp {
			cond.maxTemp = degrees
		}

	default:
		cond.status = false
	}
}
