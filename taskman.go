package main

import (
	"bufio"
	"fmt"
	//"io";
	"flag"
	"io/ioutil"
	"time"
	//"encoding/binary"
	"os"
	"os/user"
	"strconv"
	"strings"
)

// Reading files requires checking most calls for errors. This helper will streamline our error checks below.

func check(e error) {
	if e != nil {
		panic(e)
	}
}

const name string = "Cristian"
const quote string = "HARD WORK BEATS TALENT"

type task struct {
	daily       bool
	mandatory   bool
	due_date    string
	date_added  string
	description string
}

// Checks if the program has been run today
// Initializes the main slice
func initi(dailies_cache, tasks_cache, last_run_cache, today_cache string) (int, []task) {

	current_date := time.Now().Format("20060102")
	byte_array, _ := ioutil.ReadFile(last_run_cache)
	last_run_date := string(byte_array)
	today_bin, _ := ioutil.ReadFile(today_cache)
	today, _ := strconv.Atoi(string(today_bin))
	var dailies []task

	if current_date == last_run_date {
	} else {
		today = 0
		fmt.Println("new day")
		dailies = load_csv(dailies_cache)
		ioutil.WriteFile(last_run_cache, []byte(current_date), 0644)
		ioutil.WriteFile(today_cache, []byte(string(today)), 0644)
	}

	return today, dailies

}

// returns an array that is or is not daily
func array_grep(task_list []task, value bool) []task {
	var final []task
	for i, j := range task_list {

		_, _, day := time.Now().Date()

		should_write := true

		if j.due_date != "everyday" {
			isdaily := j.daily
			taskiseven := j.due_date == "even"
			dayiseven := day%2 == 0

			should_write = !isdaily || (!taskiseven && !dayiseven) || (taskiseven && dayiseven)
		}

		matches_input := j.daily == value
		//fmt.Println(should_write,j)

		if should_write && matches_input {
			fmt.Print("[#", i, "] ")
			if j.mandatory == true {
				fmt.Print("! ")
			} else {
				fmt.Print("  ")
			}
			final = append(final, j)
			fmt.Print(j.description, " (", j.due_date, ")")
			fmt.Print("\n")
		}
	}

	fmt.Print("\n")
	return final
}

// Presents the task list
func present(task_list []task, today int, username string) {
	// TODO: Find out how to clear screen

	fmt.Print("Hello ", username, ". Today you completed ", today, " tasks.\n")
	//fmt.Println("\n\n",quote)
	fmt.Println("\nDailies\n=======\n")
	array_grep(task_list, true)
	fmt.Println("Non-dailies\n=======\n")
	array_grep(task_list, false)

}

func do_task(a []task, i int) []task {
	// TODO: Log done tasks to a file
	a = append(a[:i], a[i+1:]...)
	return a
}

func add_task() task {
	var daily bool
	var mandatory bool
	var due_date string
	var task_name string

	fmt.Print("daily? (y/N): ")
	input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	if input == "y\n" {
		//fmt.Println("user has input true")
		daily = true
	} else {
		//fmt.Println("user has input false")
		daily = false
	}

	if daily == true {
		fmt.Print("days (Mon,Tue,Wed): ")
	} else {
		fmt.Print("mandatory? (Y/n): ")
	}

	input, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	if input == "n\n" {
		//fmt.Println("user has input false")
		mandatory = false
	} else {
		//fmt.Println("user has input true")
		mandatory = true
	}

	fmt.Println("due date (yyy-mm-dd): ")
	due_date, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	if due_date == "\n" {
		due_date = "no date"
	}
	due_date = strings.TrimSuffix(due_date, "\n")

	fmt.Println("task name: ")
	task_name, _ = bufio.NewReader(os.Stdin).ReadString('\n')
	task_name = strings.TrimSuffix(task_name, "\n")

	result := task{daily, mandatory, due_date, "defaultDate", task_name}
	return result
}

func read_action() {
	for {
		fmt.Print("> ")
		action, _ := bufio.NewReader(os.Stdin).ReadString('\n')
		fmt.Println(action)
		//if action == "add\n"
	}
}

