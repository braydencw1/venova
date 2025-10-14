# 🎧 Audio Streaming Service

This service captures local audio (via PulseAudio, ALSA, or VB-Cable) and streams it to the audio server using FFmpeg.

## 🔧 Environment Variables

| Variable | Description | Default |
|-----------|--------------|----------|
| `AUDIO_FORMAT` | Audio backend (`pulse`, `alsa`, or `dshow` on Windows) | `pulse` |
| `AUDIO_DEVICE` | Input device name (e.g. `VirtualCable.monitor` or `default`) | `default` |
| `AUDIO_SERVER_IP` | Target server IP | `127.0.0.1` |
| `AUDIO_SERVER_PORT` | Target server port | `5005` |
| `FFMPEG_PATH` | Path to FFmpeg executable | `ffmpeg` |

## 💻 Platform Examples

**Windows**
```cmd
set AUDIO_FORMAT=dshow
set AUDIO_DEVICE=audio=CABLE Output (VB-Audio Virtual Cable)
```

**Linux**
```cmd
export AUDIO_FORMAT=pulse
export AUDIO_DEVICE=VirtualCable.monitor
```