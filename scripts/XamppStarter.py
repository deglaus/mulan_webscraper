from sys import platform
import subprocess

if platform == "linux" or platform == "linux2":
    print("Starting xampp based on standard Linux path")
    subprocess.run(["sudo", "/opt/lampp/lampp", "start"])
elif platform == "darwin":
    print("Starting Xampp based on standard MacOS/OS-X path")
    subprocess.run(["SUDO", "/Applications/XAMPP/xamppfiles/xampp", "start"])
elif platform == "win32":
    print("Starting Xampp based on standard Windows path")
    subprocess.run(["start", "C:\\xampp\\xampp-control.exe"])
