import sys, subprocess, secrets, io, os, base64, random, string, binascii, re
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.backends import default_backend
from Crypto.Cipher import AES
from Crypto.Util.Padding import pad, unpad
from utils import copy, paste, encrypt, decrypt, password_logic, handle_flags
prefix_text = "&&"
PASSWORD_FILE = "/tmp/key"

# Main logic
clipboard_content = paste()
password = password_logic(clipboard_content)
handle_flags()

if prefix_text.encode() in clipboard_content:
    try:
        content = clipboard_content.decode()
        content = re.sub(r"@[^&]*&&|@.*$|<.*?>", "", content).replace(" ", "")
        text = decrypt(content, password) or "Incorrect password."
        print("Decrypted text:", text)
    except UnicodeDecodeError:
        print("had to remove image from clipboard, please run again")
    copy("")
else:
    try:
        text = encrypt(input("Text Input: "), password)
    except KeyboardInterrupt:
        sys.exit(0)
    copy(prefix_text + text)
    print("Encrypted text copied to clipboard.")
