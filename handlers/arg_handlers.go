package handlers

import (
	"fmt"
	"notebook/service"
	"strings"
)

// ArgsHandler структура для обработки аргументов командной строки
type ArgsHandler struct {
	ArgNames []string
	NoteText string
	Today    bool
}

// NewArgsHandler создает новый ArgsHandler и парсит аргументы командной строки
func NewArgsHandler(args []string) *ArgsHandler {
	var argNames []string
	var noteText string
	today := false

	for i, arg := range args {
		if arg == "today" {
			today = true
			continue
		}

		if arg == "--" {
			// Все последующие аргументы считаем текстом заметки
			noteText = strings.Join(args[i+1:], " ")
			break
		}

		// До разделителя -- все считаем аргументами
		argNames = append(argNames, arg)
	}

	// Если текст заметки пустой, но есть больше одного аргумента, считаем последний аргумент текстом заметки
	if noteText == "" && len(argNames) > 1 {
		noteText = argNames[len(argNames)-1]
		argNames = argNames[:len(argNames)-1]
	}

	return &ArgsHandler{
		ArgNames: argNames,
		NoteText: noteText,
		Today:    today,
	}
}

func (h *ArgsHandler) HandleArgs(flagService *service.ArgsService) error {
	if h.Today && len(h.ArgNames) > 0 {
		return flagService.ProcessGetTodayNotesByArg(h.ArgNames[0])
	}

	if len(h.ArgNames) > 0 && h.NoteText != "" {
		return flagService.ProcessSaveNoteToArgs(h.ArgNames, h.NoteText)
	}

	if len(h.ArgNames) > 0 && h.NoteText == "" {
		return flagService.ProcessGetNotesByArg(h.ArgNames[0])
	}

	return fmt.Errorf("аргументы и текст не могут быть пустыми")
}
