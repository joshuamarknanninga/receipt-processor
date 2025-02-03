# Receipt Processor API

A RESTful API service built with Go, Docker, and MongoDB to process receipts and calculate reward points based on specific rules.

## Features
- Process receipts via POST endpoint
- Retrieve points for receipts via GET endpoint
- Data persistence using MongoDB (in-memory during service lifetime)
- Dockerized setup for easy deployment
- Input validation and error handling

## Technologies Used
- **Go** (Golang) - Core backend language
- **Gin** - HTTP web framework
- **MongoDB** - NoSQL database for receipt storage
- **Docker** - Containerization and service orchestration
- **UUID** - Unique identifier generation

## API Endpoints
| Method | Endpoint                | Description                     |
|--------|-------------------------|---------------------------------|
| POST   | `/receipts/process`     | Submit receipt for processing   |
| GET    | `/receipts/{id}/points` | Get points for processed receipt|

## Installation

### Prerequisites
- Docker and Docker Compose
- Go 1.21+ (optional for local development)

### Setup
1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/receipt-processor.git
   cd receipt-processor
   ```

2. Start services using Docker:
   ```bash
   docker-compose up
   ```

3. The API will be available at:
   ```http
   http://localhost:8080
   ```

## API Specification

### Process Receipt
**Request:**
```bash
curl -X POST http://localhost:8080/receipts/process \
  -H "Content-Type: application/json" \
  -d '{
    "retailer": "Target",
    "purchaseDate": "2022-01-01",
    "purchaseTime": "13:01",
    "items": [
      {"shortDescription": "Mountain Dew 12PK", "price": "6.49"}
    ],
    "total": "6.49"
  }'
```

**Response:**
```json
{"id": "7fb1377b-b223-49d9-a31a-5a02701dd310"}
```

### Get Points
**Request:**
```bash
curl http://localhost:8080/receipts/7fb1377b-b223-49d9-a31a-5a02701dd310/points
```

**Response:**
```json
{"points": 28}
```

## Development Journey

This project represents a significant learning experience in building a production-ready API service. While I had fundamental knowledge of Go and Docker, implementing them together with MongoDB required substantial research and practice.

**Key Learning Points:**
1. **Go Web Development**: Through extensive Googling and AI-assisted research, I learned to:
   - Structure Go projects effectively
   - Implement clean REST API handlers with Gin
   - Manage dependencies using Go modules
   - Handle concurrent operations safely

2. **Docker Integration**: The Docker setup was created through:
   - Official Docker documentation study
   - Trial-and-error with multi-stage builds
   - Learning container networking patterns
   - Optimizing image sizes for production

3. **MongoDB Implementation**: Chosen for its document-oriented nature that aligns well with receipt data. While alternatives like:
   - **In-memory storage** (for simpler cases)
   - **Redis** (for pure cache needs)
   - **PostgreSQL** (for relational data)
   could be viable options, MongoDB provided excellent practice with:
   - BSON data marshaling
   - Official MongoDB Go driver usage
   - Document-oriented design patterns
   - Containerized database management

   <https://streamable.com/hublu5>

   <https://streamable.com/11l9uc>

**Challenges Overcome:**
- Implementing complex validation rules
- Creating efficient Docker-compose workflows
- Debugging container networking issues
- Ensuring proper database connection handling
- Mastering Go's context management

## Testing
1. **Using Insomnia/Postman**:
   - Import the provided test receipts from `test_receipts/` directory
   - Verify responses and error handling

2. **Manual Testing**:
   ```bash
   # Process receipt
   curl -X POST http://localhost:8080/receipts/process -H "Content-Type: application/json" -d @test_receipts/receipt-1.json

   # Get points
   curl http://localhost:8080/receipts/<ID>/points
   ```

3. **Database Inspection**:
   ```bash
   docker exec -it receipt-processor-mongodb-1 mongosh
   use receipts
   db.receipts.find().pretty()
   ```

## Acknowledgments
This project was made possible through:
- **Go Documentation** - Essential language reference
- **Docker Community Resources** - Containerization guidance
- **MongoDB University** - Database implementation patterns
- **AI Assistance** - Debugging help and code suggestions
- **Stack Overflow** - Community problem-solving insights

Special thanks to the open-source community for providing invaluable learning resources that helped bridge knowledge gaps throughout development.

## License
MIT License - See [LICENSE](LICENSE) for details.

---

_This project demonstrates contemporary backend development practices while highlighting the learning journey of implementing multiple technologies in tandem. The choice of MongoDB reflects both practical application of learned skills and exploration of NoSQL patterns in Go applications._
```