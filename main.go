package main

import (
	"bufio"
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/go-sql-driver/mysql"
	"github.com/olekukonko/tablewriter"
)




var db *sql.DB


//hodling row data in db using struct

type Tasks struct{
  ID int64
  Name string
  Status bool
}

//mark true the row user specifies for now
func markCompleted(id int64) {
  _, err := db.Exec("UPDATE tasks SET status = ? WHERE id = ?", true, id)
  if err != nil{
    log.Fatal(err)
  }
  fmt.Println("Updated successfully")
  
}

func addNewTasks(name string) (int64, error) {
  var t Tasks
  t.Name = name
  res , err := db.Exec("INSERT INTO tasks (name, status) VALUES (?, ?)", t.Name, false)
  if err != nil{
    return 0, fmt.Errorf("add tasks: %v", err)
  }
  id , err := res.LastInsertId()
  if err != nil{
    return 0, fmt.Errorf("add tasks: %v", err)
  }
  return id, nil
}

func getTaskItems(id int64) (Tasks, error) {
  var t Tasks

  row := db.QueryRow("SELECT * FROM tasks WHERE id = ?", id)
  if err := row.Scan(&t.ID, &t.Name, &t.Status); err != nil{
    //if query was successful but return empty set
    if err == sql.ErrNoRows{
      return t, fmt.Errorf("no such id : %d", id)
    }
    return t, fmt.Errorf("no such row : %d", id)
  }
  return t, nil

}


func wholeDbFromLastEntry() (string, error) {
  if db == nil{
    log.Fatal("Db connection not initialized")
  }
  var tasks []Tasks
  rows, err  := db.Query("SELECT * FROM tasks ORDER BY id DESC")
  if  err != nil{
    return "", fmt.Errorf("wholeDbFromLastEntry: %v", err)
  } 
  defer rows.Close()

  for rows.Next(){
    var t Tasks
    if err := rows.Scan(&t.ID, &t.Name, &t.Status); err != nil{

      return "", fmt.Errorf("wholeDbFromLastEntry row loop: %v", err)
    }
    tasks = append(tasks, t)

  }
  // if errors occured during iterations
  if err := rows.Err(); err != nil{
    return "", fmt.Errorf("error during iterations : %v", err)

  }
  var buf bytes.Buffer
  table := tablewriter.NewWriter(&buf)

  table.SetHeader([]string{"Task Name", "Completed"})

  for t, s := range tasks{
    if s.Name == ""{
      continue
    }
    
    table.Append([]string{s.Name, fmt.Sprintf("%v", s.Status) })
    fmt.Println("t", t)
    fmt.Println("s", s)
  }
  table.Render()
  return buf.String(), nil
}
   



func main(){

    //connect to DB
  cfg := mysql.Config{
    User : "denzil",
    Passwd : "pass",
    Net : "tcp",
    Addr : "127.0.0.1:3306",
    DBName: "tasklists",
    AllowNativePasswords: true,
    
  }
  var er error
  db, er = sql.Open("mysql", cfg.FormatDSN())
  if er != nil {
    log.Fatal(er)
  }

  pingErr := db.Ping()
  if pingErr != nil {
    log.Fatal(pingErr)
  }
  fmt.Println("Connected")


  // Have a better system 
  //for now whether user want to create or update db
  reader := bufio.NewReader(os.Stdin)
  // add cool tasks --help thing
  //NOTE user input for displaying usage options 
  fmt.Print("Enter tasks --help or y/n to use the program thank you: ")
  val, _ := reader.ReadString('\n')
  val = strings.TrimSpace(val)
  fmt.Println("val", val)
    taksGuide := `
    Welcome to tasks a TODO CLI app made for programmers.
    Developer : Denzil Pinto 
    Usage: 
      e - Exit the program
      c - Create your tasks
      u - Update or delete your tasks
      v - View your tasks
    Request: Use lower case alphabets only.
    Complaints:
      visit github.com/icodeologist/tasksCLI-IO
    `
  if val == "y"{
    fmt.Println(taksGuide)
  }else{
    fmt.Println("Program Exited")
    return

  }
  //Get user input to create tasks
  //NOTE user input for Usage optiongs
  readerForCRUD := bufio.NewReader(os.Stdin)
  fmt.Print("Enter the valid usage options : ")
  valAgain , _ := readerForCRUD.ReadString('\n')
  valAgain = strings.TrimSpace(valAgain)

  //for create or c
  switch valAgain {
  case "c":
    // user input here again for task name
    //IF user enter c then we need to get the name of tasks that the 
    //user want to create"
    //NOTE Probably last one for task names
    r := bufio.NewReader(os.Stdin)
    fmt.Print("Enter the task you want to create: ")
    taskNameFromUser , _ := r.ReadString('\n')
    taskNameFromUser = strings.TrimSpace(taskNameFromUser)

    //once we get the input
    //call create function
    id, err := addNewTasks(taskNameFromUser)
    if err != nil{
      fmt.Println(err)
    }else{
      fmt.Println("Task added successfully")
    }
    //show db after successfully getting id back
    data , err  := getTaskItems(id)
    if err != nil{
      fmt.Println("oops : ", err) 
    }else{
      fmt.Println(data)
    }
  // view entire db by recently added data
   case "v": 
    data, err := wholeDbFromLastEntry()
    if err != nil{
      fmt.Errorf("error wholeDbFromLastEntry : %v" ,err)
    }else{
    fmt.Println(data)
    }
  case "u", "d":
    data, err := wholeDbFromLastEntry()
    if err != nil{
      fmt.Errorf("error wholeDbFromLastEntry : %v" ,err)
    }else{
    fmt.Println(data)
    }
    id := bufio.NewReader(os.Stdin)
    fmt.Print("Enter the ID you want to mark completed :")
    idTOedit , _ := id.ReadString('\n')
    idTOedit = strings.TrimSpace(idTOedit)
    num , err := strconv.Atoi(idTOedit)
    if err != nil{
      log.Fatal(err)
    }
    markCompleted(int64(num))
    
    updatedData, err := wholeDbFromLastEntry()
    if err != nil{
      log.Fatal(err)
    }
    fmt.Println(updatedData)
  default: 
    fmt.Println("Tafkdkjsfkdsj")

  }




}




//TODO 
//Return the db with nice formatted way Idk json maybe
//Fix some logic while user experience thats it
