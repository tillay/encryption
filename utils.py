import sys, subprocess, os, base64, random, string
from Crypto.Util.Padding import pad, unpad
from hashlib import sha256
from Crypto.Cipher import AES

PASSWORD_FILE = "/tmp/key"
prefix_password = "@@"

def copy(data):
    cmd = ["wl-copy"] if os.environ.get("XDG_SESSION_TYPE") == "wayland" else ["xclip", "-selection", "clipboard"]
    subprocess.run(cmd, input=data.encode() if isinstance(data, str) else data)

def passwd():
    return PASSWORD_FILE
def paste():
    cmd = ["wl-paste", "--no-newline"] if os.environ.get("XDG_SESSION_TYPE") == "wayland" else ["xclip", "-o", "-selection", "clipboard"]
    return subprocess.run(cmd, stdout=subprocess.PIPE).stdout

def encrypt(plaintext, passphrase):
    key, iv = sha256(passphrase.encode()).digest(), os.urandom(AES.block_size)
    return base64.b64encode(iv + AES.new(key, AES.MODE_CBC, iv).encrypt(pad(plaintext.encode(), AES.block_size))).decode()
def decrypt(ciphertext, passphrase):
    try:
        decoded = base64.b64decode(ciphertext)
        return unpad(AES.new(sha256(passphrase.encode()).digest(), AES.MODE_CBC, decoded[:AES.block_size]).decrypt(decoded[AES.block_size:]), AES.block_size).decode()
    except (ValueError, KeyError): return None

def password_logic(clipboard_content):
    if clipboard_content.startswith(prefix_password.encode()):
        with open(PASSWORD_FILE, 'w') as f: f.write(clipboard_content.decode())
        print(f"New password saved to {PASSWORD_FILE}.")
        copy("")
        sys.exit(0)
    if os.path.exists(PASSWORD_FILE): return open(PASSWORD_FILE).read()
    try:
        pw = input("Password to store until next reboot: ")
    except KeyboardInterrupt:
        sys.exit(0)
    with open(PASSWORD_FILE, 'w') as f: f.write(pw)
    return pw

def handle_flags():
    if len(sys.argv) > 1:
        if sys.argv[1] == "-n":
            pw = input("Password to store until next reboot: ")
            copy(pw)
            with open(PASSWORD_FILE, 'w') as f: f.write(pw)
            print(f"New password saved to {PASSWORD_FILE} and copied to clipboard.")
            sys.exit(0)
        elif sys.argv[1] == "-p":
            try:
                password = prefix_password+''.join(random.choices(string.ascii_letters + string.digits, k=int(sys.argv[2])))
            except IndexError:
                print("Please a password length after the -p flag")
                sys.exit(1)
            copy(password)
            with open(PASSWORD_FILE, 'w') as f: f.write(password)
            print(f"New RANDOM password saved to {PASSWORD_FILE} and copied to clipboard.")
            sys.exit(0)
handle_flags()
