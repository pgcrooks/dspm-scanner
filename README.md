# dspm-scanner

## Architecture

```mermaid
graph TD;
    LocalOrch-->ScraperService;
    LocalOrch-->DataService;
    LocalOrch-->Scanner;

    ScraperService-->ScraperLocal;
    ScraperService-->ScraperAWS;

    ScraperLocal-->ScraperLocalWorker0;
    ScraperLocal-->ScraperLocalWorkerN;

    ScraperAWS-->ScraperAWSWorker0;
    ScraperAWS-->ScraperAWSWorkerN;

    Scanner-->ScannerWorker0;
    Scanner-->ScannerWorkerN;
```

### Scraper
Scrapes metadata from:
- AWS S3 buckets
- Local disk

### DataStore
Data service, backed by:
- SQLite

### Scanner
Scan assets described by metadata.

## Build & Test
```bash
goreleaser release --snapshot --clean

go test -v ./...
```
