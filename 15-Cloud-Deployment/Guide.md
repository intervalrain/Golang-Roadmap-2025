# Chapter 15.1: Cloud Platforms

## 前言
雲端運算徹底改變了應用程式的建置和部署方式。三大公有雲供應商——**Amazon Web Services (AWS)**、**Google Cloud Platform (GCP)** 和 **Microsoft Azure**——提供了廣泛的服務，從虛擬機、資料庫到機器學習和物聯網。本節將對它們進行簡要介紹。

## 主流平台概覽
- **AWS:** 市場領導者，提供最全面、最成熟的服務。對於 Go 開發者，核心服務包括:
    - **EC2 (Elastic Compute Cloud):** 虛擬伺服器。
    - **S3 (Simple Storage Service):** 物件儲存。
    - **RDS (Relational Database Service):** 託管式關聯式資料庫。
    - **EKS (Elastic Kubernetes Service):** 託管式 Kubernetes 服務。
    - **Lambda:** 無伺服器運算。

- **Google Cloud Platform (GCP):** 以其在 Kubernetes、大數據和機器學習領域的強大實力而聞名。核心服務包括:
    - **Compute Engine:** 虛擬機。
    - **Cloud Storage:** 物件儲存。
    - **Cloud SQL:** 託管式 MySQL, PostgreSQL, 和 SQL Server。
    - **GKE (Google Kubernetes Engine):** 業界領先的託管式 Kubernetes 服務。
    - **Cloud Functions:** 無伺服器運算。

- **Microsoft Azure:** 在企業市場中佔有重要地位，並與微軟生態系統深度整合。核心服務包括:
    - **Virtual Machines:** 虛擬機。
    - **Blob Storage:** 物件儲存。
    - **Azure SQL Database:** 託管式 SQL Server。
    - **AKS (Azure Kubernetes Service):** 託管式 Kubernetes 服務。
    - **Azure Functions:** 無伺服器運算。

---

# Chapter 15.2: Kubernetes

## 前言
**Kubernetes (K8s)** 是一個開源的容器編排平台，用於自動化容器化應用程式的部署、擴展和管理。它已成為在雲端運行應用程式的事實標準。

## 核心概念
- **Cluster (叢集):** 一組運行 Kubernetes 的節點 (Node)，分為 Master Node 和 Worker Node。
- **Node (節點):** 一台虛擬機或實體機，是 Kubernetes 的工作單元。
- **Pod:** Kubernetes 中最小的可部署單元。一個 Pod 可以包含一個或多個容器，它們共享儲存和網路資源。
- **Deployment:** 一個描述所需狀態的物件，例如運行多少個應用程式的副本。Kubernetes 會自動管理 Pod，以確保實際狀態與所需狀態相符。
- **Service:** 一個抽象層，定義了訪問一組 Pod 的方式。Service 提供了一個穩定的端點 (IP 位址和 DNS 名稱)，即使後端的 Pod 發生變化。
- **Ingress:** 管理對叢集中 Service 的外部訪問，通常是 HTTP。Ingress 可以提供負載平衡、SSL 終止和基於名稱的虛擬主機。

---

# Chapter 15.3: Load Balancing

## 前言
**Load Balancing (負載平衡)** 是將傳入的網路流量分配到多個後端伺服器（或應用程式實例）的過程。這是建構高可用性、高可靠性系統的關鍵。

## 為何需要負載平衡？
- **提高可用性:** 如果一個伺服器發生故障，負載平衡器會自動將流量重新導向到健康的伺服器。
- **增加擴展性:** 您可以輕鬆地在後端新增更多伺服器來處理增加的流量。
- **提升效能:** 將流量分配到多個伺服器，可以減少單一伺服器的負載，從而降低延遲。

在雲端環境中，負載平衡通常由託管服務提供，例如 AWS 的 **Elastic Load Balancer (ELB)**、GCP 的 **Cloud Load Balancing** 和 Azure 的 **Load Balancer**。在 Kubernetes 中，`Service` 和 `Ingress` 物件也扮演了負載平衡的角色。

---

# Chapter 15.4: Database Deployment

## 前言
在雲端部署資料庫時，您有兩種主要選擇：使用託管服務或自行管理。

## 託管式資料庫服務 (Managed Database Services)
例如 AWS RDS, GCP Cloud SQL, Azure SQL Database。這是大多數情況下的建議選擇。
- **優點:**
    - **簡化管理:** 雲端供應商會處理硬體佈建、軟體修補、備份和高可用性等繁瑣任務。
    - **高可用性和可靠性:** 內建的故障轉移和備份功能。
    - **易於擴展:** 只需點擊幾下即可擴展資料庫的計算和儲存資源。
- **缺點:**
    - **成本較高:** 比自行管理更昂貴。
    - **靈活性較低:** 對底層配置的控制較少。

## 在虛擬機上自行管理 (Self-managed on VMs)
在 EC2 或 Compute Engine 實例上自行安裝和管理資料庫（例如 PostgreSQL 或 MongoDB）。
- **優點:**
    - **完全控制:** 您可以完全控制資料庫的配置和版本。
    - **成本較低:** 可能比託管服務更具成本效益。
- **缺點:**
    - **管理複雜:** 您需要自行負責所有管理任務，包括安全性、備份和擴展。

---

# Chapter 15.5: Serverless

## 前言
**Serverless (無伺服器運算)** 是一種雲端運算模型，雲端供應商會動態管理應用程式所需的基礎設施。您只需編寫和部署程式碼，而無需擔心伺服器的佈建或管理。Go 的快速啟動時間和高效能使其成為 Serverless 的絕佳選擇。

## 核心概念
- **Function as a Service (FaaS):** Serverless 的核心。您將程式碼組織成獨立的函式，這些函式由事件觸發（例如 HTTP 請求、檔案上傳到 S3、新訊息到達佇列）。
- **按需付費:** 您只需為函式實際運行的時間付費，當函式沒有運行時，不收取任何費用。
- **自動擴展:** 平台會根據需求自動擴展或縮減您的函式實例數量。

## 主流平台
- **AWS Lambda:** 最流行和成熟的 FaaS 平台。與 AWS 生態系統緊密整合。
- **Google Cloud Functions:** GCP 的 FaaS 產品，同樣提供事件驅動的運算能力。
- **Azure Functions:** Microsoft Azure 的對應服務。

### Go on AWS Lambda 範例
```go
package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/lambda"
)

type MyEvent struct {
	Name string `json:"name"`
}

func HandleRequest(event MyEvent) (string, error) {
	return fmt.Sprintf("Hello %s!", event.Name), nil
}

func main() {
	lambda.Start(HandleRequest)
}
```

## 結語
恭喜您完成了雲端部署章節！您現在對如何在 AWS、GCP 和 Azure 等主要雲端平台上部署和管理 Go 應用程式有了全面的了解，從傳統的虛擬機到現代的 Kubernetes 和 Serverless 架構。

下一步，我們將探討 **微服務與架構**，學習如何設計和建構大型、可擴展且具韌性的分散式系統。
