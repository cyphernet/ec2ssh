package ui

import (
	"text/tabwriter"
	"os"
	"fmt"
	"strconv"
	"bufio"
	"strings"

	"github.com/cyphernet/ec2ssh/instance"
)

func Show(ec2Instances instance.Instances) {
	w := new(tabwriter.Writer)
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	fmt.Fprintln(w, "Index\tName\tIP\tInstance Type\tState")
	for i, c := range ec2Instances {
		fmt.Fprintln(w, strconv.Itoa(i)+"\t"+c.Name+"\t"+c.IP+"\t"+c.InstanceType+"\t"+c.State)
	}
	w.Flush()
}

func Get() int {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Enter index: ")
	text, _ := reader.ReadString('\n')
	i, err := strconv.Atoi(strings.TrimSpace(text))
	if err != nil {
		fmt.Println("Invalid entry")
		os.Exit(0)
	}

	return i
}
