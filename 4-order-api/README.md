# purpleschool

## Для Windows (PowerShell)
1. Откройте PowerShell в корне проекта
2. Выполните:
```powershell
# Запуск приложения
.\scripts\windows\run.ps1
# Выполнение миграций
.\scripts\windows\migrations.ps1

## Для masOS 
1. Задаём права для выполнения скриптов
chmod +x scripts/macos/*.sh
# Запуск приложения
./scripts/macos/run.sh
# Выполнение миграций
./scripts/macos/migrations.sh

## Для запуска e2e теста, необходимо создать тестовую базу данных и на ней выполнить миграции, указав в файле конфигурации .env данные для тестовой базы. 