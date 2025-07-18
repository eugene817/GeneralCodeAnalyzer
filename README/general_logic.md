## Общая логика работы приложения

Список использованых технологий
- Логика приложения написана на языке `golang` версии `1.23`
- В качестве REST API используется фреймворк `Echo`
- В качестве фронтенда используется связка:
    - `Tailwindcss` для написания стилей
    - `htmx` для запросов на бэкенд
    - `htmx-json-enc` расширение для отправки данных на бэкенд в формате json
    - `html` для написания структуры фронтенда
    - `JavaScript` для написания вспомогательных функций фронтенда (замена формуляров на `Monaco`)
    - `Monaco ediot` для более удобного вписания кода на фронтенде
- `Docker` для контейнеризации самого приложения
- `Docker-compose` для подключения внешних сервисов, таких как бд
- `DIND` - docker in docker, для того чтобы приложение могло запускать контейнеры
- `Adminer` - для автоматического создание admin страницы приложения
находясь внутри контейнера
- `Cowdocs` - библиотека для упрощения и ускроения работы с управлением и созданием контейнеров

---

Логика приложения предельно проста.

Со стороны пользователя:
- Пользователь производит регистрацию а затем логин
- По желанию пользователь выбирает какой код он хочет проанализировать
- Пользователь копирует и вставляет или пишет код, а затем выбирает какой тип анализа
он хочет выполнить
- Если пользователь выполняет статический анализ, то получает данные статического анализа на той же
странице где код, чтобы можно было бы исправить ошибки
- Если пользователь выполняет динамический анализ, то пользователя перенаправляют на страницу со метриками
динамического анализа, откуда он может вернуться либо на главную страницу либо на страницу анализа кода

Что происходит внутри:
- Приложение запускается в виде нескольких контейнеров: 
    - Контейнер с приложением
    - Контейнер с базой данных
    - Контейнер dind
- Перед запуском происходит компиляция приложения, а так же загрузка нужных образов контейнеров, чтобы не делать этого
при запросе, а уже иметь готовые для работы образы
- После запуска открывается доступ к приложению, а именно: localhost:8080/ - страница доступная без аутентификации
- Каждая страница или каждое изменение в странице которое видит пользователь это определенный endpoint, поскольку
приложения является backend-driven, тоесть это шаблоны с вставкой данных из бэкенда, но поскольку используется htmx,
это позволяет не менять всего html, а только частей его, что ускоряет работу приложения и позволяет делать более user-friendly
интерфейсы, при этом не углубляясь в сложности написания фронтенда через полноценные фреймворки
- При логине и регистрации происходят стандартные: шифрование пароля, генерация JWT токена, уставление cookie и 
перенаправление на страницы доступные после аутентификации
- После аутентификации, бэкенд проверяет JWT токен при каждом запросе
- При анализе кода, независимо статический или динамический, запрос отправляется на бэкенд. Где происходит загрузка этого кода
в уже подготовленный контейнер с подготовленной конфигурацией. Контейнер запускается, из него забираются метрики, так называемые логи
- Так же отправляется запрос на внешний сервис ИИ `ollama`, который так же анализирует код и отправляет данные обратно на бэкенд
(стандартный выход), все это обрабатывается, вставляется в шаблон и отправляется на фронтенд
- Посколку встроенное в `echo` распараллеливание запросов позволяет обрабатывать большое кол-во одновременно и конфигурация dind и библиотека
cowdocs настроены на параллельную работу, приложение может запускать большое кол-во контейнеров одновременно и в асинхронном формате обрабатывать их
состояние и работу. А так же тот факт что приложение написано на языке `go`, который компиллируется, и структура работы приложения выполнена на основе
максимально независимых микросервисов говорит о том что приложение способно на стабильную и эффективную работу с большим колличеством запросов

