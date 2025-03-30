# Internal Transfers Application

This application implements an internal transfers system with HTTP endpoints using Golang and PostgreSQL.

## Features

- **Account Creation Endpoint** (`POST /accounts`): Create a new account with an initial balance.
- **Account Query Endpoint** (`GET /accounts/{account_id}`): Retrieve the current balance of a given account.
- **Transaction Submission Endpoint** (`POST /transactions`): Process a transfer between two accounts.

## Requirements

- Go (1.16 or later)
- PostgreSQL

## Setup & Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/shbd845/TransactionSystem.git
   cd TransactionSystem

2. **Create a .env file**

   ```
   Copy .envexample file and rename it to .env
   Then provide details of your postgres server

3. **Run the application**

    ```
    go run cmd/main.go 