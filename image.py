import sys, subprocess, io, os, base64, secrets
from PIL import Image, UnidentifiedImageError
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
# Constants
PASSWORD_FILE = "/tmp/key"
prefix_image = "££"
decrypted_image_path = "/tmp/decrypted_image.png"
message_file_path = os.path.expanduser("~/Downloads/message.txt")

# Open image in browser
def open_image():
    subprocess.Popen('xdg-open /tmp/decrypted_image.png', shell=True)

# Clipboard functions
def copy_to_clipboard(data):
    cmd = ["wl-copy"] if os.environ.get("XDG_SESSION_TYPE") == "wayland" else ["xclip", "-selection", "clipboard"]
    subprocess.run(cmd, input=data.encode() if isinstance(data, str) else data)

def get_from_clipboard():
    cmd = ["wl-paste", "--no-newline"] if os.environ.get("XDG_SESSION_TYPE") == "wayland" else ["xclip", "-o", "-selection", "clipboard"]
    return subprocess.run(cmd, stdout=subprocess.PIPE).stdout

def derive_key(password):
    kdf = PBKDF2HMAC(algorithm=hashes.SHA256(), length=32, salt=password.encode(), iterations=100000, backend=default_backend())
    return kdf.derive(password.encode())

def encrypt_image(image, key):
    iv = secrets.token_bytes(16)
    img_data = io.BytesIO(); image.save(img_data, format="PNG")
    cipher = Cipher(algorithms.AES(key), modes.CFB8(iv), backend=default_backend())
    encrypted_data = cipher.encryptor().update(img_data.getvalue()) + cipher.encryptor().finalize()
    return base64.b64encode(iv + encrypted_data)

def decrypt_image(encrypted_data, key):
    decoded = base64.b64decode(encrypted_data)
    iv, encrypted_img = decoded[:16], decoded[16:]
    cipher = Cipher(algorithms.AES(key), modes.CFB8(iv), backend=default_backend())
    decrypted_data = cipher.decryptor().update(encrypted_img) + cipher.decryptor().finalize()
    return Image.open(io.BytesIO(decrypted_data))

def process_message_file():
    if os.path.exists(message_file_path):
        with open(message_file_path, 'rb') as f:
            return f.read()
    return None

clipboard_content = process_message_file() or get_from_clipboard()
if not os.path.exists(PASSWORD_FILE):
    password = input("Enter password: ")
    with open(PASSWORD_FILE, 'w') as file:
        file.write(password)
else:
    with open(PASSWORD_FILE, 'r') as file:
        password = file.read()
if clipboard_content.startswith(prefix_image.encode()):
    try:
        key = derive_key(password)
        encrypted_image = clipboard_content[len(prefix_image):]
        decrypted_image = decrypt_image(encrypted_image, key)
        decrypted_image.save(decrypted_image_path)
        copy_to_clipboard(decrypted_image_path)
        print(f"Decrypted image copied to clipboard")
        open_image()
        if clipboard_content == process_message_file():
            os.remove(message_file_path)
        if not os.environ.get("XDG_SESSION_TYPE") == "wayland":
            os.system(f"cat {decrypted_image_path} | xclip -selection clipboard -target image/png -i")
            os.remove(decrypted_image_path)
    except (UnidentifiedImageError, ValueError, OSError):
        print("click the funny little download button to get the encrypted image")
        sys.exit(1)
else:
    try:
        image = Image.open(io.BytesIO(clipboard_content))
        key = derive_key(password)
        encrypted_image = encrypt_image(image, key)
        copy_to_clipboard(prefix_image.encode() + encrypted_image)
        print("Image encrypted and copied to clipboard.")
    except (UnidentifiedImageError, ValueError) as e:
        print("Error: Clipboard does not contain a valid encrypted or decrypted image. Use encrypt.py for text encryption.")
        sys.exit(1)
