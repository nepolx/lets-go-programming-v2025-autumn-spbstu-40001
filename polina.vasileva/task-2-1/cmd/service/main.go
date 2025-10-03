package main

import "fmt"

const MINTEMP = 15
const MAXTEMP = 30

type ConditionerT struct {
	minTemp int
	maxTemp int
	status bool //0 - bad, 1 - good
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
			var degrees int
			_, err = fmt.Scan(&sign)
			if err != nil {
				fmt.Println("Invalid input")
				return
			}
			_, err = fmt.Scan(&degrees)
			if err != nil {
				fmt.Println("Invalid input")
				return
			}
			if conditioner.status {
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
}

func changeTemp(cond *ConditionerT, sign string, degrees int) {
	switch sign {
	case ">=":
		if degrees >= cond.minTemp {
			if degrees <= cond.maxTemp {
				cond.minTemp = degrees
			} else {
			if degrees < cond.minTemp {
				cond.status = true
			} else {
				cond.status = false
			}
		}
	}
	case "<=":
		if degrees <= cond.maxTemp {
			if degrees >= cond.minTemp {
				cond.maxTemp = degrees
			} else {
			if degrees > cond.maxTemp {
				cond.status = true
			} else {
				cond.status = false
			}
		}
	}
	default:
		cond.status = false
	}
}