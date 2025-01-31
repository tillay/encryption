import os, sys, getpass
from pathlib import Path
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
from cryptography.hazmat.primitives import hashes, padding
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes

BLOCK_SIZE, KEY_SIZE, IV_SIZE, SALT_SIZE = 128, 32, 16, 16

def derive_key(password: str, salt: bytes) -> bytes:
    return PBKDF2HMAC(algorithm=hashes.SHA256(), length=KEY_SIZE, salt=salt, iterations=100000, backend=default_backend()).derive(password.encode())

def caesar_cipher(text: str, shift: int) -> str:
    return ''.join(chr((ord(char) + shift) % 256) for char in text)

def encrypt_name(name: str, password: str) -> str:
    shift = sum(ord(char) for char in password) % 256
    return caesar_cipher(name, shift)

def decrypt_name(encrypted_name: str, password: str) -> str:
    shift = sum(ord(char) for char in password) % 256
    return caesar_cipher(encrypted_name, -shift)

def encrypt_file(filepath: Path, password: str):
    salt = os.urandom(SALT_SIZE)
    key = derive_key(password, salt)
    iv = os.urandom(IV_SIZE)
    padder = padding.PKCS7(BLOCK_SIZE).padder()
    cipher = Cipher(algorithms.AES(key), modes.CBC(iv), backend=default_backend())
    encrypted_data = cipher.encryptor().update(padder.update(filepath.read_bytes()) + padder.finalize()) + cipher.encryptor().finalize()
    enc_name = encrypt_name(filepath.name, password)
    (filepath.parent / (enc_name + '.enc')).write_bytes(salt + iv + encrypted_data)
    filepath.unlink()

def decrypt_file(filepath: Path, password: str):
    enc_data = filepath.read_bytes()
    salt, iv, key = enc_data[:SALT_SIZE], enc_data[SALT_SIZE:SALT_SIZE + IV_SIZE], derive_key(password, enc_data[:SALT_SIZE])
    cipher = Cipher(algorithms.AES(key), modes.CBC(iv), backend=default_backend())
    unpadder = padding.PKCS7(BLOCK_SIZE).unpadder()
    decrypted_data = unpadder.update(cipher.decryptor().update(enc_data[SALT_SIZE + IV_SIZE:]) + cipher.decryptor().finalize()) + unpadder.finalize()
    dec_name = decrypt_name(filepath.stem, password)
    (filepath.parent / dec_name).write_bytes(decrypted_data)
    filepath.unlink()

def process(password: str, files):
    operation = "decrypt" if all(file.suffix == '.enc' for file in files if file.is_file()) else "encrypt"
    print("Decrypting..." if operation == "decrypt" else "Encrypting...")
    for file in files:
        if file.is_file():
            try:
                (decrypt_file if operation == "decrypt" and file.suffix.endswith(".enc") else encrypt_file)(file, password)
            except ValueError:
                print("Incorrect password")
                sys.exit(1)

def encrypt_recursive():
    directory = Path(input("Directory to encrypt or decrypt: ").strip()).resolve()
    if not directory.is_dir(): print("Invalid directory!"); sys.exit(1)
    process(getpass.getpass("Enter password: "), list(directory.rglob("*")))
    print("Done!")

def encrypt_singular():
    file_path = Path(input("File to encrypt or decrypt: "))
    if not file_path.is_file(): print("Invalid file path!"); sys.exit(1)
    process(getpass.getpass("Enter password: "), [file_path])
    print("Done!")

if len(sys.argv) != 2:
    print("Usage: python script.py -f (single file) or -r (recursive directory)")
    sys.exit(1)
if sys.argv[1] == "-f":
    encrypt_singular()
elif sys.argv[1] == "-r":
    encrypt_recursive()
else:
    print("Flag -f to operate on a single file, flag -r to operate on all files in a directory recursively.")
