# General Code Analyzer

- This is my diploma project

The core functionality is to analyze SQL code with statistical and dynamic analyzis,
The whole diploma project is gonna be here (in polish and english languages)

## The idea

The idea is to make general code analyzer, maybe after that train ai models to recognise
"ideal code" and to make life of begginer programmers eazier


## How to run
(For now)

If you already have golang installed:
- Simply write in the directory Backend
`go run .`

## How to use
(For now)

Send POST request for the localhost:8090/analyze
with sql query and sql init db
```Go
// incoming request
type AnalyzeRequest struct {
  SQLQuery string `json:"sql_query"`
  InitSQL string `json:"init_sql"`
}
```


