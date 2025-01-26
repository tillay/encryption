import sys, re
from utils import copy, paste, encrypt, decrypt, password_logic, handle_flags, passwd
prefix_text = "&&"

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
