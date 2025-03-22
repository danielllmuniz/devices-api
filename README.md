# Device API

## 📌 About the Project
This API was developed in Golang and allows device management. The main operations include:

- Create a device
- Update a device
- Apply a patch to a device
- Get all devices
- Get devices by brand
- Get devices by state
- Get a device by ID
- Delete a device

The API was developed using:
- **Golang**
- **Chi** (lightweight router for HTTP APIs)
- **Tern** (migration manager for SQL databases)
- **SQLC** (safe and efficient SQL query generator for Go)

## 🚀 How to Run the Project
The API runs inside a Docker container and can be easily started using the `Makefile`.

### 1️⃣ Prerequisites
Make sure you have installed:
- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- [Make](https://www.gnu.org/software/make/)

### 2️⃣ Running the Application
1. Clone this repository:
   ```sh
   git clone https://github.com/danielllmuniz/devices-api.git
   cd device-api
   ```

2. Copy the `.env.example` file to `.env` and configure your environment variables:
   ```sh
   cp .env.example .env
   ```

3. Start the application containers:
   ```sh
   make up
   ```

4. Run the database migrations:
   ```sh
   make migrate
   ```

Now the API will be running at `http://localhost:8000` 🚀

## 📜 Available Endpoints
### Devices (`/devices`)
| Method  | Route                         | Description |
|---------|-------------------------------|-------------|
| `POST`  | `/devices`                    | Create a device |
| `PUT`   | `/devices/{id}`               | Update a device |
| `PATCH` | `/devices/{id}`               | Apply a patch to a device |
| `GET`   | `/devices`                    | Get all devices |
| `GET`   | `/devices?brand=brandName`    | Get devices by brand |
| `GET`   | `/devices?state=stateName`    | Get devices by state |
| `GET`   | `/devices/{id}`               | Get a device by ID |
| `DELETE`| `/devices/{id}`               | Delete a device |

## 🛠 Technologies Used
- **Golang** - Main programming language of the project
- **Chi** - Lightweight HTTP router for APIs
- **Tern** - Database migration management
- **SQLC** - SQL query code generation for Go
- **Docker** - Application containerization

## 📄 License
This project is licensed under the MIT License.