# encryption
Source code for all of the end to end encryption projects I have been doing with some friends

reccomended dependencies:

`xclip, wl-clipboard, python-pyrcyptodome, python-pillow`

minimal.py only requires pycryptodome.

# Usage (encrypt.py)

The default location for stored password is /tmp/key, which is wiped on reboot because its in ram.

This can be changed my modifying the source code for encrypt.py

The python script automatically checks for whats on your clipboard:

- if it starts with `&&` it parses it as encrypted text

- if it starts with `@@` it makes it the new password

Flags:

- `-p`: create new password, overwrite current saved password with it, and copy the new one to the clipboard. Add a number with the length of the password after

- `-n`: manually type in a new password

support is only for Linux with X11 or Wayland (minimal.py works on all os'es)

For sharing passwords between people securely use pgpcli.

# Usage (image.py)

The python script automatically checks for what to do:

- if a file ~/Downloads/message.txt exists, it parses that as an encrypted image

- if the clipboard starts with `££` it parses it as an encrypted image

- if the clipboard is an png image, it encryptes the image (Other image types probably coming soonish maybe)

- if it starts with none of these, it does nothing


# Testing protocol:

For encrypt.py:

See what happens when you run the script with:

- Nothing on your clipboard (should prompt for text)

- Encrypted text from previous test

- Encrypted text already made on previous versions

- Password starting with password prefix

- Weird junk and file formats on clipboard

- `-p` flag (make new random password)

- `-n` flag (prompt manually for new password)

- packages `xclip, wl-clipboard, python-pyrcyptodome` not installed

- Not on a Linux system

- On both Wayland and X11

- On a Linux system but no desktop environment (in tty)

For image.py:

- random junk on clipboard (including bad formats)

- Decrypted PNG image

- Encrypted PNG image from previous test

- Encrypted PNG image from previous version

- Encrypted PNG image from short copy from discord

- Broken PNG image from attempted short copy from discord

- Encrypted PNG image from download button on discord (message.txt)

- Corrupted message.txt
