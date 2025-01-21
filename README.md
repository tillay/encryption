# encryption
Source code for all of the end to end encryption projects I have been doing with some friends

reccomended dependencies:

`xclip, wl-clipboard, python-pycryptodome, python-pillow (for image)`

minimal.py only requires pycryptodome.

# General Usage
to install:

```git clone https://github.com/tillay8/encryption&&cd encryption```

to run a python script:

`python3 <script.py> <optional flag>`

running a go script (pgpcli) is explained later

note: support of encrypt.py and image.py are only for Linux with X11 or Wayland due to clipboard utility needs

# Usage (encrypt.py)

The default location for stored password is /tmp/key, which is wiped on reboot because its in ram.

This can be changed my modifying the source code for encrypt.py

The python script automatically checks for whats on your clipboard:

- if it starts with `&&` it parses it as encrypted text

- if it starts with `@@` it makes it the new password

# Usage (image.py)

The python script automatically checks for what to do:

- if a file ~/Downloads/message.txt exists, it parses that as an encrypted image

- if the clipboard starts with `££` it parses it as an encrypted image

- if the clipboard is an png image, it encryptes the image (Other image types probably coming soonish maybe)

- if it starts with none of these, it does nothing

# Flags 

these work on image.py, encrypt.py, and utils.py

- `-p`: create new password, overwrite current saved password with it, and copy the new one to the clipboard. Add a number with the length of the password after. Example: ```python3 utils.py -p 420```

- `-n`: manually type in a new password

# Usage (minimal.py)

minimal.py is designed to work on a lot more OS'es due to no clipboard integration.

to use, type your message into the text box. if its unencrypted, it will encrypt it, and if it is encrypted, it will decrypt it. 

# Usage (pgpcli)

For sharing passwords between people securely use pgpcli.

pgpcli uses the PGP asymetric standard for encryption and provides a wrapper for PGP encryption

to compile:

make sure you are in the pgpcli subdirectory, then run ```go build -o pgpcli```

this will create a pgpcli binary, which can be "dot slashed" or run by ```./pgpcli```

The menu is pretty straightfoward. By default, keys are stored in ~/wpgp, but that can be changed

pgpcli also checks your clipboard to see if you have a public or private key on your clipboard, and will run the proper steps if so. 
