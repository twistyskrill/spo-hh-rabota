
  # Handyman Service Website Layout

  This is a code bundle for Handyman Service Website Layout. The original project is available at https://www.figma.com/design/Y0B7UeWoHDVtDqKHm6dyuO/Handyman-Service-Website-Layout.

  ## Running the code

  Run `npm i` to install the dependencies.

  Run `npm run dev` to start the development server.

  ## Как запустить AI-функционал

Для работы функций «Придумать описание» и «Узнать рыночную цену» используется локальная нейросеть `llama3`, запущенная через `Ollama`. Это позволяет генерировать тексты абсолютно бесплатно и оффлайн. 

### Шаг 1. Установка Ollama
Скачайте и установите легковесный сервер Ollama:
- **Windows**: Запустить в терминале PowerShell команду `irm https://ollama.com/install.ps1 | iex` или скачать инсталлятор с [ollama.com/download](https://ollama.com/download)
- **Linux/macOS**: Запустить команду `curl -fsSL https://ollama.com/install.sh | sh`

### Шаг 2. Загрузка модели llama3
После установки Ollama откройте терминал и выполните команду:
```bash
ollama pull llama3
```
*Загрузка займёт некоторое время (размер модели ~4.7 ГБ).*

### Шаг 3. Запуск сервера Ollama
Убедитесь, что сервер запущен. Обычно он запускается автоматически как фоновый процесс или служба OS.
Если вы остановили его вручную, запустите команду:
```bash
ollama serve
```
REST API будет доступен по локальному адресу `http://localhost:11434`.
  