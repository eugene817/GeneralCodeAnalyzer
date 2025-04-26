# General Code Analyzer

- This is my diploma project

The core functionality is to analyze SQL code with statistical and dynamic analyzis,
The whole diploma project is gonna be here (in polish and english languages)

## The idea

The idea is to make general code analyzer, maybe after that train ai models to recognise
"ideal code" and to make life of begginer programmers eazier


## How to run

Simply run:

```bash
docker-compose up --build
```


#Docs

@ru

Начнем распутывать это запутанное и очень абстрактное приложение с `main.go`

### `main.go` 

Это файл который практически является "высокоуровневым сердцем" приложения

Что здесь происходит:

0. Описание импортов внешних библиотек и так называемого `package`, который нужен для работы с теми же импортами

#### func main()

1. Инициализация мэнеджера из библиотеки `Cowdocs` (перейдем к этому позже) и создание инстанции API
```go
  mgr, err := container.NewDockerManager()
	if err != nil {
		fmt.Errorf("failed to create Docker manager: %v", err)
    os.Exit(1)
	}
  mng := api.NewAPI(mgr)
```

2. Проверка наличия нужных контейнеров для работы приложения (внутри если их нет, они устанавливаются)
```go
  // Ensure the images are available
  Images := []string{
    "python:3",
    "keinos/sqlite3",
  }
  
  for {
        if err := mng.Ping(); err == nil {
            break
        }
        log.Println("waiting for Docker daemon…")
        time.Sleep(1500 * time.Millisecond)
  }

  if err := mng.EnsureImages(Images); err != nil {
    log.Fatalf("failed to pull initial images: %v", err)    
    os.Exit(1)
  }
```
В цикле for в свою очередь происходит очень интересная вещь, поскольку приложение запускается
через `docker-compose`, то `docker-compose` должен как-то обслужить создание и запуск новых контейнеров
внутри себя. Для этого существует модуль `dind` - docker in docker, ожидание запуска которого и происходит
в цикле for


3. Создание инстанции главного интерфейса web-фреймворка `Echo` и инициализация базы данных
```go
  e := echo.New()
  db, err := database.InitDB()
```

4. Создание инстанции интерфейса сервиса `svc` и инстанции обработчика `h`, а после запускается
функция регистрации маршрутов через обработчик
```go
  svc := services.NewService(mng)
  h := handlers.NewHandler(svc, db)
  h.RegisterRoutes(e)
```

5. Поскольку приложение не имеет "полноценного" фронтенда, а работает на `html`, `tailwindcss`, `htmx`
то требуется обработка статических файлов (в данном случае стилей) и шаблонов `html`. Что и происходит здесь:
```go
	templates.RegisterTemplatesRoutes(e)
	e.Static("/static", "./api/static")
```

6. Ну и поскольку приложение запускается через фреймворк `Echo`, то здесь происходит его запуск:
```go
	port := config.GetPort()
	e.Logger.Fatal(e.Start(port))
```

Дальше я бы перешел к сервисам:

### services/

Здесь находится несколько важных файлов:

#### Ядро
- `services.go`
#### Экзекуторы (исполнители)
- `llm_executor.go`
- `python_executor.go`
- `sql_executor.go`
#### Генерация рекомендаций
- `recommendations.go`
#### Инструменты
- `utils.go`


#### Ядро
- `services.go`
В ядре не происходит многого, описывается импорт, интерфейс сервиса и функиця для создания его инстанции

#### Экзекуторы
- `llm_executor.go`
