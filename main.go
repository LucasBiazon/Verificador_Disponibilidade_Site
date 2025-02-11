package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/user"
)

var (
	dirName string = "checker-site"
	dataName string = "data.json"
	responseName string = "response.json"	
)

type CliHandler interface{
	OpenInterface()
	OpenInputFile()
	OpenResponseFile()
}
type Cli struct{}
func (Cli) OpenInterface(){
	fmt.Println("abrindo interface")
}
func (Cli) OpenInputFile(){
}
func (Cli) OpenResponseFile(){
	fmt.Println("abrindo response")
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

type Site struct {
	Name string `json:"name"`
	Url  string `json:"url"`
}

type Sites struct {
	Sites []Site `json:"sites"`
}


type FileSystem interface{
	Stat(name string)(os.FileInfo, error)
	MkdirAll(path string, perm os.FileMode) error
	Create(name string) (*os.File, error)
}

type FileSystemImpl struct{}
func (FileSystemImpl) Stat(name string)(os.FileInfo, error){
	return os.Stat(name)
}
func (FileSystemImpl) MkdirAll(path string, perm os.FileMode) error{
	return os.MkdirAll(path, perm)
}
func (FileSystemImpl) Create(name string) (*os.File, error){
	return os.Create(name)
}

type DirManager struct {
	fs FileSystem
	username  string 

}
func (d *DirManager) CreateDir() {
	dirPath := fmt.Sprintf("/home/%s/%s", d.username, dirName)
	dataPath := fmt.Sprintf("%s/%s", dirPath, dataName)
	responsePath := fmt.Sprintf("%s/%s", dirPath, responseName)

	if _, err := d.fs.Stat(dirPath); os.IsNotExist(err) {
		err := d.fs.MkdirAll(dirPath, os.ModePerm)
		if err != nil {
			log.Fatalf("Erro ao criar o diretório: %v", err)
		}
	}


	if _, err := d.fs.Stat(dataPath); os.IsNotExist(err) {
		file, err := d.fs.Create(dataPath)
		if err != nil {
			log.Fatalf("Erro ao criar o arquivo: %v", err)
		}
		defer file.Close()
		defaultData := Sites{
			Sites: []Site{
				{
					Name: "site example",
					Url:  "https://github.com/lucasBiazon",
				},
			},
		}
		jsonData, err := json.MarshalIndent(defaultData, "", "  ")
		if err != nil {
			log.Fatalf("Erro ao gerar JSON padrão: %v", err)
		}

		if _, err := file.Write(jsonData); err != nil {
			log.Fatalf("Erro ao escrever JSON no arquivo data.json: %v", err)
		}
	}

	if _, err := d.fs.Stat(responsePath); os.IsNotExist(err) {
		file, err := d.fs.Create(responsePath)
		if err != nil {
			log.Fatalf("Erro ao criar o arquivo: %v", err)
		}
		defer file.Close()
	}
}

type FileJsonReader interface {
	ReadJsonFile(path string) ([]byte, error)
}

type FileReaderImpl struct{}

func (FileReaderImpl) ReadJsonFile(path string) ([]byte, error) {
	jsonFile, err := os.OpenFile(path, os.O_RDONLY|os.O_CREATE, 0644)
	if err != nil {
		fmt.Print(err)
		return nil, errors.New(fmt.Sprint("Error: ", err))
	}
	defer jsonFile.Close()
	return io.ReadAll(jsonFile)
}

func ReadDataFile(fr FileJsonReader, path string) (*Sites, error) {
	data, err := fr.ReadJsonFile(path)
	if err != nil {
		return nil, err
	}
	var sites Sites
	if len(data) > 0 {
		if err := json.Unmarshal(data, &sites); err != nil {
			return nil, errors.New(fmt.Sprint("Error: ", err))
		}

		for _, v := range sites.Sites {
			if v.Name == "" {
				return &Sites{}, errors.New("site name is required")
			}
			if v.Url == "" {
				return &Sites{}, errors.New("site url is required")
			}
		}

		return &sites, nil
	}
	return &Sites{}, nil
}


type Checker interface{
	Checker(url string) (bool)
}


type Result struct {
	string
	bool
}

type CheckerImpl struct{}
func (CheckerImpl) Checker(url string) (bool){
	r, err := http.Get(url)
	if err != nil{
		fmt.Println("error when checking the site: ",url)
		return false
	}
	if r.StatusCode != 200{
		return false
	}
	return true
}

func WebSiteChecker(c Checker, sites *Sites) (map[string]bool){
	results := make(map[string]bool)
	canalResult := make(chan Result)

	for _, site := range sites.Sites{
		go func(n, u string){
			canalResult <- Result{n, c.Checker(u)}
		}(site.Name, site.Url)
	}
	for i := 0; i < len(sites.Sites); i++{
		result := <-canalResult
		results[result.string] = result.bool
	}

	return results
}

func main(){
	fs := FileSystemImpl{}
	usr, err := user.Current()
	if err != nil{
		log.Fatalln("Error em pegar dados de usuário: ", err)
	}

	dirManager := DirManager{
		fs:       fs,
		username: usr.Username,
	}
	dirManager.CreateDir()


	cli := &Cli{}
	RunCli(os.Stdin, os.Stdout, cli)
}