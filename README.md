<a id="readme-top"></a>

<h1 align="center">DNS Forwarder</h1>

<p align="center">
  A simple and efficient DNS forwarder written in Go!
  <br />
  <br />
  <a href="https://github.com/kartmos/dns-forwarder/issues/new?labels=bug&template=bug-report.md">Report a Bug</a>
  &middot;
  <a href="https://github.com/kartmos/dns-forwarder/issues/new?labels=enhancement&template=feature-request.md">Request a Feature</a>
</p>

<!-- TABLE OF CONTENTS -->
<details>
  <summary>Table of Contents</summary>
  <ol>
    <li><a href="#about-the-project">About the Project</a></li>
    <li><a href="#features">Features</a></li>
    <li><a href="#getting-started">Getting Started</a></li>
    <li><a href="#usage">Usage</a></li>
  </ol>
</details>

<!-- ABOUT THE PROJECT -->
## About the Project

DNS Forwarder is a simple and efficient DNS forwarder written in Go that redirects DNS queries to a specified DNS server. The project demonstrates the basics of working with network requests and the DNS protocol in Go. Key features include:

* UDP support
* Configurable DNS server address
* Request logging
* Simple configuration

<p align="right">(<a href="#readme-top">back to top</a>)</p>

### Features

- DNS query forwarding
- UDP support
- Configurable DNS server address
- Request logging
- Simple configuration

<!-- GETTING STARTED -->
## Getting Started

### Prerequisites

- Go 1.24.0+
- Internet connection

### Installation

1. Clone the repository.
2. Navigate to the project directory.
3. Build the binary in the project's root folder.

<p align="right">(<a href="#readme-top">back to top</a>)</p>

<!-- USAGE -->
## Usage with Docker-compose

```zsh
docker-compose -f build/deploy/docker-compose.yml up -d
```
<p align="right">(<a href="#readme-top">back to top</a>)</p> 