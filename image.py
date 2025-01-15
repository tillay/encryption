import sys, subprocess, io, os, base64, secrets
from PIL import Image, UnidentifiedImageError
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from utils import copy, paste, password_logic, handle_flags
prefix_image = "££"
PASSWORD_FILE = "/tmp/key"
decrypted_image_path = "/tmp/decrypted_image.png"
message_file_path = os.path.expanduser("~/Downloads/message.txt")
image_viewer = ("")
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
clipboard_content = process_message_file() or paste()
password = password_logic(clipboard_content)

if clipboard_content.startswith(prefix_image.encode()):
    try:
        key = derive_key(password)
        encrypted_image = clipboard_content[len(prefix_image):]
        decrypted_image = decrypt_image(encrypted_image, key)
        decrypted_image.save(decrypted_image_path)
        if (image_viewer == ""):
            os.system(f'xdg-open {decrypted_image_path}')
            xdg_default = subprocess.run(["xdg-mime", "query", "default", "image/png"], capture_output=True, text=True).stdout.strip()
            print(f"Opened decrypted image with {xdg_default}")
        else:
            print(f"Opened decrypted image with {image_viewer}")
            os.system(f"{image_viewer} {decrypted_image_path}")
        try:
            os.remove(message_file_path)
        except:
            print(f"Somehow managed to decrypt image from clipboard. Use the download button next time")
    except (UnidentifiedImageError, ValueError, OSError):
        print("click the funny little download button to get the encrypted image")
        sys.exit(1)
else:
    try:
        image = Image.open(io.BytesIO(clipboard_content))
        key = derive_key(password)
        encrypted_image = encrypt_image(image, key)
        copy(prefix_image.encode() + encrypted_image)
        print("Image encrypted and copied to clipboard.")
    except (UnidentifiedImageError, ValueError) as e:
        print("Error: Clipboard does not contain a valid encrypted or decrypted image. Use encrypt.py for text encryption.")
        sys.exit(1)
