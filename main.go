package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)
type CliHandler interface{
	OpenInterface()
	OpenInputFile()
	OpenResponseFile()
}
func Espace(output io.Writer){
	fmt.Fprint(output, "\n\n\n")
}

func RunCli(input io.Reader, output io.Writer, handler CliHandler){
	scanner := bufio.NewScanner(input)
	fmt.Fprintln(output, "Enter some text (press Ctrl+D or Ctrl+Z to end):")
	fmt.Fprintln(output, "\t Digite 1 - para abrir interface\n" +
							"\t Digite 2 - para abrir arquivo de entrada\n" +
							"\t Digite 3 - para abrir arquivo de resposta")
	fmt.Print("\nYour choice: ")

	for scanner.Scan(){
		text := scanner.Text()
		Espace(output)
		switch text{
		case "1":
			fmt.Fprintln(output,"->Abrindo interface<-")
			handler.OpenInterface()
		case "2":
			fmt.Fprintln(output,"->Abrindo arquivo de entrada<-")
			handler.OpenInputFile()
		case "3":
			fmt.Fprintln(output,"->Abrindo arquivo de resposta<-")
			handler.OpenResponseFile()
		default:
			fmt.Fprintln(output,"Ops... Não há essa opção")
			fmt.Fprintln(output,"\t Digite 1 - para abrir interface\n" +
									"\t Digite 2 - para abrir arquivo de entrada\n" +
									"\t Digite 3 - para abrir arquivo de resposta")
		}
		Espace(output)
		fmt.Print("\nYour new choice: ")

		if err := scanner.Err(); err != nil {
			fmt.Println("Error:", err)
		}
	}
}

type Cli struct{}

func (Cli) OpenInterface(){
	fmt.Println("abrindo interface")
}

func (Cli) OpenInputFile(){
	fmt.Println("abrindo input")
}

func (Cli) OpenResponseFile(){
	fmt.Println("abrindo response")
}

func main(){
	cli := &Cli{}
	RunCli(os.Stdin, os.Stdout, cli)
}