import os, sys, getpass
from pathlib import Path
from cryptography.hazmat.primitives.kdf.pbkdf2 import PBKDF2HMAC
from cryptography.hazmat.primitives import hashes
from cryptography.hazmat.backends import default_backend
from cryptography.hazmat.primitives.ciphers import Cipher, algorithms, modes
from cryptography.hazmat.primitives import padding

BLOCK_SIZE = 128
KEY_SIZE = 32
IV_SIZE = 16
SALT_SIZE = 16

def derive_key(password: str, salt: bytes) -> bytes:
    kdf = PBKDF2HMAC(
        algorithm=hashes.SHA256(),
        length=KEY_SIZE,
        salt=salt,
        iterations=100000,
        backend=default_backend()
    )
    return kdf.derive(password.encode())

def encrypt_name(name: str, key: bytes) -> bytes:
    iv = os.urandom(IV_SIZE)
    cipher = Cipher(algorithms.AES(key), modes.CBC(iv), backend=default_backend())
    encryptor = cipher.encryptor()
    padder = padding.PKCS7(BLOCK_SIZE).padder()
    padded_name = padder.update(name.encode()) + padder.finalize()
    encrypted_name = encryptor.update(padded_name) + encryptor.finalize()
    return iv + encrypted_name

def decrypt_name(encrypted_name: bytes, key: bytes) -> str:
    iv = encrypted_name[:IV_SIZE]
    encrypted_data = encrypted_name[IV_SIZE:]
    cipher = Cipher(algorithms.AES(key), modes.CBC(iv), backend=default_backend())
    decryptor = cipher.decryptor()
    decrypted_name = decryptor.update(encrypted_data) + decryptor.finalize()
    unpadder = padding.PKCS7(BLOCK_SIZE).unpadder()
    unpadded_name = unpadder.update(decrypted_name) + unpadder.finalize()
    return unpadded_name.decode().rstrip('\0')

def encrypt_file(filepath: Path, password: str):
    with open(filepath, 'rb') as file:
        file_data = file.read()
    salt = os.urandom(SALT_SIZE)
    key = derive_key(password, salt)
    encrypted_name = encrypt_name(filepath.name, key)
    enc_filepath = filepath.parent / (encrypted_name.hex() + '.enc')
    iv = os.urandom(IV_SIZE)
    padder = padding.PKCS7(BLOCK_SIZE).padder()
    padded_data = padder.update(file_data) + padder.finalize()
    cipher = Cipher(algorithms.AES(key), modes.CBC(iv), backend=default_backend())
    encryptor = cipher.encryptor()
    encrypted_data = encryptor.update(padded_data) + encryptor.finalize()
    with open(enc_filepath, 'wb') as enc_file:
        enc_file.write(salt + iv + encrypted_data)
    filepath.unlink()

def decrypt_file(filepath: Path, password: str):
    with open(filepath, 'rb') as enc_file:
        enc_data = enc_file.read()
    salt = enc_data[:SALT_SIZE]
    iv = enc_data[SALT_SIZE:SALT_SIZE + IV_SIZE]
    encrypted_data = enc_data[SALT_SIZE + IV_SIZE:]
    key = derive_key(password, salt)
    encrypted_name = bytes.fromhex(filepath.stem)
    decrypted_name = decrypt_name(encrypted_name, key)
    dec_filepath = filepath.parent / decrypted_name
    cipher = Cipher(algorithms.AES(key), modes.CBC(iv), backend=default_backend())
    decryptor = cipher.decryptor()
    decrypted_data = decryptor.update(encrypted_data) + decryptor.finalize()
    unpadder = padding.PKCS7(BLOCK_SIZE).unpadder()
    original_data = unpadder.update(decrypted_data) + unpadder.finalize()
    with open(dec_filepath, 'wb') as dec_file:
        dec_file.write(original_data)
    filepath.unlink()

def process_files_in_directory(directory: Path, password: str):
    files = list(directory.rglob("*"))
    if all(file.suffix == '.enc' for file in files if file.is_file()):
        operation = "decrypt"
        print("Decrypting...")
    else:
        operation = "encrypt"
        print("Encrypting...")
    for file in files:
        if file.is_file():
            if operation == "encrypt" and not file.suffix.endswith(".enc"):
                encrypt_file(file, password)
            elif operation == "decrypt" and file.suffix.endswith(".enc"):
                try:
                    decrypt_file(file, password)
                except ValueError:
                    print("incorrect password")
                    sys.exit(1)

directory_path = input("Directory to encrypt or decrypt: ").strip()
password = getpass.getpass("Enter password: ")
directory = Path(directory_path).resolve()
if not directory.is_dir():
    print("Invalid directory!")
    sys.exit(1)
process_files_in_directory(directory, password)
print("done!")
