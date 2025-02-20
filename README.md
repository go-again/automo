# AutoMo

`automo` is a lightweight cross-platform utility that prevents your system from going idle by automatically moving your mouse cursor. Perfect for keeping your system active during long downloads, presentations, or when you need to appear "active" in messaging apps.

## Features

- **Two Operating Modes**:
  - **Normal Mode**: Gently moves the mouse cursor in a subtle circular pattern
  - **Zen Mode**: Simulates mouse activity without visible cursor movement
- **Smart Detection**: Only activates after 30 seconds of mouse inactivity
- **Configurable**: Adjust check interval to your needs
- **Resource Friendly**: Minimal CPU and memory usage
- **Simple Console Interface**: Clear status messages and easy controls
- **OS Support**: macOS and Windows

## Security Notes

### Windows Security Note
We do not distribute pre-built binaries because Windows executables should ideally be signed with a valid certificate to prevent being flagged as potentially harmful by security software. Without proper signing, executables may be:
- Blocked by Windows SmartScreen/Defender
- Flagged by antivirus software
- Automatically quarantined
- Trigger false-positive malware alerts

For these reasons, we strongly recommend building the application from source code yourself. This ensures you have full control and transparency over the executable you're running.

### macOS Security Note
macOS requires explicit permission for applications to control your computer. This is a built-in security feature that affects all applications, including `automo`. Before running:

1. You'll need to grant accessibility permissions to `automo`
2. Navigate to **System Settings** > **Privacy & Security** > **Accessibility**
3. Click the "+" button to add `automo` to the list of allowed applications

This security measure protects your system by ensuring that only applications you explicitly trust can control your computer. The permission request is standard for any application that needs to simulate mouse or keyboard input, including common tools like Automator and browsers.

## Getting Started

### Building from Source

#### Prerequisites
- Go 1.20 or later (tested with 1.23+)
- macOS or Windows (10/11 recommended) operating system
- For macOS: Xcode Command Line Tools

#### macOS Development Setup
1. Install Xcode Command Line Tools:
```bash
xcode-select --install
```
This will install essential development tools required for building applications on macOS, including:
- C compiler (needed for CGo)
- Git
- Make and other build tools

If you already have Xcode installed, you can skip this step as the Command Line Tools are included with Xcode.

2. Verify the installation:
```bash
xcode-select -p
```
This should return the path to your Command Line Tools installation.


#### Build Steps

1. Clone the repository:
```bash
git clone https://github.com/go-again/automo.git
cd automo
```

2. Build the application:

For console application (with visible window):
```bash
# macOS
go build -o automo
# Windows
go build -o automo.exe
```

For background application (no visible window):
```bash
# Windows
go build -ldflags -H=windowsgui -o automo_silent.exe
```

## Usage

### Basic Usage

Simply run the executable:
```bash
# macOS
./automo
# Windows
automo.exe
```

### Command Line Options

Run in zen mode (no visible cursor movement):
```bash
# macOS
./automo -zen
# Windows
automo.exe -zen
```

Set custom check interval (in seconds):
```bash
# macOS
./automo -interval 10
# Windows
automo.exe -interval 10
```

Enable debug output:
```bash
# macOS
./automo -debug
# Windows
automo.exe -debug
```

Combine options:
```bash
# macOS
./automo -zen -interval 10 -debug
# Windows
automo.exe -zen -interval 10 -debug
```

### Available Flags
- `-zen`: Enable zen mode (simulates movement without moving cursor)
- `-interval`: Set check interval in seconds (default: 5)
- `-debug`: Enable debug output (shows detected user activity)

### Important Note for macOS Users

macOS requires explicit permission for applications to control your computer. This is a built-in security feature that affects all applications, including AutoMo. Before running:

1. You'll need to grant accessibility permissions to `automo`
2. Navigate to **System Settings** > **Privacy & Security** > **Accessibility**
3. Click the "+" button to add `automo` to the list of allowed applications

This security measure protects your system by ensuring that only applications you explicitly trust can control your computer. The permission request is standard for any application that needs to simulate mouse or keyboard input, including common tools like Automator and browsers.

## How It Works

`automo` monitors your keyboard and mouse activity at regular intervals (default: every 5 seconds). If no activity is detected for 30 seconds:
- In normal mode: Moves the cursor in a small circular pattern (±5 pixels)
- In zen mode: Simulates mouse movement without moving the cursor

The program detects various types of activity:
- Keyboard input (key presses and releases)
- Mouse movement
- Mouse clicks
- Mouse wheel scrolling

The program shows console output to confirm it's working and can be stopped anytime with Ctrl+C.

## Tips
- Use zen mode if you don't want visible cursor movement
- Adjust the interval based on your needs (-interval flag)
- Use debug mode to see what activity is being detected
- Run from command prompt to see status messages
- Use the silent build if you want to run it in the background
- Create a shortcut to `automo.exe` with your preferred flags for quick access

## Troubleshooting

### Common Issues on Windows
1. "Access Denied" error
    - Run the application as administrator
2. No console output
    - Make sure you're not using the silent build
3. Application not working
    - Check if antivirus is blocking the application
    - Verify you have the required Windows permissions
    - Try running from a different directory

### Common Issues on macOS
1. "Operation not permitted" error
    - Check System Settings > Privacy & Security > Accessibility
    - Make sure `automo` is in the allowed applications list
    - Try removing and re-adding the permissions

2. Build fails with CGo errors
    - Verify Xcode Command Line Tools are installed
    - Run `xcode-select --install` if needed
    - Make sure you have the latest OS updates

3. Terminal says "cannot be opened because the developer cannot be verified"
    - Right-click (or Control-click) the application and select "Open"
    - Click "Open" in the confirmation dialog
    - This only needs to be done once

4. Mouse movement not working
    - Check Activity Monitor to verify `automo` is running
    - Try restarting the application
    - Ensure accessibility permissions are granted
    - Check if running from a protected directory (like `/Applications`)

## License

MIT License - Feel free to use and modify as needed.

## Contributing

Contributions are welcome! Feel free to submit issues and pull requests.

## Acknowledgments

Inspired by various mouse jiggler implementations, including [MouseJiggler](https://github.com/arkane-systems/mousejiggler)
