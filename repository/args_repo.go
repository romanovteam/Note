package repository

import (
	"errors"
	"gorm.io/gorm"
	"time"
)

type Arg struct {
	ID    uint   `gorm:"primary_key"`
	Name  string `gorm:"unique;not null"`
	Notes []Note `gorm:"many2many:arg_notes;constraint:OnDelete:CASCADE;"`
}

type Note struct {
	ID        uint   `gorm:"primary_key"`
	Text      string `gorm:"not null"`
	Args      []*Arg `gorm:"many2many:arg_notes;constraint:OnDelete:CASCADE;"`
	CreatedAt time.Time
}

// ArgsRepo интерфейс для работы с аргументами и заметками
type ArgsRepo interface {
	AddArg(name string) (*Arg, error)
	GetArg(name string) (*Arg, error)
	AddNoteToArgs(argNames []string, text string) error
	GetNotesByArgName(argName string) ([]Note, error)
	GetAllArgs() ([]Arg, error)
	DeleteArg(name string) error
	DeleteAllNotesAndArgs() error
	GetTodayNotesByArgName(argName string) ([]Note, error)
}

// GormArgsRepo реализация интерфейса ArgsRepo с использованием GORM
type GormArgsRepo struct {
	db *gorm.DB
}

// NewGormArgsRepo создаем новый GormArgsRepo
func NewGormArgsRepo(db *gorm.DB) *GormArgsRepo {
	return &GormArgsRepo{db: db}
}

// AddArg добавляем новый аргумент в базу данных
func (g *GormArgsRepo) AddArg(name string) (*Arg, error) {
	arg := &Arg{Name: name}
	err := g.db.FirstOrCreate(arg, Arg{Name: name}).Error
	return arg, err
}

// GetArg достаем аргумент из базы данных по имени
func (g *GormArgsRepo) GetArg(name string) (*Arg, error) {
	var arg Arg
	err := g.db.Where("name = ?", name).First(&arg).Error
	return &arg, err
}

// AddNoteToArgs добавляем новую заметку, связанную с несколькими аргументами
func (g *GormArgsRepo) AddNoteToArgs(argNames []string, text string) error {
	var args []*Arg
	for _, name := range argNames {
		arg, err := g.GetArg(name)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				// Если аргумент не найден, создаем его
				arg, err = g.AddArg(name)
				if err != nil {
					return err
				}
			} else {
				return err
			}
		}
		args = append(args, arg)
	}

	note := &Note{Text: text, Args: args}
	return g.db.Create(note).Error
}

// GetNotesByArgName получает все заметки, связанные с аргументом по его имени
func (g *GormArgsRepo) GetNotesByArgName(argName string) ([]Note, error) {
	var arg Arg
	err := g.db.Where("name = ?", argName).Preload("Notes").First(&arg).Error
	if err != nil {
		return nil, err
	}

	var notes []Note
	err = g.db.Model(&arg).Association("Notes").Find(&notes)
	return notes, err
}

// GetAllArgs достаем все аргументы из базы данных
func (g *GormArgsRepo) GetAllArgs() ([]Arg, error) {
	var args []Arg
	err := g.db.Find(&args).Error
	return args, err
}

// DeleteArg удаляем аргумент из базы данных
func (g *GormArgsRepo) DeleteArg(name string) error {
	return g.db.Where("name = ?", name).Delete(&Arg{}).Error
}

func (g *GormArgsRepo) DeleteAllNotesAndArgs() error {
	var args []Arg

	// Получаем все аргументы
	if err := g.db.Find(&args).Error; err != nil {
		return err
	}

	// Удаляем все ассоциации между аргументами и заметками
	for _, arg := range args {
		if err := g.db.Model(&arg).Association("Notes").Clear(); err != nil {
			return err
		}
	}

	// Затем удаляем все записи из таблицы Notes
	if err := g.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Note{}).Error; err != nil {
		return err
	}

	// Наконец, удаляем все записи из таблицы Args
	if err := g.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&Arg{}).Error; err != nil {
		return err
	}

	return nil
}

func (g *GormArgsRepo) GetTodayNotesByArgName(argName string) ([]Note, error) {
	var notes []Note
	startOfDay := time.Now().Truncate(24 * time.Hour)

	err := g.db.Joins("JOIN arg_notes ON notes.id = arg_notes.note_id").
		Joins("JOIN args ON args.id = arg_notes.arg_id").
		Where("args.name = ?", argName).
		Where("notes.created_at >= ?", startOfDay).
		Find(&notes).Error

	return notes, err
}
