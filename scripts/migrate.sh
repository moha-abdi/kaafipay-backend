#!/bin/bash

# Load environment variables from .env file
if [ -f .env ]; then
    export $(cat .env | grep -v '#' | xargs)
fi

# Check if DATABASE_URL is set
if [ -z "$DATABASE_URL" ]; then
    echo "Error: DATABASE_URL environment variable is not set"
    exit 1
fi

# Default command is "up" if no argument is provided
COMMAND=${1:-up}

# Function to show usage
show_usage() {
    echo "Usage: ./scripts/migrate.sh <command> [args]"
    echo ""
    echo "Commands:"
    echo "  up                     Run all pending migrations"
    echo "  down                   Roll back all migrations"
    echo "  create <name>          Create a new migration"
    echo "  version                Show current migration version"
    echo "  force <version>        Force set database version"
    echo ""
    echo "Examples:"
    echo "  ./scripts/migrate.sh up"
    echo "  ./scripts/migrate.sh create add_users_table"
    echo "  ./scripts/migrate.sh force 12"
}

case $COMMAND in
    "up")
        echo "Running all pending migrations..."
        if ! migrate -database "${DATABASE_URL}" -path internal/db/migrations up; then
            echo "Error: Migration failed. If database is in dirty state, use 'force' command to fix it."
            echo "Example: ./scripts/migrate.sh force <version>"
            exit 1
        fi
        ;;
    "down")
        echo "Rolling back all migrations..."
        if ! migrate -database "${DATABASE_URL}" -path internal/db/migrations down; then
            echo "Error: Migration rollback failed"
            exit 1
        fi
        ;;
    "create")
        if [ -z "$2" ]; then
            echo "Error: Migration name is required for create command"
            echo "Usage: ./scripts/migrate.sh create <migration_name>"
            exit 1
        fi
        echo "Creating new migration files..."
        migrate create -ext sql -dir internal/db/migrations -seq "$2"
        ;;
    "version")
        echo "Checking current migration version..."
        migrate -database "${DATABASE_URL}" -path internal/db/migrations version
        ;;
    "force")
        if [ -z "$2" ]; then
            echo "Error: Version number is required for force command"
            echo "Usage: ./scripts/migrate.sh force <version>"
            exit 1
        fi
        echo "Forcing database version to $2..."
        migrate -database "${DATABASE_URL}" -path internal/db/migrations force "$2"
        ;;
    "help"|"-h"|"--help")
        show_usage
        ;;
    *)
        echo "Invalid command: $COMMAND"
        show_usage
        exit 1
        ;;
esac 