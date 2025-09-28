# farming-server

## Installation

1. Clone the repository:

    ```bash
    git clone https://github.com/hrutik1235/farming-server
    ```

2. Install **Air** (if not already installed):

   ```bash
   go install github.com/air-verse/air@latest
   ```

---

## Running the Application

> **Note:** The `PORT` environment variable is **mandatory**.
> The API will fail to start if `PORT` is not set.

### Development (with Air)

```bash
PORT=6000 air .
```

## 🧪 Test the API

```bash
curl http://localhost:6000/api/v1/health
```

## 📄 License

MIT License
