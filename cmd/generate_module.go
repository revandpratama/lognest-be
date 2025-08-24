package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/rs/zerolog/log"
)

type TemplateData struct {
	OriginalModuleName    string
	ModuleName            string 
	CapitalizedModuleName string 
	GoModulePath          string 
}

func GenerateModule(moduleName string) {
	const goModulePath = "github.com/revandpratama/lognest"

	sanitizedModuleName := strings.ReplaceAll(moduleName, "-", "")

	// e.g., "user-profile" -> ["user", "profile"] -> ["User", "Profile"] -> "UserProfile"
	parts := strings.Split(moduleName, "-")
	var capitalizedParts []string
	for _, part := range parts {
		if len(part) > 0 {
			capitalizedParts = append(capitalizedParts, strings.ToUpper(string(part[0]))+part[1:])
		}
	}
	capitalizedModuleName := strings.Join(capitalizedParts, "")

	data := TemplateData{
		OriginalModuleName:    moduleName,
		ModuleName:            sanitizedModuleName,
		CapitalizedModuleName: capitalizedModuleName,
		GoModulePath:          goModulePath,
	}

	log.Info().Msgf("Generating module: %s\n", moduleName)

	// Define the directory and file structure
	structure := map[string]string{
		"entity":     entityTemplate,
		"repository": repositoryTemplate,
		"usecase":    usecaseTemplate,
		"handler":    handlerTemplate,
	}

	for dir, fileTemplate := range structure {
		// Create directory
		fullDirPath := filepath.Join("internal", "modules", moduleName, dir)
		if err := os.MkdirAll(fullDirPath, os.ModePerm); err != nil {
			log.Fatal().Err(err).Msgf("Failed to create directory %s: %v", fullDirPath, err)
		}

		// Create file
		fileName := fmt.Sprintf("%s.go", dir)
		fullFilePath := filepath.Join(fullDirPath, fileName)
		if err := createFileFromTemplate(fullFilePath, fileTemplate, data); err != nil {
			log.Fatal().Err(err).Msgf("Failed to create file %s: %v", fullFilePath, err)
		}
		log.Printf("  ✓ Created %s", fullFilePath)
	}

	log.Info().Msg("\nModule generation complete! ✨")
}

func createFileFromTemplate(path, tmpl string, data TemplateData) error {
	// Check if file already exists
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.Printf("  ! Skipped %s (already exists)", path)
		return nil
	}

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	t, err := template.New("file").Parse(tmpl)
	if err != nil {
		return err
	}

	return t.Execute(file, data)
}

const entityTemplate = `package entity

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/revandpratama/lognest/config"
	"gorm.io/gorm"
)

// {{.CapitalizedModuleName}} represents the data structure for a {{.ModuleName}}.
type {{.CapitalizedModuleName}} struct {
	ID          uuid.UUID      ` + "`" + `gorm:"type:uuid;primary_key" json:"id"` + "`" + `
	CreatedAt   time.Time      ` + "`" + `gorm:"not null" json:"created_at"` + "`" + `
	UpdatedAt   time.Time      ` + "`" + `gorm:"not null" json:"updated_at"` + "`" + `
	DeletedAt   gorm.DeletedAt ` + "`" + `gorm:"index" json:"-"` + "`" + `
}

// TableName sets the table name for the {{.CapitalizedModuleName}}.
func ({{.CapitalizedModuleName}}) TableName() string {
	return fmt.Sprintf("%s.%s", config.ENV.LOGNEST_SCHEMA, "{{.ModuleName}}s")
}

func (p *{{.CapitalizedModuleName}}) BeforeCreate(tx *gorm.DB) error {
	if p.ID == uuid.Nil {
		uuidGenerated, err := uuid.NewV7()
		if err != nil {
			return err
		}
		p.ID = uuidGenerated
	}
	return nil
}
`

const repositoryTemplate = `package repository

import (
	"gorm.io/gorm"
)

// {{.CapitalizedModuleName}}Repository defines the interface for database operations for a {{.CapitalizedModuleName}}.
type {{.CapitalizedModuleName}}Repository interface {
	
}

type {{.ModuleName}}Repository struct {
	db *gorm.DB
}

// New{{.CapitalizedModuleName}}Repository creates a new instance of {{.CapitalizedModuleName}}Repository.
func New{{.CapitalizedModuleName}}Repository(db *gorm.DB) {{.CapitalizedModuleName}}Repository {
	return &{{.ModuleName}}Repository{db: db}
}

// NOTE: The following are example implementations. You will need to adjust them.
`

const usecaseTemplate = `package usecase

import (
	"{{.GoModulePath}}/internal/modules/{{.OriginalModuleName}}/repository"
)

// {{.CapitalizedModuleName}}Usecase defines the business logic interface for a {{.CapitalizedModuleName}}.
type {{.CapitalizedModuleName}}Usecase interface {

}

type {{.ModuleName}}Usecase struct {
	repo repository.{{.CapitalizedModuleName}}Repository
}

// New{{.CapitalizedModuleName}}Usecase creates a new instance of {{.CapitalizedModuleName}}Usecase.
func New{{.CapitalizedModuleName}}Usecase(repo repository.{{.CapitalizedModuleName}}Repository) {{.CapitalizedModuleName}}Usecase {
	return &{{.ModuleName}}Usecase{repo: repo}
}


`

const handlerTemplate = `package handler

import (
	"{{.GoModulePath}}/internal/modules/{{.OriginalModuleName}}/usecase"
)

// {{.CapitalizedModuleName}}Handler defines the HTTP handler interface for a {{.CapitalizedModuleName}}.
type {{.CapitalizedModuleName}}Handler interface {
	// TODO: Add other handler methods (FindByUserID, FindBySlug, etc.)
}

type {{.ModuleName}}Handler struct {
	usecase usecase.{{.CapitalizedModuleName}}Usecase
}

// New{{.CapitalizedModuleName}}Handler creates a new instance of {{.CapitalizedModuleName}}Handler.
func New{{.CapitalizedModuleName}}Handler(usecase usecase.{{.CapitalizedModuleName}}Usecase) {{.CapitalizedModuleName}}Handler {
	return &{{.ModuleName}}Handler{usecase: usecase}
}

`
