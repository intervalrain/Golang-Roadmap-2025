# 第 1 章：開發環境設定

歡迎來到 Go 語言學習之路的第一章！在本章節中，我們將引導您完成 Go 開發環境的設定。一個好的開發環境是高效學習和開發的基礎。

## 學習重點

- 安裝 Go 語言
- 設定您的工作區 (Workspace)
- 選擇並設定一個程式碼編輯器
- 撰寫並執行您的第一個 Go 程式

---

## 1. 安裝 Go

首先，我們需要從官方網站下載並安裝 Go。

1.  **前往官方下載頁面**：[https://go.dev/dl/](https://go.dev/dl/)
2.  **選擇您的作業系統版本**：根據您的電腦（Windows, macOS, or Linux）下載對應的安裝檔。
3.  **執行安裝**：
    *   **macOS**: 開啟下載的 `.pkg` 檔案，並依照提示完成安裝。預設會安裝在 `/usr/local/go`。
    *   **Windows**: 執行 `.msi` 安裝檔，依照提示完成。預設會安裝在 `C:\Program Files\Go`。
    *   **Linux**: 將下載的 `.tar.gz` 檔案解壓縮到 `/usr/local`。

4.  **驗證安裝**：開啟一個新的終端機（Terminal）或命令提示字元（Command Prompt），輸入以下指令：

    ```bash
    go version
    ```

    如果安裝成功，您應該會看到類似 `go version go1.xx.x darwin/amd64` 的版本資訊。

---

## 2. 設定工作區

從 Go 1.11 版本開始，Go 引入了 `Go Modules` 的概念，這使得專案可以存在於 `GOPATH` 之外的任何地方，`GOPATH` 的重要性因此降低了。我們強烈建議所有新專案都使用 Go Modules。

您的預設 `GOPATH` 通常在 `$HOME/go`（macOS/Linux）或 `%USERPROFILE%\go`（Windows）。您可以透過 `go env` 指令查看相關的環境變數。

為了方便管理，建議您可以建立一個專門用來存放 Go 專案的資料夾，例如 `~/go-projects`。

---

## 3. 選擇並設定編輯器

一個好的程式碼編輯器可以大幅提升您的開發效率。

- **Visual Studio Code (VS Code)**：這是目前最受歡迎的選擇，免費且強大。
    1.  [下載並安裝 VS Code](https://code.visualstudio.com/)。
    2.  在 VS Code 中，前往擴充套件市場，搜尋並安裝官方的 **Go** 擴充套件 (`golang.Go`)。
    3.  安裝後，VS Code 可能會提示您安裝額外的 Go 工具（如 `gopls`, `gofmt` 等），請務必點選 "Install All"。

- **GoLand**：由 JetBrains 公司開發的專業 Go IDE，功能非常強大，但需要付費。

---

## 4. 您的第一個 Go 程式

現在，讓我們來撰寫並執行一個簡單的 "Hello, World!" 程式來確認環境是否設定成功。

1.  **建立專案資料夾**：
    在您喜歡的位置（例如我們剛才提到的 `~/go-projects`）建立一個新的資料夾，例如 `hello-world`。

2.  **初始化 Go Module**：
    在終端機中，進入 `hello-world` 資料夾，並執行以下指令：

    ```bash
    go mod init example.com/hello-world
    ```
    `example.com/hello-world` 是這個模組的路徑，您可以換成您自己的。執行成功後，資料夾內會多一個 `go.mod` 檔案。

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
    在終端機中，執行以下指令：

    ```bash
    go run main.go
    ```

    如果一切順利，您將會在螢幕上看到 `Hello, World!` 的輸出。

---

## 結論

恭喜！您已經成功設定了 Go 的開發環境，並執行了您的第一個 Go 程式。這是一個非常重要的里程碑。

接下來，請您親自動手操作一遍，確保每個步驟都能順利完成。在下一個階段，我們將會對本章節做一個總結。
