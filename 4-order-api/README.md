### E-Commerce Backend API (Go)
Backend service in Go for an e-commerce platform with code-based authorization. The service generates verification codes and validates them, emulating the SMS login process. It includes product management and protected endpoints for order placement.

## For Windows (PowerShell)
1. Open PowerShell in the project root directory
# Start the application
.\scripts\windows\run.ps1
# Run migrations
.\scripts\windows\migrations.ps1

## For macOS 
1. Make scripts executable
chmod +x scripts/macos/*.sh
# Start the application
./scripts/macos/run.sh
# Run migrations
./scripts/macos/migrations.sh

# Before running e2e tests:
1. Create a test database
2. Run migrations against it
3. Configure test database connection in the .env file
4. Create a .env file in the e2e directory