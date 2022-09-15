from sys import platform
import subprocess

if platform == "linux" or platform == "linux2":
    print("Stopping xampp based on standard Linux path")
    subprocess.run(["sudo", "/opt/lampp/lampp", "stop"])
elif platform == "darwin":
    print("Stopping Xampp based on standard MacOS/OS-X path")
    subprocess.run(["SUDO", "/Applications/XAMPP/xamppfiles/xampp", "stop"])
elif platform == "win32":
    print("Stopping Xampp based on standard Windows path")
    subprocess.run(["stop", "C:\\xampp\\xampp-control.exe"])
