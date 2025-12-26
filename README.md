Server Monitor Go
==========

<img alt="Apache 2.0" src="https://img.shields.io/github/license/LittleJake/server-monitor-go?style=for-the-badge">

<img alt="GitHub Repo stars" src="https://img.shields.io/github/stars/LittleJake/server-monitor-go?style=for-the-badge">

基于Golang & Gin 的服务器监控平台，添加多语言支持。

数据由服务器端Agent采集，Redis存储相关数据。

[Python Agent](https://github.com/LittleJake/server-monitor-script)

[Golang Agent](https://github.com/LittleJake/server-monitor-agent-go)


### 使用

#### 编译

```bash
git clone https://github.com/LittleJake/server-monitor-go
cd server-monitor-go
go build .
cp .env.example .env
```

#### 配置Redis数据源

```bash
vim .env
```

### 界面演示

<img width="2478" height="1254" alt="image" src="https://github.com/user-attachments/assets/c90677aa-5620-48a2-a933-12d35931723e" />

<img width="2478" height="1254" alt="image" src="https://github.com/user-attachments/assets/83a7d317-bbab-4e57-8ce7-92eb253c65cc" />


### Demo

[Demo](https://monitor-next.littlejake.net)

### 数据结构

#### Collection

```json
{
    "Disk": {
        // Based on different mountpoint.
        "mountpoint": {
            "total": "0.00",
            "used": "0.00",
            "free": "0.00",
            "percent": 0
        },
    },
    "Memory": {
        "Mem": {
            "total": "0.00",
            "used": "0.00",
            "free": "0.00",
            "percent": 0.0
        },
        "Swap": {
            "total": "0.00",
            "used": "0.00",
            "free": "0.00",
            "percent": 0.0
        }
    },
    "Load": {
        // Metrics based on platform.
        "metric": 0.0
    },
    "Network": {
        "RX": {
            "bytes": 0,
            "packets": 0
        },
        "TX": {
            "bytes": 0,
            "packets": 0
        }
    },
    "Thermal": {
        // Celsius
        "sensor": 0.0,
    },
    "Battery": {
        "percent": 0.0,
    }
}
```

#### Info

```json
{
    "CPU": "",
    "System Version": "",
    "IPV4": "masked ipv4",
    "IPV6": "masked ipv6",
    "Uptime": "time in readable form",
    "Connection": "",
    "Process": "",
    "Load Average": "",
    "Update Time": "",
    "Country": "extract from ip-api.com",
    "Country Code": "extract from ip-api.com",
    "Throughput": "Gigabytes",
}
```

### 开源协议

[Apache 2.0](LICENSE)

### 鸣谢

[MDUI](https://mdui.org)

### Sponsors

Thanks for the amazing VM server provided by [DartNode](https://dartnode.com?via=1).

 <a href="https://dartnode.com?via=1"><img src="https://raw.githubusercontent.com/LittleJake/LittleJake/master/images/dartnode.png" width="150"></a>


