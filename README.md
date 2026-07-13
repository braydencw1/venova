## 🧩 Environment Variables

These variables configure the bot, database, Minecraft, RCON, and audio server.

---

### 🔐 Core
| Variable | Description | Default |
|-----------|--------------|----------|
| `TOKEN` | Bot authentication token | — |
| `BOT_ID` | Unique bot ID | — |

---

### 🗄️ Database
| Variable | Description | Default |
|-----------|--------------|----------|
| `DB_HOST` | Database host address | — |
| `DB_DB` | Database name | — |
| `DB_USER` | Database username | — |
| `DB_PASS` | Database password | — |
| `DB_PORT` | Database port | — |

---

### ⛏️ Minecraft / SSH
| Variable | Description | Default |
|-----------|--------------|----------|
| `MC_HOST` | Minecraft server host | — |
| `MC_PORT` | Minecraft server port | — |
| `MC_USER` | SSH username for connecting to MC host | — |
| `MC_SSH_KEY` | SSH private key contents for MC host | — |

---

### 🧰 RCON
| Variable | Description | Default |
|-----------|--------------|----------|
| `RCON_HOST` | RCON host address | — |
| `RCON_PORT` | RCON port | — |
| `RCON_PASS` | RCON password | — |

---

### � Audio Server
| Variable | Description | Default |
|-----------|--------------|----------|
| `AUDIO_SERVER_PORT` | Audio server port | `5005` |

---

### �💡 Example `.env`
```env
TOKEN=your-bot-token
BOT_ID=1234567890

DB_HOST=localhost
DB_DB=mydatabase
DB_USER=myuser
DB_PASS=mypassword
DB_PORT=5432

MC_HOST=minecraft.example.com
MC_PORT=25565
MC_USER=mcadmin
MC_SSH_KEY="-----BEGIN OPENSSH PRIVATE KEY-----\n...\n-----END OPENSSH PRIVATE KEY-----"

RCON_HOST=127.0.0.1
RCON_PORT=25575
RCON_PASS=secret

AUDIO_SERVER_PORT=5005
