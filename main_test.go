package main

import (
	"bytes"
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