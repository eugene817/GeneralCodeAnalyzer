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

## Examples of work
- 1

Request
```bash
curl -X POST http://localhost:8090/analyze \
-H "Content-Type: application/json" \
-d '{"sql_query": "SELECT department, COUNT(*) AS employee_count FROM employees WHERE salary > 50000 GROUP BY department HAVING COUNT(*) > 1 ORDER BY employee_count DESC;", "init_sql": "CREATE TABLE employees (id INTEGER, name TEXT, department TEXT, salary REAL); INSERT INTO employees (id, name, department, salary) VALUES (1, '\''Alice'\'', '\''HR'\'', 60000), (2, '\''Bob'\'', '\''Engineering'\'', 80000), (3, '\''Charlie'\'', '\''Engineering'\'', 75000), (4, '\''Diana'\'', '\''HR'\'', 55000), (5, '\''Eve'\'', '\''Marketing'\'', 40000), (6, '\''Frank'\'', '\''Engineering'\'', 70000);"}' | jq
```

Answer
```bash
{
  "analysis": "QUERY PLAN\n|--SCAN employees\n|--USE TEMP B-TREE FOR GROUP BY\n`--USE TEMP B-TREE FOR ORDER BY",
  "metrics": {
    "execution_time": 0.240345139,
    "memory_difference": 0,
    "memory_used_after": 0,
    "memory_used_before": 0
  },
  "recommendations": [
    "No indexes are used in this query. Consider adding indexes."
  ],
  "result": "Engineering|3\nHR|2"
}
```


