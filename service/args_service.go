package service

import (
	"errors"
	"fmt"
	"notebook/logger"
	"notebook/repository"
)

type ArgsService struct {
	repo   repository.ArgsRepo
	logger logger.Logger
}

// NewArgsService создаем новый ArgsService
func NewArgsService(repo repository.ArgsRepo, logger logger.Logger) *ArgsService {
	return &ArgsService{repo: repo, logger: logger}
}

// ProcessSaveNoteToArgs сохраняет заметку и связывает ее с несколькими аргументами
func (s *ArgsService) ProcessSaveNoteToArgs(argNames []string, text string) error {
	if len(argNames) == 0 || text == "" {
		return errors.New("аргументы и текст не могут быть пустыми")
	}

	if text == "all" {
		for _, argName := range argNames {
			err := s.ProcessGetNotesByArg(argName)
			if err != nil {
				return err
			}
		}
		return nil
	}

	err := s.repo.AddNoteToArgs(argNames, text)
	if err != nil {
		s.logger.LogError(err)
		return err
	}

	fmt.Printf("Заметка '%s' добавлена к аргументам: %v\n", text, argNames)
	return nil
}

// ProcessGetNotesByArg выводит все заметки, связанные с конкретным аргументом
func (s *ArgsService) ProcessGetNotesByArg(argName string) error {
	notes, err := s.repo.GetNotesByArgName(argName)
	if err != nil {
		s.logger.LogError(err)
		return err
	}

	if len(notes) == 0 {
		fmt.Printf("Заметки, связанные с аргументом '%s', отсутствуют.\n", argName)
		return nil
	}

	fmt.Printf("Заметки, связанные с аргументом '%s':\n", argName)
	for _, note := range notes {
		fmt.Printf("- %s\n", note.Text)
	}
	return nil
}

func (s *ArgsService) ProcessDeleteAll() error {
	err := s.repo.DeleteAllNotesAndArgs()
	if err != nil {
		s.logger.LogError(err)
		return err
	}

	fmt.Println("Все заметки и аргументы удалены.")
	return nil
}

func (s *ArgsService) ProcessGetTodayNotesByArg(argName string) error {
	notes, err := s.repo.GetTodayNotesByArgName(argName)
	if err != nil {
		s.logger.LogError(err)
		return err
	}

	if len(notes) == 0 {
		fmt.Printf("Заметки, связанные с аргументом '%s', сделанные сегодня, отсутствуют.\n", argName)
		return nil
	}

	fmt.Printf("Заметки, связанные с аргументом '%s', сделанные сегодня:\n", argName)
	for _, note := range notes {
		fmt.Printf("- %s\n", note.Text)
	}
	return nil
}
