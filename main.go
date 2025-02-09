package main

import (
	"bufio"
	"fmt"
	"os"
)


func Espace(){
	fmt.Print("\n\n\n")
}

func main(){
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("Enter some text (press Ctrl+D or Ctrl+Z to end):")
	fmt.Println("\t Digite 1 - para abrir interface\n" +
							"\t Digite 2 - para abrir arquivo de entrada\n" +
							"\t Digite 3 - para abrir arquivo de resposta")
	fmt.Print("Your choice: ")
	for scanner.Scan(){
		text := scanner.Text()

		Espace()

		switch text{
			case "1":
				fmt.Println(" -Iniciando server...")
				fmt.Println(" -Abrindo interface...")
			case "2":
				fmt.Println(" -Abrindo arquivo de entrada...")
			case "3": 
				fmt.Println(" -Abrindo arquivo de respostas...")
			default:
				fmt.Println("Ops... Não há essa opção")
				fmt.Println("\t Digite 1 - para abrir interface\n" +
							"\t Digite 2 - para abrir arquivo de entrada\n" +
							"\t Digite 3 - para abrir arquivo de resposta")
		}

		Espace()
		fmt.Print("Your new choice: ")

		if err := scanner.Err(); err != nil {
			fmt.Println("Error:", err)
		}
	}
}