func tidy(S []task) []task {

	length := len(S) - 2

	for {
		inorder := true
		for i := 0; i <= length; i++ {
			if S[i].date_added > S[i+1].date_added {
				small := S[i+1]
				big := S[i]
				S[i] = small
				S[i+1] = big
				inorder = false
			}
		}

		if inorder == true {
			break
		}
	}

	for {
		inorder := true
		for i := 0; i <= length; i++ {
			if S[i].due_date > S[i+1].due_date {
				small := S[i+1]
				big := S[i]
				S[i] = small
				S[i+1] = big
				inorder = false
			}
		}

		if inorder == true {
			break
		}
	}

	for {
		inorder := true
		for i := 0; i <= length; i++ {
			a := S[i].mandatory
			b := S[i+1].mandatory
			result := !a && b
			if result {
				small := S[i+1]
				big := S[i]
				S[i] = small
				S[i+1] = big
				inorder = false
			}
		}

		if inorder == true {
			break
		}
	}

	for {
		inorder := true
		for i := 0; i <= length; i++ {
			a := S[i].daily
			b := S[i+1].daily
			result := !a && b
			if result {
				small := S[i+1]
				big := S[i]
				S[i] = small
				S[i+1] = big
				inorder = false
			}
		}

		if inorder == true {
			break
		}
	}

	return S
}

func load_csv(tasks_cache string) []task {

	var task_list []task

	f, _ := os.Open(tasks_cache)
	defer f.Close()

	var lines []string

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	for _, j := range lines {
		daily := false
		daily_raw := strings.Split(j, ",")[0]
		if daily_raw == "true" {
			daily = true
		}

		mandatory := false
		mandatory_raw := strings.Split(j, ",")[1]
		if mandatory_raw == "true" {
			mandatory = true
		}

		due_date := strings.Split(j, ",")[2]

		date_added := strings.Split(j, ",")[3]

		description := strings.Split(j, ",")[4]

		slice := task{daily, mandatory, due_date, date_added, description}

		task_list = append(task_list, slice)
	}

	return task_list
}

func save_csv(task_list []task, tasks_cache, today_cache string, today int) []task {

	f, _ := os.Create(tasks_cache)
	defer f.Close()

	var list []string
	var new_task_list []task

	for index, value := range task_list {

		write := true

		for j := range list {
			if value.description == list[j] {
				write = false
			}
		}

		if write == true {
			if index != 0 {
				fmt.Fprint(f, "\n")
			}

			fmt.Fprint(f, value.daily, ",", value.mandatory, ",", value.due_date, ",", value.date_added, ",", value.description)
			list = append(list, value.description)
			new_task_list = append(new_task_list, value)
		}
	}
	ioutil.WriteFile(today_cache, []byte(strconv.Itoa(today)), 0644)

	return new_task_list
}

func main() {

	user, _ := user.Current()

	tasks_cache := strings.Join([]string{"/home/", string(user.Username), "/.cache/gtm/tasks.csv"}, "")
	dailies_cache := strings.Join([]string{"/home/", string(user.Username), "/.cache/gtm/dailies.csv"}, "")
	last_run_cache := strings.Join([]string{"/home/", string(user.Username), "/.cache/gtm/last_run_date.txt"}, "")
	today_cache := strings.Join([]string{"/home/", string(user.Username), "/.cache/gtm/today.txt"}, "")
	task_list := load_csv(tasks_cache)

	today, dailies := initi(dailies_cache, tasks_cache, last_run_cache, today_cache)

	task_list = append(task_list, dailies...)

	wordPtr := flag.String("word", "no date", "a string")
	numPtr := flag.Int("num", 0, "the id")
	flag.Parse()

	if *wordPtr == "add" {
		task_list = append(task_list, add_task())
	} else if *wordPtr == "do" {
		task_list = do_task(task_list, *numPtr)
		today = today + 1
	} else if *wordPtr == "not" {
		task_list = do_task(task_list, *numPtr)
	} else if *wordPtr == "show" {
	} else if *wordPtr == "next" {
		j := task_list[0]
		if j.mandatory == true {
			fmt.Print("! ")
		}
		fmt.Print(j.description, " (", j.due_date, ")")
		fmt.Print("\n")
	} else {
		fmt.Println("command not recognized")
	}

	if *wordPtr != "next" {
		task_list = tidy(task_list)
		task_list = save_csv(task_list, tasks_cache, today_cache, today)
		present(task_list, today, user.Name)
	}
}
