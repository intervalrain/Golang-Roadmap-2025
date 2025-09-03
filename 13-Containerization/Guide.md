# Chapter 13.1: Docker Basics

## 前言
歡迎來到容器化的世界！在本章節中，我們將介紹 **Docker** 的基礎知識，這是一個強大的工具，可以幫助您將應用程式及其所有依賴項打包到一個標準化的單元中，稱為「容器」。學習 Docker 將使您的 Go 應用程式部署更加一致、可靠且高效。

## 什麼是容器？
您可以將容器想像成一個輕量級、可執行的軟體包，其中包含運行應用程式所需的一切：程式碼、執行時、系統工具、系統函式庫和設定。

與傳統的虛擬機 (VM) 不同，容器直接在主機的操作系統核心上運行，這使得它們非常輕巧且啟動速度極快。

## 核心概念
- **Image (映像檔):** 一個唯讀的模板，用於建立容器。映像檔包含了應用程式的程式碼以及所有依賴項。您可以從 Docker Hub 上拉取現有的映像檔，也可以自己建立。
- **Container (容器):** 映像檔的運行實例。您可以建立、啟動、停止、移動或刪除容器。
- **Dockerfile:** 一個文字檔案，其中包含了一系列指令，用於自動化地建立 Docker 映像檔。
- **Docker Hub:** 一個由 Docker, Inc. 提供的雲端服務，用於儲存和分享 Docker 映像檔。

---

# Chapter 13.2: Dockerfile 最佳實務

## 前言
`Dockerfile` 是建立映像檔的藍圖。一個編寫良好的 `Dockerfile` 不僅可以讓映像檔更小，還能加快建置速度並提高安全性。本節將介紹撰寫 `Dockerfile` 的最佳實務。

## 核心原則
1.  **保持映像檔輕量:**
    - 使用官方的輕量級基礎映像檔，例如 `golang:1.22-alpine`。Alpine Linux 是一個極簡的發行版，非常適合用於容器。
2.  **利用快取機制:**
    - 將不常變動的指令放在 `Dockerfile` 的前面，例如安裝依賴項。將經常變動的指令（例如複製原始碼）放在後面。這樣可以充分利用 Docker 的建置快取。
3.  **使用 `.dockerignore`:**
    - 建立一個 `.dockerignore` 檔案，以排除不需要複製到映像檔中的檔案和目錄，例如 `.git`、`*.md` 或本地開發環境的設定檔。

---

# Chapter 13.3: Multi-stage Builds

## 前言
Multi-stage builds (多階段建置) 是一個非常強大的功能，可以讓您在最終的映像檔中只保留必要的產物，從而大幅縮小映像檔的體積。

## 如何運作？
您可以在一個 `Dockerfile` 中使用多個 `FROM` 指令。每個 `FROM` 指令都代表一個新的建置階段。您可以從一個階段將檔案複製到另一個階段，從而將最終的應用程式與其建置環境分開。

### 範例
```dockerfile
# --- Build Stage ---
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /main .

# --- Final Stage ---
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /main .
CMD ["./main"]
```
在這個範例中：
- **builder stage:** 使用 `golang` 映像檔來編譯 Go 應用程式。
- **final stage:** 使用輕量的 `alpine` 映像檔，並只從 `builder` 階段複製編譯好的執行檔。最終的映像檔將不包含任何 Go 工具鏈或原始碼。

---

# Chapter 13.4: Docker Compose

## 前言
當您的應用程式由多個服務組成時（例如一個 Web 伺服器和一個資料庫），手動管理它們會變得很麻煩。**Docker Compose** 是一個工具，可讓您使用一個 `YAML` 檔案來定義和運行多容器的 Docker 應用程式。

## 核心概念
- **`docker-compose.yml`:** 一個 YAML 檔案，用於配置應用程式的服務、網路和儲存卷。
- **Services:** 組成應用程式的各個容器，例如 `web`、`db`、`redis`。
- **Networks:** 讓服務之間可以互相通訊。
- **Volumes:** 用於持久化容器中的資料。

---

# Chapter 13.5: Container Security

## 前言
雖然容器提供了隔離性，但確保其安全性仍然至關重要。本節將介紹一些容器安全性的基本原則。

## 最佳實務
1.  **使用非 Root 使用者:**
    - 在 `Dockerfile` 中建立一個非 root 使用者，並使用 `USER` 指令切換到該使用者來運行應用程式。這可以減少潛在的攻擊面。
2.  **最小權限原則:**
    - 不要授予容器不必要的權限。例如，避免使用 `--privileged` 旗標。
3.  **掃描映像檔:**
    - 定期使用工具（如 `Trivy` 或 `Docker Scout`）掃描您的映像檔，以查找已知的漏洞。
4.  **管理 Secret:**
    - 不要將敏感資訊（如 API 金鑰或密碼）硬編碼到 `Dockerfile` 或映像檔中。使用 Docker secrets 或其他 secret 管理工具。

## 結語
恭喜您完成了容器化章節的學習！您現在已經掌握了使用 Docker 將 Go 應用程式容器化的基礎知識，從建立映像檔到管理多容器應用程式，再到確保容器的安全性。

下一步，我們將探索 **DevOps 與 CI/CD**，學習如何自動化建置、測試和部署您的容器化應用程式。
