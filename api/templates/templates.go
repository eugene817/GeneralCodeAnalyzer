package templates

import (
  "html/template"
  "io"
  "github.com/labstack/echo/v4"
  "log"
  "os"
  "path/filepath"
  "strings"
)


type Template struct {
  Templates *template.Template
}

func NewTemplate() *Template {
  tmpl := template.New("")
  err := filepath.Walk("api/templates", func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }
    if !info.IsDir() && strings.HasSuffix(info.Name(), ".html") {
      _, err := tmpl.ParseFiles(path)
      if err != nil {
        return err
      }
    }
    return nil
  })
  if err != nil {
    log.Fatal(err)
  }
  return &Template{
    Templates: tmpl,
  }
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
  return t.Templates.ExecuteTemplate(w, name, data)
}



