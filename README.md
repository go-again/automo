# AutoMo

AutoMo is a lightweight Windows utility that prevents your system from going idle by automatically moving your mouse cursor. Perfect for keeping your system active during long downloads, presentations, or when you need to appear "active" in messaging apps.

## Features

- **Two Operating Modes**:
  - **Normal Mode**: Gently moves the mouse cursor in a subtle zigzag pattern
  - **Zen Mode**: Simulates mouse activity without visible cursor movement
- **Smart Detection**: Only activates after 30 seconds of mouse inactivity
- **Configurable**: Adjust check interval to your needs
- **Resource Friendly**: Minimal CPU and memory usage
- **Simple Console Interface**: Clear status messages and easy controls

## Important Security Note

We do not distribute pre-built binaries because Windows executables should ideally be signed with a valid certificate to prevent being flagged as potentially harmful by security software. Without proper signing, executables may be:
- Blocked by Windows SmartScreen/Defender
- Flagged by antivirus software
- Automatically quarantined
- Trigger false-positive malware alerts

For these reasons, we strongly recommend building the application from source code yourself. This ensures you have full control and transparency over the executable you're running.

## Getting Started

### Building from Source

#### Prerequisites
- Go 1.20 or later (tested with 1.20+)
- Windows operating system (Windows 10/11 recommended)

#### Build Steps

1. Clone the repository:
```bash
git clone https://github.com/go-again/automo.git
cd automo
```

2. Build the application:

For console application (with visible window):
```bash
go build -o automo.exe
```

For background application (no visible window):
```bash
go build -ldflags -H=windowsgui -o automo_silent.exe
```

## Usage

### Basic Usage

Simply run the executable:
```bash
automo.exe
```

### Command Line Options

Run in zen mode (no visible cursor movement):
```bash
automo.exe -zen
```

Set custom check interval (in seconds):
```bash
automo.exe -interval 10
```

Combine options:
```bash
automo.exe -zen -interval 10
```

### Available Flags
- `-zen`: Enable zen mode (simulates movement without moving cursor)
- `-interval`: Set check interval in seconds (default: 5)

## How It Works

AutoMo monitors your mouse cursor position at regular intervals (default: every 5 seconds). If no movement is detected for 30 seconds:
- In normal mode: Moves the cursor in a small zigzag pattern (±4 pixels)
- In zen mode: Simulates mouse movement without moving the cursor

The program shows console output to confirm it's working and can be stopped anytime with Ctrl+C.

## Tips
- Use zen mode if you don't want visible cursor movement
- Adjust the interval based on your needs (-interval flag)
- Run from command prompt to see status messages
- Use the silent build if you want to run it in the background
- Create a shortcut to automo.exe with your preferred flags for quick access

## Troubleshooting

### Common Issues
1. "Access Denied" error
   - Run the application as administrator
2. No console output
   - Make sure you're not using the silent build
3. Application not working
   - Check if antivirus is blocking the application
   - Verify you have the required Windows permissions

## License

MIT License - Feel free to use and modify as needed.

## Contributing

Contributions are welcome! Feel free to submit issues and pull requests.

## Acknowledgments

Inspired by various mouse jiggler implementations, including [MouseJiggler](https://github.com/arkane-systems/mousejiggler)
