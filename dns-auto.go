package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os/exec"
	"strings"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func ping(ip string) bool {

	out, _ := exec.Command("ping", ip, "-c 1", "-i 3", "-w 10").Output()

	if strings.Contains(string(out), "100% packet loss") {
		fmt.Println("== Disabling " + ip)
		return false
	} else {
		fmt.Println("== Enabling " + ip)
		return true
	}
}

func fix(ip string, enabled bool, data string) string {

	if enabled {
		return strings.Replace(string(data), "# Unreachable "+ip, "nameserver "+ip, -1)
	} else {
		return strings.Replace(string(data), "nameserver "+ip, "# Unreachable "+ip, -1)
	}
}

func main() {

	flag.Parse()

	path := "/etc/resolv.conf"
	data, err := ioutil.ReadFile(path)
	contents := string(data)
	check(err)
	fmt.Println("== Current /etc/resolv.conf:")
	fmt.Println(contents)

	if len(flag.Args()) < 1 {
		fmt.Println("No IPs provided on command line, exiting.")
		return
	}

	for i, ip := range flag.Args() {
		fmt.Println("== Checking #", i, ip)
		contents = fix(ip, ping(ip), contents)
	}

	fmt.Println()
	fmt.Println("== Modified /etc/resolv.conf:")
	fmt.Println(contents)

	err = ioutil.WriteFile(path, []byte(contents), 0)
	check(err)
}
