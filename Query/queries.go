package main

import (
	"fmt"

	"database/sql"

	_ "github.com/go-sql-driver/mysql"
)

// Employee is ...
type Employee struct {
	Employee string
}

func getManagers() {
	db, err := sql.Open("mysql", "root:root@/employees")
	checkErr(err)

	query := `SELECT CONCAT(
		' Name - ',m.first_name,' ',m.last_name,' ','Title - ',t.title,' ','Departament - ',d.dept_name,' ','Salary - ',s.salary
		) manager 
			FROM employees m JOIN 
			dept_manager dm ON dm.emp_no = m.emp_no JOIN 
			dept_emp de ON m.emp_no = de.emp_no JOIN 
			departments d ON de.dept_no = d.dept_no AND dm.dept_no = de.dept_no JOIN 
			titles t ON m.emp_no = t.emp_no JOIN 
			salaries s ON m.emp_no = s.emp_no 
			WHERE dm.to_date > NOW() AND t.to_date > NOW() AND s.to_date > NOW() 
			ORDER BY manager;`
	rows, err := db.Query(query)
	checkErr(err)
	employee := Employee{}
	employees := []Employee{}

	for rows.Next() {
		var manager string
		err = rows.Scan(&manager)
		checkErr(err)
		employee.Employee = manager
		employees = append(employees, employee)
		fmt.Println(employee.Employee)
		defer db.Close()
	}
}

func getEmployees() {
	db, err := sql.Open("mysql", "root:root@/employees")
	checkErr(err)

	query := `
	SELECT CONCAT(
		'Congratulation with job anniversary!!!  
		Departament - ',d.dept_name, ', Title - ',t.title, ', 
		Name - ', e.first_name, '  ', e.last_name, ', Hire date - ', 
		e.hire_date, ', Experience - ', DATE_FORMAT(FROM_DAYS(DATEDIFF(NOW(), e.hire_date)),'%Y')+0
		) employee 
			FROM employees e JOIN 
			dept_emp de ON e.emp_no = de.emp_no JOIN 
			departments d ON de.dept_no = d.dept_no AND e.emp_no = de.emp_no JOIN 
			titles t ON e.emp_no = t.emp_no 
			WHERE (de.to_date > NOW() AND t.to_date > NOW() AND
			MONTH(now()) = MONTH(e.hire_date))
			ORDER BY employee;`
	rows, err := db.Query(query)
	checkErr(err)
	employee := Employee{}
	employees := []Employee{}

	for rows.Next() {
		var manager string
		err = rows.Scan(&manager)
		checkErr(err)
		employee.Employee = manager
		employees = append(employees, employee)
		fmt.Println(employee.Employee)
		defer db.Close()
	}
}

func getDepartmens() {
	db, err := sql.Open("mysql", "root:root@/employees")
	checkErr(err)
	query := `
	SELECT CONCAT(
		' Department - ',dept.dept_name, ', Employees_Number - ',COUNT(emp.emp_no), ', Dept_Salary - ', SUM(sal.salary))
		FROM departments dept JOIN
		 dept_emp cdept on dept.dept_no = cdept.dept_no JOIN
		  employees emp on emp.emp_no = cdept.emp_no JOIN
		   salaries sal ON emp.emp_no = sal.emp_no 
		   WHERE cdept.to_date > NOW() AND sal.to_date > NOW() 
		   GROUP BY dept.dept_name 
		   ORDER BY dept.dept_name;`
	rows, err := db.Query(query)
	checkErr(err)
	employee := Employee{}
	employees := []Employee{}

	for rows.Next() {
		var manager string
		err = rows.Scan(&manager)
		checkErr(err)
		employee.Employee = manager
		employees = append(employees, employee)
		fmt.Println(employee.Employee)
		defer db.Close()
	}
}

func main() {
	getManagers()
	getEmployees()
	getDepartmens()
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
