## SQLite Schema
```sql
CREATE TABLE dhcp_leases (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    mac_address TEXT UNIQUE NOT NULL,
    ip_address TEXT UNIQUE NOT NULL
);
```
