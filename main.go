package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	fmt.Println("Printing test")

	if len(os.Args) < 2 {
		log.Fatal("Error, try someting else")
	}

	buff, err := os.ReadFile(os.Args[1])

	if err != nil {
		log.Fatal("Error. try something else")
	}

	var CpuState = CpuState{Memory: buff}

	for ; int(CpuState.PC) < len(buff); CpuState.PC++ {
		fmt.Printf("0x%x\n", CpuState.RegA)
		UpdateState(&CpuState, buff[CpuState.PC])
	}

	CpuState.Memory = make([]byte, 0)
	fmt.Printf("%+v\n", CpuState)
}
