# Chapter 14.1: Git & GitHub

## 前言
**Git** 是一個強大的分散式版本控制系統，而 **GitHub** 是目前最流行的 Git 儲存庫託管服務。它們是現代軟體開發的基石，讓團隊能夠高效協作、追蹤變更並管理程式碼歷史。本節將回顧其核心概念與重要的協作流程。

## 核心概念
- **Repository (儲存庫):** 您的專案資料夾，包含所有檔案和它們的修訂歷史。
- **Commit (提交):** 將您對檔案的變更儲存到儲存庫的快照。
- **Branch (分支):** 一個獨立的開發線，讓您可以開發新功能而不影響主線 (通常是 `main` 或 `master` 分支)。
- **Merge (合併):** 將一個分支的變更整合到另一個分支。
- **Pull Request (PR):** 一個請求，通知團隊成員您已經完成了一個功能的開發，並希望將您的變更合併到主分支中。這是進行程式碼審查 (Code Review) 的主要機制。

## GitHub Flow
GitHub Flow 是一個輕量級、基於分支的工作流程，非常適合大多數專案：
1.  在 `main` 分支中的任何內容都是可部署的。
2.  要進行新工作，請從 `main` 建立一個描述性的分支 (例如 `new-feature` 或 `bug-fix`)。
3.  在本地提交到該分支，並定期將您的工作推送到伺服器上的同名分支。
4.  當您需要回饋或準備好合併時，開啟一個 Pull Request。
5.  在審查並批准後，您可以將其合併到 `main` 分支。
6.  一旦合併到 `main`，就應該立即部署。

---

# Chapter 14.2: GitHub Actions

## 前言
**GitHub Actions** 是一個強大的 CI/CD 平台，直接整合在您的 GitHub 儲存庫中。它允許您自動化軟體開發工作流程，例如在每次提交或發布時自動建置、測試和部署您的 Go 應用程式。

## 核心概念
- **Workflow (工作流程):** 一個可配置的自動化過程，由一個或多個 `job` 組成。工作流程定義在儲存庫的 `.github/workflows` 目錄下的 YAML 檔案中。
- **Event (事件):** 觸發工作流程運行的特定活動，例如 `push`、`pull_request` 或 `schedule`。
- **Job (作業):** 在同一個 `runner` 上執行的一組 `step`。預設情況下，多個 job 會並行運行。
- **Step (步驟):** 一個獨立的任務，可以執行命令或一個 `action`。
- **Action (動作):** 一個可重複使用的程式碼單元，可以被包含在 `step` 中。您可以建立自己的 action，或使用 GitHub Marketplace 中的 action。
- **Runner (執行器):** 一個安裝了 GitHub Actions runner 應用程式的伺服器，用於運行您的工作流程。

### Go 工作流程範例
一個典型的 Go 專案 CI 工作流程可能包含以下步驟：
1.  **Checkout:** 拉取您的程式碼。
2.  **Setup Go:** 設定特定版本的 Go 環境。
3.  **Lint:** 執行靜態程式碼分析。
4.  **Test:** 運行單元測試。
5.  **Build:** 編譯您的應用程式。

```yaml
# .github/workflows/ci.yml
name: Go CI

on: [push, pull_request]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v4

    - name: Set up Go
      uses: actions/setup-go@v5
      with:
        go-version: '1.22'

    - name: Lint
      run: go vet ./...

    - name: Test
      run: go test -v ./...

    - name: Build
      run: go build -v ./...
```

---

# Chapter 14.3: Build Automation

## 前言
建置自動化是 CI/CD 流程的核心部分。對於 Go 專案，這通常意味著編譯程式碼、處理依賴項，並將應用程式打包成可部署的產物，例如一個二進制執行檔或一個 Docker 映像檔。

## 使用 `go build`
`go build` 是最基本的建置工具。您可以使用不同的旗標來為不同的作業系統和架構進行交叉編譯。

```bash
# 為 Linux 編譯
GOOS=linux GOARCH=amd64 go build -o my-app-linux

# 為 Windows 編譯
GOOS=windows GOARCH=amd64 go build -o my-app.exe
```

## 使用 `Makefile`
對於更複雜的建置流程，`Makefile` 是一個非常實用的工具。它允許您定義一系列的任務，例如 `build`, `test`, `lint`, `clean` 等。

```makefile
.PHONY: build test lint clean

build:
	go build -o my-app .

test:
	go test ./...

lint:
	golangci-lint run

clean:
	rm -f my-app
```

---

# Chapter 14.4: Deployment Strategies

## 前言
部署策略是在將新版本的軟體發布到生產環境時所採用的方法。選擇正確的策略可以最小化停機時間、降低風險，並確保用戶體驗的流暢。

## 常見策略
- **Recreate (滾動更新):** 最簡單的策略。先關閉舊版本的應用程式，然後再啟動新版本。會導致短暫的停機時間。
- **Rolling Update (滾動更新):** 逐一用新版本的實例替換舊版本的實例。可以避免停機時間，但會在一段時間內同時運行新舊兩個版本。
- **Blue-Green Deployment (藍綠部署):** 同時部署兩個完全相同的環境：「藍色」（當前版本）和「綠色」（新版本）。當新版本準備就緒時，只需將流量從藍色切換到綠色即可。如果出現問題，可以快速切換回來。
- **Canary Deployment (金絲雀部署):** 將一小部分流量（例如 5%）引導到新版本，同時大部分流量仍由舊版本處理。如果新版本表現穩定，則逐步增加流量，直到所有流量都轉移到新版本。這種方法可以及早發現問題並限制其影響範圍。

---

# Chapter 14.5: Environment Management

## 前言
有效的環境管理對於確保軟體品質和可靠的部署至關重要。開發團隊通常會使用多個環境來隔離不同的開發階段。

## 標準環境
- **Development (開發環境):** 開發人員在本地機器上使用的環境，用於日常的開發和測試。
- **Testing/Staging (測試/預備環境):** 一個模擬生產環境的鏡像，用於在部署前進行全面的整合測試、效能測試和使用者驗收測試 (UAT)。
- **Production (生產環境):** 實際提供給終端使用者使用的環境。

## 管理環境配置
應用程式在不同環境中通常需要不同的配置（例如資料庫連接字串、API 金鑰）。管理這些配置的最佳實務是：
- **使用環境變數:** 將配置從程式碼中分離出來，透過環境變數注入。這是符合 [十二因子應用 (The Twelve-Factor App)](https://12factor.net/config) 原則的做法。
- **使用設定檔:** 為每個環境提供一個設定檔（例如 `config.dev.json`, `config.prod.json`），並在應用程式啟動時加載相應的檔案。
- **集中式配置管理:** 對於大型系統，可以使用像 HashiCorp Consul 或 AWS Parameter Store 這樣的工具來集中管理和分發配置。

## 結語
恭喜您完成了 DevOps 與 CI/CD 章節！您現在了解了版本控制的最佳實務，學會了如何使用 GitHub Actions 自動化您的工作流程，並掌握了不同的部署策略和環境管理方法。

下一步，我們將進入 **雲端部署** 的世界，學習如何將您的容器化 Go 應用程式部署到主流的雲端平台。 
