package main

import (
	"bytes"
	"os"
	"strings"
	"testing"
)

type MockHandler struct{
	InterfaceCalled   bool
	InputFileCalled   bool
	ResponseFileCalled bool
}

func (m *MockHandler ) OpenInterface(){
	m.InterfaceCalled = true
}
func (m *MockHandler ) OpenInputFile(){
	m.InputFileCalled = true
}
func (m *MockHandler ) OpenResponseFile(){
	m.ResponseFileCalled = true
}

func TestRunCi(t *testing.T){
	input := strings.NewReader("1\n2\n3\n4\n")
	output := &bytes.Buffer{}
	mockHandler := &MockHandler{}
	RunCli(input, output, mockHandler)

	got := output.String()
	want := []string{
		"Digite 1 - para abrir interface",
		"Digite 2 - para abrir arquivo de entrada",
		"Digite 3 - para abrir arquivo de resposta",
		"->Abrindo interface<-",
		"->Abrindo arquivo de entrada<-",
		"->Abrindo arquivo de resposta<-",
		"Ops... Não há essa opção",
	}
	for _, w := range want {
		if !strings.Contains(got, w) {
			t.Errorf("Esperava encontrar %q na saída, mas não foi encontrado", w)
		}
	}
	if !mockHandler.InterfaceCalled {
		t.Errorf("Esperava que OpenInterface fosse chamada, mas não foi")
	}
	if !mockHandler.InputFileCalled {
		t.Errorf("Esperava que OpenInputFile fosse chamada, mas não foi")
	}
	if !mockHandler.ResponseFileCalled {
		t.Errorf("Esperava que OpenResponseFile fosse chamada, mas não foi")
	}
}
type MockFileSystem struct {
	StatCalls     map[string]bool
	MkdirAllCalls map[string]bool
	CreateCalls   map[string]bool
}
func (m *MockFileSystem) Stat(name string) (os.FileInfo, error) {
	// Simula a existência do diretório/arquivo se já tiver sido criado
	if _, exists := m.StatCalls[name]; exists {
		return nil, nil
	}
	return nil, os.ErrNotExist
}

func (m *MockFileSystem) MkdirAll(path string, perm os.FileMode) error {
	m.MkdirAllCalls[path] = true
	m.StatCalls[path] = true // Marca o diretório como existente
	return nil
}

func (m *MockFileSystem) Create(name string) (*os.File, error) {
	m.CreateCalls[name] = true
	m.StatCalls[name] = true 
	

	tmpFile, err := os.CreateTemp("", "dummy")
	if err != nil {
			return nil, err
	}
	return tmpFile, nil
}
func TestCreateDir(t *testing.T) {
    // Valores esperados
    dirPath := "/home/testuser/checker-site"
    dataPath := "/home/testuser/checker-site/data.json"
    responsePath := "/home/testuser/checker-site/response.json"

    t.Run("Deve criar o diretório se não existir", func(t *testing.T) {
        mockFS := &MockFileSystem{
            StatCalls:     make(map[string]bool),
            MkdirAllCalls: make(map[string]bool),
            CreateCalls:   make(map[string]bool),
        }
        // Injeta o username esperado
        dirManager := &DirManager{fs: mockFS, username: "testuser"}

        dirManager.CreateDir()

        if !mockFS.MkdirAllCalls[dirPath] {
            t.Errorf("Esperava que o diretório %s fosse criado, mas não foi chamado", dirPath)
        }
    })

    t.Run("Não deve recriar o diretório se já existir", func(t *testing.T) {
        mockFS := &MockFileSystem{
            StatCalls:     map[string]bool{dirPath: true}, // Diretório já existe
            MkdirAllCalls: make(map[string]bool),
            CreateCalls:   make(map[string]bool),
        }
        dirManager := &DirManager{fs: mockFS, username: "testuser"}

        dirManager.CreateDir()

        if mockFS.MkdirAllCalls[dirPath] {
            t.Errorf("O diretório %s já existia, mas MkdirAll foi chamado erroneamente", dirPath)
        }
    })

    t.Run("Deve criar os arquivos se não existirem", func(t *testing.T) {
        mockFS := &MockFileSystem{
            StatCalls:     make(map[string]bool),
            MkdirAllCalls: make(map[string]bool),
            CreateCalls:   make(map[string]bool),
        }
        dirManager := &DirManager{fs: mockFS, username: "testuser"}

        dirManager.CreateDir()

        if !mockFS.CreateCalls[dataPath] {
            t.Errorf("Esperava que o arquivo %s fosse criado, mas não foi chamado", dataPath)
        }
        if !mockFS.CreateCalls[responsePath] {
            t.Errorf("Esperava que o arquivo %s fosse criado, mas não foi chamado", responsePath)
        }
    })

    t.Run("Não deve recriar arquivos se já existirem", func(t *testing.T) {
        mockFS := &MockFileSystem{
            StatCalls: map[string]bool{
                dataPath:     true, // Arquivo já existe
                responsePath: true, // Arquivo já existe
            },
            MkdirAllCalls: make(map[string]bool),
            CreateCalls:   make(map[string]bool),
        }
        dirManager := &DirManager{fs: mockFS, username: "testuser"}

        dirManager.CreateDir()

        if mockFS.CreateCalls[dataPath] {
            t.Errorf("O arquivo %s já existia, mas Create foi chamado erroneamente", dataPath)
        }
        if mockFS.CreateCalls[responsePath] {
            t.Errorf("O arquivo %s já existia, mas Create foi chamado erroneamente", responsePath)
        }
    })
}
