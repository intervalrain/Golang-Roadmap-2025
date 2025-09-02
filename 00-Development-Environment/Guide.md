# Chapter 0: Development Environment Setup

歡迎來到 Go 語言學習之路的第一章！在本章節中，我們將引導您完成 Go 開發環境的設定。一個好的開發環境是高效學習和開發的基礎。

## 學習重點

- 安裝 Go
- 設定您的 Workspace
- 選擇並設定一個 Code Editor
- 撰寫並執行您的第一個 Go 程式

---

## 1. Install Go

首先，我們需要從官方網站下載並安裝 Go。

1.  **前往官方下載頁面**：[https://go.dev/dl/](https://go.dev/dl/)
2.  **選擇您的作業系統版本**：根據您的電腦（Windows, macOS, or Linux）下載對應的安裝檔。
3.  **執行安裝**：
    *   **macOS**: 開啟下載的 `.pkg` 檔案，並依照提示完成安裝。預設會安裝在 `/usr/local/go`。
    *   **Windows**: 執行 `.msi` 安裝檔，依照提示完成。預設會安裝在 `C:\Program Files\Go`。
    *   **Linux**: 將下載的 `.tar.gz` 檔案解壓縮到 `/usr/local`。

4.  **驗證安裝**：開啟一個新的 **Terminal** 或 **Command Prompt**，輸入以下指令：

    ```bash
    go version
    ```

    如果安裝成功，您應該會看到類似 `go version go1.xx.x darwin/amd64` 的版本資訊。

---

## 2. Setup Workspace

從 Go 1.11 版本開始，Go 引入了 **Go Modules** 的概念，這使得專案可以存在於 `GOPATH` 之外的任何地方，`GOPATH` 的重要性因此降低了。我們強烈建議所有新專案都使用 **Go Modules**。

您的預設 `GOPATH` 通常在 `$HOME/go`（macOS/Linux）或 `%USERPROFILE%\go`（Windows）。您可以透過 `go env` 指令查看相關的環境變數。

為了方便管理，建議您可以建立一個專門用來存放 Go 專案的資料夾，例如 `~/go-projects`。

---

## 3. Editor and Extensions

一個好的 **Code Editor** 可以大幅提升您的開發效率。我們推薦使用 Visual Studio Code。

### 3.1. Install Visual Studio Code

- **下載與安裝**：前往 [Visual Studio Code 官方網站](https://code.visualstudio.com/) 下載並安裝適合您作業系統的版本。

### 3.2. Install Go Extensions

為了在 VS Code 中獲得最好的 Go 開發體驗（例如：程式碼自動完成、定義跳轉、語法檢查等），您需要安裝官方的 Go **Extension**。

1.  **開啟 VS Code**。
2.  點擊左側活動欄中的 **Extensions** 圖示。
3.  在搜尋框中輸入 `Go`。
4.  找到由 **Go Team at Google** 開發的 **Extension**（通常是第一個），點擊 **Install**。
5.  安裝完成後，VS Code 右下角可能會彈出提示，詢問是否安裝額外的 Go 分析工具（如 `gopls`, `gofmt` 等）。請務必點選 **Install All**，這些工具對於開發至關重要。

### 3.3. Other Editor Options

- **GoLand**：由 JetBrains 公司開發的專業 Go IDE，功能非常強大，但需要付費。

---

## 4. Your First Go Program

現在，讓我們來撰寫並執行一個簡單的 "Hello, World!" 程式來確認環境是否設定成功。

1.  **建立專案資料夾**：
    在您喜歡的位置（例如我們剛才提到的 `~/go-projects`）建立一個新的資料夾，例如 `hello-world`。

2.  **初始化 Go Module**：
    在 **Terminal** 中，進入 `hello-world` 資料夾，並執行以下指令：

    ```bash
    go mod init example.com/hello-world
    ```
    `example.com/hello-world` 是這個 **module** 的路徑，您可以換成您自己的。執行成功後，資料夾內會多一個 `go.mod` 檔案。

3.  **建立 `main.go` 檔案**：
    在資料夾中建立一個名為 `main.go` 的檔案，並貼上以下程式碼：

    ```go
    package main

    import "fmt"

    func main() {
        fmt.Println("Hello, World!")
    }
    ```

4.  **執行程式**：
    在 **Terminal** 中，執行以下指令：

    ```bash
    go run main.go
    ```

    如果一切順利，您將會在螢幕上看到 `Hello, World!` 的輸出。

---

## Conclusion

恭喜！您已經成功設定了 Go 的開發環境，並執行了您的第一個 Go 程式。這是一個非常重要的里程碑。

接下來，請您親自動手操作一遍，確保每個步驟都能順利完成。在下一個階段，我們將會對本章節做一個總結。