# 📊 ClickHouse & Flat File Data Integration Tool

A **web-based application** for **bidirectional data ingestion** between **ClickHouse** databases and **flat files (CSV)**.  
Built with **Go (backend)** and **HTML/JS (frontend)**.

---

## 🚀 Features

- ✅ **Bidirectional Data Flow**
  - ClickHouse → CSV
  - CSV → ClickHouse
- ✅ **JWT Authentication** for ClickHouse
- ✅ **Schema Discovery & Column Selection**
- ✅ **Data Preview** (First 100 rows)
- ✅ **Error Handling & Progress Reporting**

---

## ⚙️ Setup & Installation

### 1. Prerequisites

- [Docker](https://www.docker.com/)
- [Go 1.21+](https://golang.org/)
- [Node.js (optional)](https://nodejs.org/) – For frontend development

### 2. Run ClickHouse via Docker

```bash
docker run -d \
  --name clickhouse-server \
  -p 8123:8123 -p 9000:9000 \
  -e CLICKHOUSE_USER=default \
  -e CLICKHOUSE_PASSWORD=password \
  clickhouse/clickhouse-server:latest



Verify ClickHouse is running:

bash
Copy
Edit
curl http://localhost:8123/ping  # Should return "Ok"
3. Configure Environment Variables
Create a .env file in the project root:

ini
Copy
Edit
CLICKHOUSE_HOST=localhost
CLICKHOUSE_PORT=8123
CLICKHOUSE_DATABASE=default
CLICKHOUSE_USER=default
CLICKHOUSE_JWT_TOKEN=your_jwt_token_here  # Optional for local testing
CLICKHOUSE_SECURE=false
4. Run the Application
bash
Copy
Edit
# Install dependencies
go mod tidy

# Start the backend server
go run cmd/server/main.go
Frontend will be available at:
🌐 http://localhost:8080

📌 API Endpoints

Endpoint	Method	Description
/api/connect	POST	Test ClickHouse/CSV connection
/api/tables	GET	List available tables
/api/columns	GET	Fetch columns for a table/file
/api/preview	POST	Preview first 100 rows
/api/ingest	POST	Start data ingestion
🔧 Usage
1. ClickHouse → CSV
Select ClickHouse as source.

Enter connection details.

Choose a table & columns.

Set output file path (e.g., output.csv).

Click "Start Ingestion".

2. CSV → ClickHouse
Select Flat File as source.

Upload a CSV file.

Map columns to the target table.

Click "Start Ingestion".

🧪 Testing
Sample CSV (test_data.csv)
csv
Copy
Edit
order_id,order_date,customer_id,product_name,quantity,unit_price,total_price,country
1001,2023-01-15,3021,Laptop,1,999.99,999.99,US
1002,2023-01-16,4172,Smartphone,2,699.50,1399.00,UK
Check Ingested Data
bash
Copy
Edit
curl "http://localhost:8123?query=SELECT * FROM your_table"
📂 Project Structure
bash
Copy
Edit
/ch2csv  
├── cmd/server/main.go         # Entry point  
├── internal/  
│   ├── handlers/api.go        # API endpoints  
│   ├── services/clickhouse.go # ClickHouse logic  
│   ├── services/csv.go        # CSV parsing  
│   └── models/models.go       # Data models  
├── static/                    # Frontend (HTML/JS/CSS)  
├── go.mod                     # Dependencies  
└── README.md                  # Project readme  