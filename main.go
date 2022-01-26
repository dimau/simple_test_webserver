package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"text/template"
)

func firstThree(s string) string {
	s = strings.TrimSpace(s)
	s = s[:3]
	return s
}

var fm = template.FuncMap{
	"uc": strings.ToUpper,
	"ft": firstThree,
}

type country struct {
	Name    string
	Capital string
}

type dataForTemplate struct {
	Countries    []country
	HeaderString string
}

func handler(conn net.Conn, tpl *template.Template, data dataForTemplate) {
	defer conn.Close()

	// read request
	request(conn)

	// write response
	respond(conn, tpl, data)
}

func request(conn net.Conn) {
	i := 0
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		ln := scanner.Text()
		fmt.Println(ln)
		if i == 0 {
			method, path := strings.Fields(ln)[0], strings.Fields(ln)[1]
			fmt.Println("***METHOD: ", method, " PATH: ", path)
		}
		if ln == "" {
			// headers are done
			break
		}
		i++
	}
}

func respond(conn net.Conn, tpl *template.Template, data dataForTemplate) {
	fmt.Fprint(conn, "HTTP/1.1 200 OK\r\n")
	fmt.Fprint(conn, "Content-Type: text/html\r\n")
	fmt.Fprint(conn, "\r\n")
	// fmt.Fprintln(conn, "Content-Length: "+strconv.Itoa(len(body)))
	err := tpl.ExecuteTemplate(conn, "tpl.gohtml", data)
	if err != nil {
		log.Fatalln(err.Error())
	}
}

func main() {

	// prepairing template container
	tpl, err := template.New("").Funcs(fm).ParseGlob("templates/*.gohtml")
	if err != nil {
		log.Fatalln(err.Error())
	}

	// prepairing data for the HTML template
	russia := country{Name: "Russia", Capital: "Moscow"}
	germany := country{Name: "Germany", Capital: "Berlin"}
	countries := []country{russia, germany}
	data := dataForTemplate{
		Countries:    countries,
		HeaderString: "This is a header string",
	}

	// running TCP server
	li, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalln(err.Error())
	}
	defer li.Close()

	for {
		conn, err := li.Accept()
		if err != nil {
			log.Fatalln(err.Error())
		}

		// each TCP request will be handled by distinct goroutine
		go handler(conn, tpl, data)
	}

}
