import sys, subprocess, secrets, io, os, base64, random, string, binascii, re
from hashlib import sha256
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.backends import default_backend
from Crypto.Cipher import AES
from Crypto.Util.Padding import pad, unpad

# Constants
PASSWORD_FILE = "/tmp/key"
prefix_text, prefix_password = "&&", "@@"

# Clipboard functions
def copy_to_clipboard(data):
    cmd = ["wl-copy"] if os.environ.get("XDG_SESSION_TYPE") == "wayland" else ["xclip", "-selection", "clipboard"]
    subprocess.run(cmd, input=data.encode() if isinstance(data, str) else data)

def get_from_clipboard():
    cmd = ["wl-paste", "--no-newline"] if os.environ.get("XDG_SESSION_TYPE") == "wayland" else ["xclip", "-o", "-selection", "clipboard"]
    return subprocess.run(cmd, stdout=subprocess.PIPE).stdout

# Password handling
def password_logic():
    if clipboard_content.startswith(prefix_password.encode()):
        with open(PASSWORD_FILE, 'w') as f: f.write(clipboard_content.decode())
        print(f"New password saved to {PASSWORD_FILE}.")
        copy_to_clipboard("")
        sys.exit(0)
    if os.path.exists(PASSWORD_FILE): return open(PASSWORD_FILE).read()
    try:
        pw = input("Password to store until next reboot: ")
    except KeyboardInterrupt:
        sys.exit(0)
    with open(PASSWORD_FILE, 'w') as f: f.write(pw)
    return pw

# Encrypt/Decrypt text
def encrypt(plaintext, passphrase):
    key, iv = sha256(passphrase.encode()).digest(), os.urandom(AES.block_size)
    return base64.b64encode(iv + AES.new(key, AES.MODE_CBC, iv).encrypt(pad(plaintext.encode(), AES.block_size))).decode()

def decrypt(ciphertext, passphrase):
    try:
        decoded = base64.b64decode(ciphertext)
        return unpad(AES.new(sha256(passphrase.encode()).digest(), AES.MODE_CBC, decoded[:AES.block_size]).decrypt(decoded[AES.block_size:]), AES.block_size).decode()
    except (ValueError, KeyError): return None
        
# Main logic
if len(sys.argv) > 1:
    if sys.argv[1] == "-n":
        pw = input("Password to store until next reboot: ")
        copy_to_clipboard(pw)
        with open(PASSWORD_FILE, 'w') as f: f.write(pw)
        print(f"New password saved to {PASSWORD_FILE} and copied to clipboard.")
        sys.exit(0)
    elif sys.argv[1] == "-p":
        try:
            password = prefix_password+''.join(random.choices(string.ascii_letters + string.digits, k=int(sys.argv[2])))
        except IndexError:
            print("Please a password length after the -p flag")
            sys.exit(1)
        copy_to_clipboard(password)
        with open(PASSWORD_FILE, 'w') as f: f.write(password)
        print(f"New random password saved to {PASSWORD_FILE} and copied to clipboard.")
        sys.exit(0)

clipboard_content = get_from_clipboard()
password = password_logic()

if prefix_text.encode() in clipboard_content:
    content = clipboard_content.decode()
    content = re.sub(r"@[^&]*&&|@.*$|<.*?>", "", content).replace(" ", "")
    text = decrypt(content, password) or "Incorrect password."
    print("Decrypted text:", text)
    copy_to_clipboard("")
else:
    try:
        text = encrypt(input("Text Input: "), password)
    except KeyboardInterrupt:
        sys.exit(0)
    copy_to_clipboard(prefix_text + text)
    print("Encrypted text copied to clipboard.")
